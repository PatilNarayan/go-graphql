package roles

import (
	"context"
	"errors"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/constants"
	"iam_services_main_v1/internal/dto"
	middleware "iam_services_main_v1/internal/middlewares"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	DB *gorm.DB
}

// CreateRole handles creating a new role.
func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.CreateRoleInput) (*models.Role, error) {

	// Extract gin.Context from GraphQL context
	ginCtx, ok := ctx.Value(middleware.GinContextKey).(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unable to get gin context")
	}

	tenantID, err := helpers.GetTenant(ginCtx)
	if err != nil {
		return nil, err
	}

	if err := ValidateTenantID(uuid.MustParse(tenantID)); err != nil {
		return nil, err
	}

	if err := validateCreateRoleInput(input); err != nil {
		return nil, err
	}

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where(&dto.Mst_ResourceTypes{Name: "Role", RowStatus: 1}).First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}
	tenantIDUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}
	// check if role already exists
	var roleExists dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{Name: input.Name, RowStatus: 1, ResourceTypeID: resourceType.ResourceTypeID, TenantID: &tenantIDUUID}).First(&roleExists).Error; err == nil {
		return nil, fmt.Errorf("role already exists")
	}

	if err := CheckPermissions(input.Permissions); err != nil {
		return nil, err
	}

	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where(&dto.Mst_ResourceTypes{ResourceTypeID: input.AssignableScopeRef, RowStatus: 1}).First(&assignableScopeRef).Error; err != nil {
		return nil, fmt.Errorf("invalid TenantID")
	}

	permitMap := make(map[string]interface{})

	newRoleID := uuid.New()
	permitMap = map[string]interface{}{
		"name": input.Name,
		"key":  newRoleID.String(),
	}

	if input.Description != nil {
		permitMap["description"] = *input.Description
	}

	//create role in permit
	pc := permit.NewPermitClient()
	_, err = pc.SendRequest(ctx, "POST", fmt.Sprintf("resources/%s/roles", assignableScopeRef.Name), permitMap)

	if err != nil {
		return nil, err
	}

	resources, err := pc.SendRequest(ctx, "GET", "resources", nil)
	if err != nil {
		return nil, err
	}

	actions := utils.GetActionMap(resources.([]interface{}), assignableScopeRef.Name)
	permission, err := GetPermissionAction(input.Permissions)
	if err != nil {
		return nil, err
	}
	res := utils.CreateActionMap(actions, permission)
	update := map[string]interface{}{
		"actions": res,
	}
	_, err = pc.SendRequest(ctx, "PATCH", fmt.Sprintf("resources/%s", assignableScopeRef.Name), update)
	if err != nil {
		return nil, err
	}

	if err := r.DB.Create(&dto.TenantResource{
		ResourceID:     newRoleID,
		ResourceTypeID: resourceType.ResourceTypeID,
		Name:           input.Name,
		RowStatus:      1,
		TenantID:       &tenantIDUUID,
		CreatedBy:      constants.DefaltCreatedBy,
		UpdatedBy:      constants.DefaltCreatedBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}).Error; err != nil {
		return nil, err
	}
	role := prepareRoleObject(input)

	role.ResourceID = newRoleID
	if err := r.DB.Create(&role).Error; err != nil {
		return nil, err
	}

	for _, permissionID := range input.Permissions {
		tx := r.DB.Create(&dto.TNTRolePermission{
			ID:           uuid.New(),
			RoleID:       role.ResourceID,
			PermissionID: uuid.MustParse(permissionID),
			RowStatus:    1,
			CreatedBy:    constants.DefaltCreatedBy,
			UpdatedBy:    constants.DefaltCreatedBy,
		})

		if tx.Error != nil {
			return nil, tx.Error
		}
	}
	if len(input.Permissions) > 0 {
		err = SetPermission(ctx, role.ResourceID.String(), assignableScopeRef.Name, permission)
		if err != nil {
			return nil, err
		}
	}
	return convertRoleToGraphQL(&role), nil
}

// UpdateRole handles updating an existing role.
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, input models.UpdateRoleInput) (*models.Role, error) {

	// Extract gin.Context from GraphQL context
	ginCtx, err := helpers.GetGinContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get gin context")
	}

	tenantID, err := helpers.GetTenant(ginCtx)
	if err != nil {
		return nil, err
	}

	if err := ValidateTenantID(uuid.MustParse(tenantID)); err != nil {
		return nil, err
	}

	var role dto.TNTRole
	if err := r.validateUpdateRoleInput(input); err != nil {
		return nil, err
	}

	// Update fields
	if input.Name != "" {
		role.Name = input.Name
	}

	if err := CheckPermissions(input.Permissions); err != nil {
		return nil, err
	}

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	permitMap := make(map[string]interface{})

	permitMap = map[string]interface{}{
		"name": input.Name,
	}

	if input.Description != nil {
		permitMap["description"] = *input.Description
	}

	//create role in permit
	pc := permit.NewPermitClient()

	resource, err := GetResourceTypeIDName(input.AssignableScopeRef)
	if err != nil {
		return nil, err
	}

	if _, err := pc.SendRequest(ctx, "PATCH", fmt.Sprintf("resources/%s/roles/%s", *resource, input.ID.String()), permitMap); err != nil {
		return nil, err
	}

	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where("resource_type_id = ? AND row_status = 1", input.AssignableScopeRef).First(&assignableScopeRef).Error; err != nil {
		return nil, fmt.Errorf("invalid TenantID")
	}

	resources, err := pc.SendRequest(ctx, "GET", "resources", nil)
	if err != nil {
		return nil, err
	}

	actions := utils.GetActionMap(resources.([]interface{}), assignableScopeRef.Name)
	permission, err := GetPermissionAction(input.Permissions)
	if err != nil {
		return nil, err
	}

	updateData := map[string]interface{}{
		"name":       input.Name,
		"version":    input.Version,
		"role_type":  dto.RoleTypeEnumCustom,
		"updated_by": constants.DefaltCreatedBy,
		"updated_at": time.Now(),
	}

	if input.Name != "" {
		updateData["name"] = input.Name
	}
	if input.Description != nil {
		updateData["description"] = *input.Description
	}

	if err := r.DB.Model(&role).Where("resource_id = ? AND row_status = 1", input.ID).UpdateColumns(updateData).Error; err != nil {
		return nil, err
	}

	var pdata []dto.TNTRolePermission
	if err := r.DB.Where("role_id = ? AND row_status = 1", input.ID).Find(&pdata).Error; err != nil {
		return nil, err
	}
	newPermissions := make([]string, 0)
	for _, pid := range input.Permissions {
		exists := false
		for _, p := range pdata {
			if p.PermissionID.String() == pid {
				exists = true
				break
			}
		}
		if !exists {
			newPermissions = append(newPermissions, pid)
			tx := r.DB.Create(&dto.TNTRolePermission{
				ID:           uuid.New(),
				RoleID:       input.ID,
				PermissionID: uuid.MustParse(pid),
				RowStatus:    1,
				CreatedBy:    constants.DefaltCreatedBy,
				UpdatedBy:    constants.DefaltCreatedBy,
			})

			if tx.Error != nil {
				return nil, tx.Error
			}

		}
	}

	newPermissionsActions, err := GetPermissionAction(newPermissions)
	if err != nil {
		return nil, err
	}

	newPermissions = append(newPermissions, permission...)

	res := utils.CreateActionMap(actions, newPermissionsActions)

	removeIDs := make([]uuid.UUID, 0)
	removeIDsList := make([]string, 0)
	for _, p := range pdata {
		exists := false
		for _, pid := range input.Permissions {
			if p.PermissionID.String() == pid {
				exists = true
				break
			}
		}
		if !exists {
			removeIDs = append(removeIDs, p.ID)
			removeIDsList = append(removeIDsList, p.PermissionID.String())
		}
	}

	removeActions, err := GetPermissionAction(removeIDsList)
	if err != nil {
		return nil, err
	}

	for _, id := range removeActions {
		if res[id] != nil {
			delete(res, id)
		}
	}

	update := map[string]interface{}{
		"actions": res,
	}

	_, err = pc.SendRequest(ctx, "PATCH", fmt.Sprintf("resources/%s", assignableScopeRef.Name), update)
	if err != nil {
		return nil, err
	}
	for _, id := range removeIDs {
		if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_permission_id = ? AND row_status = 1", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
			return nil, fmt.Errorf("failed to delete role: %w", err)
		}
	}
	if len(input.Permissions) > 0 {
		err = SetPermission(ctx, input.ID.String(), assignableScopeRef.Name, newPermissions)
		if err != nil {
			return nil, err
		}
	}

	finalRole := dto.TNTRole{}

	if err := r.DB.First(&finalRole, "resource_id = ? AND row_status = 1", input.ID).Error; err != nil {
		return nil, err
	}

	return convertRoleToGraphQL(&finalRole), nil
}

// DeleteRole handles deleting a role by ID.
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, id uuid.UUID) (bool, error) {

	// Extract gin.Context from GraphQL context
	ginCtx, ok := ctx.Value(middleware.GinContextKey).(*gin.Context)
	if !ok {
		return false, fmt.Errorf("unable to get gin context")
	}

	// Retrieve x-tenant-id from headers
	tenantID := ginCtx.GetHeader("tenantID")
	if tenantID == "" {
		return false, errors.New("tenantID not found in headers")
	}

	//validate uuid format
	if _, err := uuid.Parse(tenantID); err != nil {
		return false, fmt.Errorf("invalid tenantID: %w", err)
	}

	if err := ValidateTenantID(uuid.MustParse(tenantID)); err != nil {
		return false, err
	}

	if err := ValidateRoleID(id); err != nil {
		return false, fmt.Errorf("invalid ID: %w", err)
	}

	var roleDB dto.TNTRole
	if err := r.DB.First(&roleDB, "resource_id = ? AND row_status = 1", id).Error; err != nil {
		return false, errors.New("role not found")
	}

	updates := map[string]interface{}{
		"row_status": 0,
	}
	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where("resource_type_id = ? AND row_status = 1", roleDB.ScopeResourceTypeID).First(&assignableScopeRef).Error; err != nil {
		return false, fmt.Errorf("invalid TenantID")
	}

	pc := permit.NewPermitClient()
	// var pid interface{}
	// if data, err := pc.SendRequest(ctx, "GET", fmt.Sprintf("resources/%s/roles", assignableScopeRef.Name), nil); err != nil {
	// 	return false, err
	// } else {
	// 	data := data.([]interface{})
	// 	for _, v := range data {
	// 		v := v.(map[string]interface{})
	// 		if v["name"] == roleDB.Name {
	// 			pid = v["id"]
	// 			break
	// 		}
	// 	}
	// }

	if _, err := pc.SendRequest(ctx, "DELETE", fmt.Sprintf("resources/%s/roles/%s", assignableScopeRef.Name, roleDB.ResourceID.String()), nil); err != nil {
		return false, err
	}

	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ? AND row_status = 1", id).Updates(updates).Error; err != nil {
		return false, err
	}

	if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_id = ? AND row_status = 1", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		return false, fmt.Errorf("failed to delete role: %w", err)
	}

	if err := r.DB.Model(&dto.TenantResource{}).Where("resource_id = ? AND row_status = 1", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		return false, fmt.Errorf("failed to delete role: %w", err)
	}

	actions, err := GetPermissionResourceAction(*&roleDB.ScopeResourceTypeID)
	if err != nil {
		return false, err
	}

	updateres := make(map[string]interface{})
	res := utils.CreateActionMap(updateres, actions)
	update := map[string]interface{}{
		"actions": res,
	}
	_, err = pc.SendRequest(ctx, "PATCH", fmt.Sprintf("resources/%s", assignableScopeRef.Name), update)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CheckPermissions(permissionsIds []string) error {
	//validate permissionIds
	for _, permissionID := range permissionsIds {
		if err := ValidatePermissionID(permissionID); err != nil {
			return fmt.Errorf("invalid permission ID: %w", err)
		}
	}
	return nil
}

func CreateMstRole(tenantID uuid.UUID) error {
	var mstRole []*dto.MstRole

	if err := config.DB.Find(&mstRole).Error; err != nil {
		return err
	}

	var mstRolePermissions []*dto.MstRolePermission
	if err := config.DB.Find(&mstRolePermissions).Error; err != nil {
		return err
	}
	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		return fmt.Errorf("resource type not found: %w", err)
	}

	ResDefaultPermissions, err := AddDefaultPermissions()
	if err != nil {
		return err
	}

	for _, mrole := range mstRole {
		tenantResource := &dto.TenantResource{
			ResourceID:       uuid.New(),
			ParentResourceID: &tenantID,
			ResourceTypeID:   resourceType.ResourceTypeID,
			Name:             mrole.Name,
			RowStatus:        1,
			CreatedBy:        constants.DefaltCreatedBy,
			UpdatedBy:        constants.DefaltCreatedBy,
			CreatedAt:        time.Now(),
			TenantID:         &tenantID,
		}

		if err := config.DB.Save(tenantResource).Error; err != nil {
			return err
		}

		role := dto.TNTRole{
			ResourceID:          tenantResource.ResourceID,
			Name:                mrole.Name,
			Description:         mrole.Description,
			RoleType:            dto.RoleTypeEnumDefault,
			ScopeResourceTypeID: resourceType.ResourceTypeID,
			Version:             mrole.Version,
			CreatedAt:           mrole.CreatedAt,
			CreatedBy:           mrole.CreatedBy,
			UpdatedBy:           mrole.UpdatedBy,
			UpdatedAt:           mrole.UpdatedAt,
			RowStatus:           mrole.RowStatus,
		}

		if err := config.DB.Save(&role).Error; err != nil {
			return err
		}

		for _, permissionID := range mstRolePermissions {
			if permissionID.RoleID != mrole.RoleID {
				continue
			}
			tx := config.DB.Create(&dto.TNTRolePermission{
				ID:           uuid.New(),
				RoleID:       role.ResourceID,
				PermissionID: ResDefaultPermissions[permissionID.PermissionID],
				RowStatus:    1,
				CreatedBy:    constants.DefaltCreatedBy,
				UpdatedBy:    constants.DefaltCreatedBy,
			})

			if tx.Error != nil {
				return tx.Error
			}
		}
	}

	return nil

}

func AddDefaultPermissions() (map[uuid.UUID]uuid.UUID, error) {
	var mstPermissions []*dto.MstPermission

	if err := config.DB.Where("row_status = 1").Find(&mstPermissions).Error; err != nil {
		return nil, err
	}

	res := make(map[uuid.UUID]uuid.UUID)
	for _, mpermission := range mstPermissions {
		permission := &dto.TNTPermission{
			PermissionID: uuid.New(),
			Name:         mpermission.Name,
			ServiceID:    mpermission.ServiceID,
			Action:       mpermission.Action,
			RowStatus:    1,
			// RoleID:       *input.RoleID,
			CreatedBy: constants.DefaltCreatedBy,
			UpdatedBy: constants.DefaltUpdatedBy,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		res[mpermission.PermissionID] = permission.PermissionID
		if err := config.DB.Create(permission).Error; err != nil {
			return nil, err
		}
	}

	return res, nil
}

func DeleteDefaultRole(tenantID uuid.UUID) error {
	// Find resource type for Role
	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		return fmt.Errorf("resource type not found: %w", err)
	}

	// Start a transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var roleResources []dto.TenantResource
	if err := tx.Where("tenant_id = ? AND resource_type_id = ? AND row_status = 1", tenantID, resourceType.ResourceTypeID).Find(&roleResources).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update row status in tenant_resources table
	if err := tx.Model(&dto.TenantResource{}).
		Where("tenant_id = ? AND resource_type_id = ? AND row_status = 1", tenantID, resourceType.ResourceTypeID).
		Update("row_status", 0).Error; err != nil {
		tx.Rollback()
		return err
	}

	roleIDs := make([]uuid.UUID, len(roleResources))
	for i, roleResource := range roleResources {
		roleIDs[i] = roleResource.ResourceID

	}

	// Update row status in tnt_roles table
	if err := tx.Model(&dto.TNTRole{}).
		Where("resource_id IN ? AND row_status = 1", roleIDs).
		Update("row_status", 0).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update row status in tnt_role_permissions table
	if err := tx.Model(&dto.TNTRolePermission{}).
		Where("role_id IN ? AND row_status = 1", roleIDs).
		Update("row_status", 0).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func SetPermission(ctx context.Context, roleName, assinableScopeName string, permissionAction []string) error {
	pc := permit.NewPermitClient()

	update := map[string]interface{}{
		"permissions": permissionAction,
	}

	_, err := pc.SendRequest(ctx, "POST", fmt.Sprintf("resources/%s/roles/%s/permissions", assinableScopeName, roleName), update)
	if err != nil {
		return err
	}

	return nil
}

func ValidateRoleID(roleId uuid.UUID) error {
	// Check if the resource ID exists in the database
	var count int64
	if err := config.DB.Model(&dto.TNTRole{}).Where("resource_id = ? AND row_status = 1", roleId).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func ValidateRole(resourceID uuid.UUID) error {
	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		return fmt.Errorf("resource type not found: %w", err)
	}
	var count int64
	if err := config.DB.Model(&dto.TenantResource{}).
		// Where("resource_id = ? AND row_status = 1 AND resource_type_id IN (?)", resourceID, resourceIds).
		Where("resource_id = ? AND row_status = 1 AND resource_type_id = ?", resourceID, resourceType.ResourceTypeID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func GetResourceTypeIDName(resourceID uuid.UUID) (*string, error) {
	var data dto.Mst_ResourceTypes
	if err := config.DB.Model(&dto.Mst_ResourceTypes{}).Where("resource_type_id = ? AND row_status = 1", resourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data.Name, nil
}

func GetPermissionAction(permissionsIds []string) ([]string, error) {
	var actions []string
	//validate permissionIds
	for _, permissionID := range permissionsIds {
		var data dto.TNTPermission
		if err := config.DB.Model(&dto.TNTPermission{}).Where("permission_id = ? AND row_status = 1", permissionID).First(&data).Error; err != nil {
			return nil, err
		}
		actions = append(actions, data.Action)
	}
	return actions, nil
}

func GetPermissionResourceAction(resourceID uuid.UUID) ([]string, error) {
	var resources []dto.TenantResource
	if err := config.DB.Model(&dto.TenantResource{}).Where("parent_resource_id = ? AND row_status = 1", resourceID).Find(&resources).Error; err != nil {
		return nil, err
	}

	resourceIds := []string{}
	for _, resource := range resources {
		resourceIds = append(resourceIds, resource.ResourceID.String())
	}

	var rolePermissions []dto.TNTRolePermission
	if err := config.DB.Model(&dto.TNTRolePermission{}).Where("role_id in (?) AND row_status = 1", resourceIds).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	permissionIds := []string{}
	for _, rolePermission := range rolePermissions {
		permissionIds = append(permissionIds, rolePermission.PermissionID.String())
	}

	actions, err := GetPermissionAction(permissionIds)
	if err != nil {
		return nil, err
	}
	return actions, nil
}

func ValidatePermissionID(permissionId string) error {
	// Check if the resource ID exists in the database
	var count int64
	if err := config.DB.Model(&dto.TNTPermission{}).Where("permission_id = ? AND row_status = 1", permissionId).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func ValidateTenantID(tenantID uuid.UUID) error {
	// Check if the resource ID exists in the database

	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ? AND row_status = 1", "Tenant").First(&resourceType).Error; err != nil {
		return fmt.Errorf("resource type not found: %w", err)
	}
	var count int64
	if err := config.DB.Model(&dto.TenantResource{}).Where("resource_id = ? AND row_status = 1 AND resource_type_id = ?", tenantID, resourceType.ResourceTypeID).Count(&count).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if count == 0 {
		return errors.New("Tenant ID does not exist")
	}
	return nil
}

func ValidateMstResType(resourceID uuid.UUID) error {
	var count int64
	if err := config.DB.Model(&dto.Mst_ResourceTypes{}).
		Where("resource_type_id = ? AND row_status = 1", resourceID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func validateCreateRoleInput(input models.CreateRoleInput) error {

	// Validate input
	if input.Name == "" || input.AssignableScopeRef == uuid.Nil || input.RoleType == "" || input.RoleType == "DEFAULT" {
		return errors.New("invalid input recieved")
	}

	if err := utils.ValidateName(input.Name); err != nil {
		return fmt.Errorf("invalid role name: %w", err)
	}

	if err := ValidateMstResType(input.AssignableScopeRef); err != nil {
		return fmt.Errorf("invalid assignableScopeRef ID")
	}

	return nil
}

func prepareRoleObject(input models.CreateRoleInput) dto.TNTRole {
	role := dto.TNTRole{}
	role.Name = input.Name
	if input.Description != nil {
		role.Description = *input.Description
	}

	role.RoleType = dto.RoleTypeEnumCustom
	role.Version = input.Version
	role.CreatedAt = time.Now()
	role.ScopeResourceTypeID = input.AssignableScopeRef
	role.CreatedBy = constants.DefaltCreatedBy
	role.UpdatedBy = constants.DefaltCreatedBy
	role.RowStatus = 1
	return role
}

func (r *RoleMutationResolver) validateUpdateRoleInput(input models.UpdateRoleInput) error {
	extRole := dto.TNTRole{}
	if err := r.DB.Where("resource_id = ? AND row_status = 1", input.ID).First(&extRole).Error; err != nil {
		return errors.New("role not found")
	}

	if input.RoleType == "" {
		return errors.New("role type is required")
	} else if input.RoleType == "DEFAULT" {
		return errors.New("role type Default is not allowed")
	}

	if extRole.ScopeResourceTypeID != input.AssignableScopeRef {
		return errors.New("AssignableScopeRef cannot be changed")
	}
	return nil
}
