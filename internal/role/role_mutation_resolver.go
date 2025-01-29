package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go_graphql/config"
	"go_graphql/gql/models"
	"go_graphql/internal/constants"
	"go_graphql/internal/dto"
	"go_graphql/internal/utils"
	"go_graphql/logger"
	"go_graphql/permit"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	DB *gorm.DB
}

// CreateRole handles creating a new role.
func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.CreateRoleInput) (*models.Role, error) {
	logger.Log.Info("Starting CreateRole")

	// Validate input
	if input.Name == "" {
		logger.Log.Warn("Role name is required")
		return nil, errors.New("role name is required")
	} else {
		if err := utils.ValidateName(input.Name); err != nil {
			logger.AddContext(err).Error("Invalid role name")
			return nil, fmt.Errorf("invalid role name: %w", err)
		}
	}

	if input.AssignableScopeRef == uuid.Nil {
		logger.Log.Warn("Resource ID is required")
		return nil, errors.New("resource ID is required")
	} else {
		if err := utils.ValidateMstResType(input.AssignableScopeRef); err != nil {
			logger.AddContext(err).Error("invalid resource ID")
			return nil, fmt.Errorf("invalid resource ID: %w", err)
		}
	}

	if input.RoleType == "" {
		logger.Log.Warn("Role type is required")
		return nil, errors.New("role type is required")
	} else if input.RoleType == "DEFAULT" {
		logger.Log.Warn("Role type Default is not allowed")
		return nil, errors.New("role type Default is not allowed")
	}

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	// check if role already exists
	var roleExists dto.TenantResource
	if err := r.DB.Where("name = ? AND row_status = 1 AND resource_type_id = ? AND parent_resource_id = ?", input.Name, resourceType.ResourceTypeID, input.AssignableScopeRef).First(&roleExists).Error; err == nil {
		logger.AddContext(err).Error("Role already exists")
		return nil, fmt.Errorf("role already exists")
	}

	if err := CheckPermissions(input.Permissions); err != nil {
		return nil, err
	}

	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where("resource_type_id = ? AND row_status = 1", input.AssignableScopeRef).First(&assignableScopeRef).Error; err != nil {
		logger.AddContext(err).Error("Invalid TenantID")
		return nil, fmt.Errorf("invalid TenantID")
	}

	permitMap := make(map[string]interface{})

	permitMap = map[string]interface{}{
		"name": input.Name,
		"key":  input.Name,
	}

	if input.Description != nil {
		permitMap["description"] = *input.Description
	}

	//create role in permit
	pc := permit.NewPermitClient()
	_, err := pc.APIExecute(ctx, "POST", fmt.Sprintf("resources/%s/roles", assignableScopeRef.Name), permitMap)

	if err != nil {
		return nil, err
	}

	resources, err := pc.APIExecute(ctx, "GET", "resources", nil)
	if err != nil {
		return nil, err
	}

	actions := utils.GetActionMap(resources.([]interface{}), assignableScopeRef.Name)
	permission, err := utils.GetPermissionAction(input.Permissions)
	if err != nil {
		return nil, err
	}
	res := utils.CreateActionMap(actions, permission)
	update := map[string]interface{}{
		"actions": res,
	}
	_, err = pc.APIExecute(ctx, "PATCH", fmt.Sprintf("resources/%s", assignableScopeRef.Name), update)
	if err != nil {
		return nil, err
	}

	role := dto.TNTRole{}
	role.Name = input.Name
	if input.Description != nil {
		role.Description = *input.Description
	}

	role.RoleType = dto.RoleTypeEnumCustom
	role.Version = input.Version
	role.CreatedAt = time.Now()
	role.ResourceTypeID = input.AssignableScopeRef
	role.CreatedBy = constants.DefaltCreatedBy
	role.UpdatedBy = constants.DefaltCreatedBy
	role.RowStatus = 1

	logger.Log.Info("Fetching resource type")

	logger.Log.Info("Creating tenant resource")
	tenantResource := &dto.TenantResource{
		ResourceID:     uuid.New(),
		ResourceTypeID: resourceType.ResourceTypeID,
		Name:           input.Name,
		RowStatus:      1,
		CreatedBy:      constants.DefaltCreatedBy,
		UpdatedBy:      constants.DefaltCreatedBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := r.DB.Create(&tenantResource).Error; err != nil {
		logger.AddContext(err).Error("Failed to create tenant resource")
		return nil, err
	}

	role.ResourceID = tenantResource.ResourceID
	logger.Log.Info("Saving role to database")
	if err := r.DB.Create(&role).Error; err != nil {
		logger.AddContext(err).Error("Failed to save role")
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
			logger.AddContext(tx.Error).Error("Failed to save role permission")
			return nil, tx.Error
		}
	}
	if len(input.Permissions) > 0 {
		err = SetPermission(ctx, role.Name, assignableScopeRef.Name, permission)
		if err != nil {
			return nil, err
		}
	}
	logger.Log.Info("Role created successfully")
	return convertRoleToGraphQL(&role), nil
}

// UpdateRole handles updating an existing role.
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, input models.UpdateRoleInput) (*models.Role, error) {
	logger.Log.Infof("Starting UpdateRole for ID: %s", input.ID)

	var role dto.TNTRole
	if err := utils.ValidateRole(input.ID); err != nil {
		return nil, err
	}
	var oldROleName dto.TNTRole
	if err := r.DB.Where("resource_id = ? AND row_status = 1", input.ID).First(&oldROleName).Error; err != nil {
		logger.AddContext(err).Warn("Role not found")
		return nil, errors.New("role not found")
	}

	oldROleName.Name = role.Name

	if input.RoleType == "" {
		logger.Log.Warn("Role type is required")
		return nil, errors.New("role type is required")
	} else if input.RoleType == "DEFAULT" {
		logger.Log.Warn("Role type Default is not allowed")
		return nil, errors.New("role type Default is not allowed")
	}

	// var roleResource dto.TenantResource

	// if err := r.DB.First(&roleResource, "resource_id = ? AND row_status = 1", input.ID).Error; err != nil {
	// 	logger.AddContext(err).Warn("Role not found")
	// 	return nil, errors.New("role not found")
	// }

	if *&oldROleName.ResourceTypeID != input.AssignableScopeRef {
		logger.Log.Warn("AssignableScopeRef cannot be changed")
		return nil, errors.New("AssignableScopeRef cannot be changed")
	}

	if oldROleName.Name != input.Name {
		logger.Log.Warn("Role name cannot be changed")
		return nil, errors.New("Role name cannot be changed")
	}

	// Update fields
	if input.Name != "" {
		role.Name = input.Name
	}

	// if input.AssignableScopeRef == uuid.Nil {
	// 	logger.Log.Warn("Resource type ID is required")
	// 	return nil, errors.New("Resource type ID is required")
	// } else {
	// 	if err := utils.ValidateMstResType(input.AssignableScopeRef); err != nil {
	// 		logger.AddContext(err).Error("Invalid Resource type ID")
	// 		return nil, fmt.Errorf("invalid Resource type ID: %w", err)
	// 	}
	// }

	if err := CheckPermissions(input.Permissions); err != nil {
		return nil, err
	}

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	mstresourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("resource_type_id = ? AND row_status = 1", input.AssignableScopeRef).First(&mstresourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	// // check if role already exists
	// var roleExists dto.TenantResource
	// if err := r.DB.Where("name = ? AND row_status = 1 AND resource_type_id = ? AND parent_resource_id = ?", input.Name, resourceType.ResourceTypeID, input.AssignableScopeRef).First(&roleExists).Error; err == nil {
	// 	logger.AddContext(err).Error("Role already exists")
	// 	return nil, fmt.Errorf("role already exists")
	// }

	permitMap := make(map[string]interface{})

	permitMap = map[string]interface{}{
		"name": input.Name,
	}

	if input.Description != nil {
		permitMap["description"] = *input.Description
	}

	//create role in permit
	pc := permit.NewPermitClient()
	var id interface{}
	if data, err := pc.APIExecute(ctx, "GET", "resources", nil); err != nil {
		return nil, err
	} else {
		data := data.([]interface{})
		for _, v := range data {
			v := v.(map[string]interface{})
			if mstresourceType.Name == v["key"].(string) {
				rolesMap := v["roles"].(map[string]interface{})
				if rolesMap[input.Name] == nil {
					return nil, fmt.Errorf("No such role exists")
				} else {
					permitrole := rolesMap[input.Name].(map[string]interface{})
					id = permitrole["id"]
				}
			}
		}
	}

	resource, err := utils.GetResourceTypeIDName(input.AssignableScopeRef)
	if err != nil {
		return nil, err
	}

	if _, err := pc.APIExecute(ctx, "PATCH", fmt.Sprintf("resources/%s/roles/%s", *resource, id), permitMap); err != nil {
		return nil, err
	}

	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where("resource_type_id = ? AND row_status = 1", input.AssignableScopeRef).First(&assignableScopeRef).Error; err != nil {
		logger.AddContext(err).Error("Invalid TenantID")
		return nil, fmt.Errorf("invalid TenantID")
	}

	resources, err := pc.APIExecute(ctx, "GET", "resources", nil)
	if err != nil {
		return nil, err
	}

	actions := utils.GetActionMap(resources.([]interface{}), assignableScopeRef.Name)
	permission, err := utils.GetPermissionAction(input.Permissions)
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating role in database for ID: %s", input.ID)
	updateData := map[string]interface{}{
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

	if err := r.DB.Model(&role).Updates(updateData).Error; err != nil {
		logger.AddContext(err).Error("Failed to update role")
		return nil, err
	}

	var pdata []dto.TNTRolePermission
	if err := r.DB.Where("role_id = ? AND row_status = 1", input.ID).Find(&pdata).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch role permissions")
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
				logger.AddContext(tx.Error).Error("Failed to save role permission")
				return nil, tx.Error
			}

		}
	}

	newPermissionsActions, err := utils.GetPermissionAction(newPermissions)
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

	removeActions, err := utils.GetPermissionAction(removeIDsList)
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

	_, err = pc.APIExecute(ctx, "PATCH", fmt.Sprintf("resources/%s", assignableScopeRef.Name), update)
	if err != nil {
		return nil, err
	}
	for _, id := range removeIDs {
		if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_permission_id = ? AND row_status = 1", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
			logger.AddContext(err).Error("Failed to delete role")
			return nil, fmt.Errorf("failed to delete role: %w", err)
		}
	}
	if len(input.Permissions) > 0 {
		err = SetPermission(ctx, role.Name, assignableScopeRef.Name, newPermissions)
		if err != nil {
			return nil, err
		}
	}

	logger.Log.Infof("Role updated successfully for ID: %s", input.ID)
	return convertRoleToGraphQL(&role), nil
}

// DeleteRole handles deleting a role by ID.
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, id uuid.UUID) (bool, error) {
	logger.Log.Infof("Starting DeleteRole for ID: %s", id)

	if err := utils.ValidateMstResType(id); err != nil {
		logger.AddContext(err).Error("Invalid ID")
		return false, fmt.Errorf("invalid ID: %w", err)
	}

	var roleDB dto.TNTRole
	if err := r.DB.First(&roleDB, "resource_id = ? AND row_status = 1", id).Error; err != nil {
		logger.AddContext(err).Warn("Role not found")
		return false, errors.New("role not found")
	}

	updates := map[string]interface{}{
		// "deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}
	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where("resource_type_id = ? AND row_status = 1", roleDB.ResourceTypeID).First(&assignableScopeRef).Error; err != nil {
		logger.AddContext(err).Error("Invalid TenantID")
		return false, fmt.Errorf("invalid TenantID")
	}

	pc := permit.NewPermitClient()
	var pid interface{}
	if data, err := pc.APIExecute(ctx, "GET", fmt.Sprintf("resources/%s/roles", assignableScopeRef.Name), nil); err != nil {
		return false, err
	} else {
		data := data.([]interface{})
		for _, v := range data {
			v := v.(map[string]interface{})
			if v["name"] == roleDB.Name {
				pid = v["id"]
				break
			}
		}
	}

	if _, err := pc.APIExecute(ctx, "DELETE", fmt.Sprintf("resources/%s/roles/%s", assignableScopeRef.Name, pid), nil); err != nil {
		return false, err
	}

	logger.Log.Infof("Marking role as deleted for ID: %s", id)
	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ? AND row_status = 1", id).Updates(updates).Error; err != nil {
		logger.AddContext(err).Error("Failed to delete role")
		return false, err
	}

	if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_id = ? AND row_status = 1", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		logger.AddContext(err).Error("Failed to delete role")
		return false, fmt.Errorf("failed to delete role: %w", err)
	}

	if err := r.DB.Model(&dto.TenantResource{}).Where("resource_id = ? AND row_status = 1", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		logger.AddContext(err).Error("Failed to delete role")
		return false, fmt.Errorf("failed to delete role: %w", err)
	}

	actions, err := utils.GetPermissionResourceAction(*&roleDB.ResourceTypeID)
	if err != nil {
		return false, err
	}

	updateres := make(map[string]interface{})
	res := utils.CreateActionMap(updateres, actions)
	update := map[string]interface{}{
		"actions": res,
	}
	_, err = pc.APIExecute(ctx, "PATCH", fmt.Sprintf("resources/%s", assignableScopeRef.Name), update)
	if err != nil {
		return false, err
	}

	logger.Log.Infof("Role deleted successfully for ID: %s", id)
	return true, nil
}

func CheckPermissions(permissionsIds []string) error {
	//validate permissionIds
	for _, permissionID := range permissionsIds {
		if err := utils.ValidatePermissionID(permissionID); err != nil {
			logger.AddContext(err).Error("Invalid permission ID")
			return fmt.Errorf("invalid permission ID: %w", err)
		}
	}
	return nil
}

func CreateMstRole(tenantID uuid.UUID) error {
	var mstRole []*dto.MstRole

	if err := config.DB.Find(&mstRole).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch roles from the database")
		return err
	}

	var mstRolePermissions []*dto.MstRolePermission
	if err := config.DB.Find(&mstRolePermissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch role permissions from the database")
		return err
	}
	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
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
			logger.AddContext(err).Error("Failed to save role")
			return err
		}

		role := dto.TNTRole{
			ResourceID:     tenantResource.ResourceID,
			Name:           mrole.Name,
			Description:    mrole.Description,
			RoleType:       dto.RoleTypeEnumDefault,
			ResourceTypeID: resourceType.ResourceTypeID,
			Version:        mrole.Version,
			CreatedAt:      mrole.CreatedAt,
			CreatedBy:      mrole.CreatedBy,
			UpdatedBy:      mrole.UpdatedBy,
			UpdatedAt:      mrole.UpdatedAt,
			RowStatus:      mrole.RowStatus,
		}

		if err := config.DB.Save(&role).Error; err != nil {
			logger.AddContext(err).Error("Failed to save role")
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
				logger.AddContext(tx.Error).Error("Failed to save role permission")
				return tx.Error
			}
		}
	}

	return nil

}

func AddDefaultPermissions() (map[uuid.UUID]uuid.UUID, error) {
	var mstPermissions []*dto.MstPermission

	if err := config.DB.Where("row_status = 1").Find(&mstPermissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch permissions from the database")
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
			logger.AddContext(err).Error("Failed to create permission")
			return nil, err
		}
	}

	return res, nil
}

func DeleteDefaultRole(tenantID uuid.UUID) error {
	// Find resource type for Role
	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
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
		logger.AddContext(err).Error("Failed to fetch role resources")
		return err
	}

	// Update row status in tenant_resources table
	if err := tx.Model(&dto.TenantResource{}).
		Where("tenant_id = ? AND resource_type_id = ? AND row_status = 1", tenantID, resourceType.ResourceTypeID).
		Update("row_status", 0).Error; err != nil {
		tx.Rollback()
		logger.AddContext(err).Error("Failed to update tenant resources")
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
		logger.AddContext(err).Error("Failed to update roles")
		return err
	}

	// Update row status in tnt_role_permissions table
	if err := tx.Model(&dto.TNTRolePermission{}).
		Where("role_id IN ? AND row_status = 1", roleIDs).
		Update("row_status", 0).Error; err != nil {
		tx.Rollback()
		logger.AddContext(err).Error("Failed to update role permissions")
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.AddContext(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}

func SetPermission(ctx context.Context, roleName, assinableScopeName string, permissionAction []string) error {
	pc := permit.NewPermitClient()

	update := map[string]interface{}{
		"permissions": permissionAction,
	}

	_, err := pc.APIExecute(ctx, "POST", fmt.Sprintf("resources/%s/roles/%s/permissions", assinableScopeName, roleName), update)
	if err != nil {
		return err
	}

	return nil
}
