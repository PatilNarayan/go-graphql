package roles

import (
	"context"
	"errors"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	DB *gorm.DB
}

func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.CreateRoleInput) (models.OperationResult, error) {
	// Extract gin.Context from GraphQL context
	//ginCtx, ok := ctx.Value(middlewares.GinContextKey).(*gin.Context)
	// if !ok {
	// 	return nil, fmt.Errorf("unable to get gin context")
	// }
	//UserID := ginCtx.MustGet("userID").(string)
	//userUUID := uuid.MustParse(UserID)
	userUUID := uuid.New()

	// tenantID, err := helpers.GetTenant(ginCtx)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := ValidateTenantID(uuid.MustParse(tenantID)); err != nil {
	// 	return nil, err
	// }

	if err := validateCreateRoleInput(input); err != nil {
		return nil, err
	}

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where(&dto.Mst_ResourceTypes{Name: "Role", RowStatus: 1}).First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}
	// tenantIDUUID, err := uuid.Parse(tenantID)
	// if err != nil {
	// 	return nil, err
	// }
	// check if role already exists
	var roleExists dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{Name: input.Name, RowStatus: 1, ResourceTypeID: resourceType.ResourceTypeID, TenantID: &input.AssignableScopeRef}).First(&roleExists).Error; err == nil {
		return nil, fmt.Errorf("role already exists")
	}

	if err := CheckPermissions(input.Permissions); err != nil {
		return nil, err
	}

	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where(&dto.Mst_ResourceTypes{ResourceTypeID: input.AssignableScopeRef, RowStatus: 1}).First(&assignableScopeRef).Error; err != nil {
		return nil, fmt.Errorf("invalid TenantID")
	}

	permissionActions, permissionData, err := GetPermissionAction(input.AssignableScopeRef.String(), input.Permissions)
	if err != nil {
		return nil, err
	}
	inputMap := helpers.StructToMap(input)
	inputMap["actions"] = permissionActions
	inputMap["AssignableScopeRef"] = assignableScopeRef
	inputMap["Permissions"] = permissionData
	inputMap["createdBy"] = userUUID.String()
	inputMap["updatedBy"] = userUUID.String()
	inputMap["createdAt"] = time.Now().Format(time.RFC3339)
	inputMap["updatedAt"] = time.Now().Format(time.RFC3339)

	permitMap := map[string]interface{}{
		"name":        input.Name,
		"key":         input.ID,
		"attributes":  inputMap,
		"permissions": permissionActions,
	}

	if input.Description != nil {
		permitMap["description"] = *input.Description
	}

	//create role in permit
	pc := permit.NewPermitClient()
	_, err = pc.SendRequest(ctx, "POST", fmt.Sprintf("resources/%s/roles", assignableScopeRef.ResourceTypeID), permitMap)

	if err != nil {
		return nil, err
	}

	if err := r.DB.Create(&dto.TenantResource{
		ResourceID:     input.ID,
		ResourceTypeID: resourceType.ResourceTypeID,
		Name:           input.Name,
		RowStatus:      1,
		TenantID:       &input.AssignableScopeRef,
		CreatedBy:      userUUID,
		UpdatedBy:      userUUID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}).Error; err != nil {
		return nil, err
	}
	role := prepareRoleObject(input, userUUID)

	role.ResourceID = input.ID
	if err := r.DB.Create(&role).Error; err != nil {
		return nil, err
	}

	for _, permissionID := range input.Permissions {
		tx := r.DB.Create(&dto.TNTRolePermission{
			ID:           uuid.New(),
			RoleID:       role.ResourceID,
			PermissionID: uuid.MustParse(permissionID),
			RowStatus:    1,
			CreatedBy:    userUUID,
			UpdatedBy:    userUUID,
		})

		if tx.Error != nil {
			return nil, tx.Error
		}
	}

	RoleQueryResolver := &RoleQueryResolver{DB: r.DB}
	data, err := RoleQueryResolver.Role(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RoleMutationResolver) UpdateRole(ctx context.Context, input models.UpdateRoleInput) (models.OperationResult, error) {
	// Extract gin.Context from GraphQL context
	//ginCtx, ok := ctx.Value(middlewares.GinContextKey).(*gin.Context)
	// if !ok {
	// 	return nil, fmt.Errorf("unable to get gin context")
	// }
	//UserID := ginCtx.MustGet("userID").(string)
	//userUUID := uuid.MustParse(UserID)
	// userUUID := uuid.New()

	// tenantID, err := helpers.GetTenant(ginCtx)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := ValidateTenantID(uuid.MustParse(tenantID)); err != nil {
	// 	return nil, err
	// }

	// if err := ValidateTenantID(uuid.MustParse(tenantID)); err != nil {
	// 	return nil, err
	// }

	var role dto.TNTRole
	if err := r.DB.Where("resource_id = ? AND row_status = 1", input.ID).First(&role).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	if err := r.validateUpdateRoleInput(input); err != nil {
		return nil, err
	}

	// Update fields
	if input.Name != "" {
		role.Name = input.Name
	}

	var assignableScopeRefData dto.Mst_ResourceTypes
	if err := r.DB.Where(&dto.Mst_ResourceTypes{ResourceTypeID: input.AssignableScopeRef, RowStatus: 1}).First(&assignableScopeRefData).Error; err != nil {
		return nil, fmt.Errorf("invalid TenantID")
	}

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	permissionActions, permissionData, err := GetPermissionAction(input.AssignableScopeRef.String(), input.Permissions)
	if err != nil {
		return nil, err
	}
	inputMap := helpers.StructToMap(input)
	inputMap["Permissions"] = permissionData
	inputMap["AssignableScopeRef"] = assignableScopeRefData
	inputMap["actions"] = permissionActions
	inputMap["created_by"] = role.CreatedBy
	inputMap["updated_by"] = role.UpdatedBy
	inputMap["created_at"] = time.Now().Format(time.RFC3339)
	inputMap["updated_at"] = time.Now().Format(time.RFC3339)

	permitMap := map[string]interface{}{
		"name":        input.Name,
		"attributes":  inputMap,
		"permissions": permissionActions,
	}

	if input.Description != nil {
		permitMap["description"] = *input.Description
	}
	//create role in permit
	pc := permit.NewPermitClient()

	if _, err := pc.SendRequest(ctx, "PATCH", fmt.Sprintf("resources/%s/roles/%s", input.AssignableScopeRef, input.ID.String()), permitMap); err != nil {
		return nil, err
	}

	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where("resource_type_id = ? AND row_status = 1", input.AssignableScopeRef).First(&assignableScopeRef).Error; err != nil {
		return nil, fmt.Errorf("invalid TenantID")
	}

	updateData := map[string]interface{}{
		"name":       input.Name,
		"version":    input.Version,
		"role_type":  dto.RoleTypeEnumCustom,
		"updated_by": role.UpdatedBy,
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
				CreatedBy:    role.CreatedBy,
				UpdatedBy:    role.UpdatedBy,
			})

			if tx.Error != nil {
				return nil, tx.Error
			}

		}
	}

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

	for _, id := range removeIDs {
		if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_permission_id = ? AND row_status = 1", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
			return nil, fmt.Errorf("failed to delete role: %w", err)
		}
	}

	RoleQueryResolver := &RoleQueryResolver{DB: r.DB}
	data, err := RoleQueryResolver.Role(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RoleMutationResolver) DeleteRole(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	// Extract gin.Context from GraphQL context
	//ginCtx, ok := ctx.Value(middlewares.GinContextKey).(*gin.Context)
	// if !ok {
	// 	return nil, fmt.Errorf("unable to get gin context")
	// }
	//UserID := ginCtx.MustGet("userID").(string)
	//userUUID := uuid.MustParse(UserID)
	// userUUID := uuid.New()

	// tenantID, err := helpers.GetTenant(ginCtx)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := ValidateTenantID(uuid.MustParse(tenantID)); err != nil {
	// 	return nil, err
	// }

	// if err := ValidateTenantID(uuid.MustParse(tenantID)); err != nil {
	// 	return nil, err
	// }

	var roleDB dto.TNTRole
	if err := r.DB.First(&roleDB, "resource_id = ? AND row_status = 1", input.ID).Error; err != nil {
		return nil, errors.New("role not found")
	}

	var assignableScopeRef dto.Mst_ResourceTypes
	if err := r.DB.Where("resource_type_id = ? AND row_status = 1", roleDB.ScopeResourceTypeID).First(&assignableScopeRef).Error; err != nil {
		return nil, fmt.Errorf("invalid TenantID")
	}

	pc := permit.NewPermitClient()

	if _, err := pc.SendRequest(ctx, "DELETE", fmt.Sprintf("resources/%s/roles/%s", assignableScopeRef.ResourceTypeID, roleDB.ResourceID.String()), nil); err != nil {
		return nil, err
	}

	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ? AND row_status = 1", input.ID).UpdateColumns(utils.UpdateDeletedMap()).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_id = ? AND row_status = 1", input.ID).UpdateColumns(utils.UpdateDeletedMap()).Error; err != nil {
		return nil, fmt.Errorf("failed to delete role: %w", err)
	}

	RoleQueryResolver := &RoleQueryResolver{DB: r.DB}
	data, err := RoleQueryResolver.Role(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return data, nil
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

func CheckPermissions(permissionsIds []string) error {
	//validate permissionIds
	for _, permissionID := range permissionsIds {
		if err := ValidatePermissionID(permissionID); err != nil {
			return fmt.Errorf("invalid permission ID: %w", err)
		}
	}
	return nil
}

func ValidatePermissionID(permissionId string) error {
	// Check if the resource ID exists in the database
	var count int64
	if err := config.DB.Model(&dto.MstPermission{}).Where("permission_id = ? AND row_status = 1", permissionId).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func GetPermissionAction(resourceID string, permissionsIds []string) ([]string, []dto.MstPermission, error) {
	var actions []string
	//validate permissionIds
	var res []dto.MstPermission
	for _, permissionID := range permissionsIds {
		var data dto.MstPermission
		if err := config.DB.Model(&dto.MstPermission{}).Where("permission_id = ? AND row_status = 1", permissionID).First(&data).Error; err != nil {
			return nil, nil, err
		}
		res = append(res, data)
		actions = append(actions, data.Action)
	}

	pc := permit.NewPermitClient()

	resourceData, err := pc.SendRequest(context.Background(), "GET", fmt.Sprintf("resources/%s", resourceID), nil)
	if err != nil {
		return nil, nil, err
	}

	actionsData := resourceData["actions"].(map[string]interface{})
	for _, val := range actions {
		valied := false
		for key := range actionsData {
			fmt.Println(key, val)
			if key == val {
				valied = true
				break
			}
		}

		if !valied {
			return nil, nil, errors.New("invalid permission action")
		}
	}

	return actions, res, nil
}

func prepareRoleObject(input models.CreateRoleInput, userID uuid.UUID) dto.TNTRole {
	role := dto.TNTRole{}
	role.Name = input.Name
	if input.Description != nil {
		role.Description = *input.Description
	}

	role.RoleType = dto.RoleTypeEnumCustom
	role.Version = input.Version
	role.CreatedAt = time.Now()
	role.ScopeResourceTypeID = input.AssignableScopeRef
	role.CreatedBy = userID
	role.UpdatedBy = userID
	role.RowStatus = 1
	return role
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
