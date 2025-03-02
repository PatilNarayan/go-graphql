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
	"iam_services_main_v1/pkg/logger"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoleAlreadyExists      = errors.New("role already exists")
	ErrInvalidRoleType        = errors.New("invalid role type")
	ErrInvalidAssignableScope = errors.New("invalid assignable scope")
	ErrInvalidPermissions     = errors.New("invalid permissions")
	ErrRoleNotFound           = errors.New("role not found")
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	DB *gorm.DB
}

// CreateRole creates a new role.
func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.CreateRoleInput) (models.OperationResult, error) {
	userUUID := uuid.New() // Replace with actual user ID from context

	tenantID, err := helpers.GetTenantID(ctx)
	if err != nil {
		return r.handleError("404", "Invalid tenant ID", err)
	}

	// resourceType, err := r.getResourceTypeByName("Role")
	// if err != nil {
	// 	return r.handleError("404", "Resource type not found", err)
	// }

	// if err := r.checkRoleExists(input.Name, resourceType.ResourceTypeID, tenantID); err != nil {
	// 	return r.handleError("409", "Role already exists", err)
	// }

	// if err := validatePermissions(input.Permissions); err != nil {
	// 	return r.handleError("400", "Invalid permissions", err)
	// }

	assignableScopeRef, err := r.getAssignableScopeRef(input.AssignableScopeRef)
	if err != nil {
		return r.handleError("404", "Assignable scope ref not found", err)
	}

	// you cant remove this DB operations because in graphql u r passing the permissions ids and in
	// permit this ids are not there its only actions are present in permit permissions or
	// need to add permissions get apis but this ids are need to present in DB also
	permissionActions, permissionData, err := r.getPermissionActions(input.AssignableScopeRef.String(), input.Permissions)
	if err != nil {
		return r.handleError("400", "Invalid permissions", err)
	}

	inputMap := r.prepareInputMap(input, permissionActions, permissionData, assignableScopeRef, tenantID, userUUID)
	permitMap := r.preparePermitMap(input, inputMap, permissionActions)

	pc := permit.NewPermitClient()
	if _, err := pc.SendRequest(ctx, "POST", fmt.Sprintf("resources/%s/roles", input.AssignableScopeRef.String()), permitMap); err != nil {
		return r.handleError("500", "Error creating role in permit", err)
	}

	if err := r.createTenantResource(input, input.AssignableScopeRef, tenantID, userUUID); err != nil {
		return r.handleError("500", "Error creating tenant resource", err)
	}

	role := r.prepareRoleObject(input, userUUID)
	if err := r.DB.Create(&role).Error; err != nil {
		return r.handleError("500", "Error creating role", err)
	}

	if err := r.createRolePermissions(input.ID, input.Permissions, userUUID); err != nil {
		return r.handleError("500", "Error creating role permissions", err)
	}

	roleQueryResolver := &RoleQueryResolver{DB: r.DB}
	return roleQueryResolver.Role(ctx, input.ID)
}

// UpdateRole updates an existing role.
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, input models.UpdateRoleInput) (models.OperationResult, error) {
	role, err := r.getRoleByID(input.ID)
	if err != nil {
		return r.handleError("500", "Error getting role", err)
	}

	if err := r.validateUpdateRoleInput(input, role); err != nil {
		return r.handleError("400", "Invalid input", err)
	}

	assignableScopeRef, err := r.getAssignableScopeRef(input.AssignableScopeRef)
	if err != nil {
		return r.handleError("404", "Assignable scope ref not found", err)
	}

	permissionActions, permissionData, err := r.getPermissionActions(input.AssignableScopeRef.String(), input.Permissions)
	if err != nil {
		return r.handleError("400", "Invalid permissions", err)
	}

	inputMap := r.prepareInputMapForUpdate(input, permissionActions, permissionData, assignableScopeRef, role)
	permitMap := r.preparePermitMapForUpdate(input, inputMap, permissionActions)

	pc := permit.NewPermitClient()
	if _, err := pc.SendRequest(ctx, "PATCH", fmt.Sprintf("resources/%s/roles/%s", input.AssignableScopeRef, input.ID.String()), permitMap); err != nil {
		return r.handleError("500", "Error updating role in permit", err)
	}

	if err := r.updateRoleDetails(role, input); err != nil {
		return r.handleError("500", "Error updating role", err)
	}

	if err := r.updateRolePermissions(input.ID, input.Permissions, role.CreatedBy, role.UpdatedBy); err != nil {
		return r.handleError("500", "Error updating role permissions", err)
	}

	roleQueryResolver := &RoleQueryResolver{DB: r.DB}
	return roleQueryResolver.Role(ctx, input.ID)
}

// DeleteRole deletes a role.
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	role, err := r.getRoleByID(input.ID)
	if err != nil {
		return r.handleError("500", "Error getting role", err)
	}

	assignableScopeRef, err := r.getAssignableScopeRef(role.ScopeResourceTypeID)
	if err != nil {
		return r.handleError("404", "Assignable scope ref not found", err)
	}

	pc := permit.NewPermitClient()
	if _, err := pc.SendRequest(ctx, "DELETE", fmt.Sprintf("resources/%s/roles/%s", assignableScopeRef.ResourceTypeID, role.ResourceID.String()), nil); err != nil {
		return r.handleError("500", "Error deleting role in permit", err)
	}

	if err := r.deleteRoleResources(input.ID); err != nil {
		return r.handleError("500", "Error deleting role resources", err)
	}

	return utils.FormatSuccess([]models.Data{})
}

// Helper Functions

func (r *RoleMutationResolver) handleError(code, message string, err error) (models.OperationResult, error) {
	em := fmt.Sprintf("%s: %v", message, err)
	logger.LogError(em)
	return utils.FormatError(utils.FormatErrorStruct(code, message, em)), nil
}

func (r *RoleMutationResolver) getResourceTypeByName(name string) (*dto.Mst_ResourceTypes, error) {
	var resourceType dto.Mst_ResourceTypes
	err := r.DB.Where("name = ? AND row_status = 1", name).First(&resourceType).Error
	if err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}
	return &resourceType, nil
}

func (r *RoleMutationResolver) checkRoleExists(name string, resourceTypeID int, tenantID *uuid.UUID) error {
	var roleExists dto.TenantResource
	err := r.DB.Where("name = ? AND resource_type_id = ? AND tenant_id = ? AND row_status = 1", name, resourceTypeID, tenantID).First(&roleExists).Error
	if err == nil {
		return ErrRoleAlreadyExists
	}
	return nil
}

func (r *RoleMutationResolver) getAssignableScopeRef(scopeRef uuid.UUID) (*dto.Mst_ResourceTypes, error) {
	var assignableScopeRef dto.Mst_ResourceTypes
	err := r.DB.Where("resource_type_id = ? AND row_status = 1", scopeRef).First(&assignableScopeRef).Error
	if err != nil {
		return nil, fmt.Errorf("assignable scope ref not found: %w", err)
	}
	return &assignableScopeRef, nil
}

func (r *RoleMutationResolver) getPermissionActions(resourceID string, permissions []string) ([]string, []dto.MstPermission, error) {
	var actions []string
	var permissionsData []dto.MstPermission

	for _, permissionID := range permissions {
		var permission dto.MstPermission
		err := config.DB.Where("permission_id = ? AND row_status = 1", permissionID).First(&permission).Error
		if err != nil {
			return nil, nil, fmt.Errorf("invalid permission ID: %w", err)
		}
		permissionsData = append(permissionsData, permission)
		actions = append(actions, permission.Name)
	}

	pc := permit.NewPermitClient()
	resourceData, err := pc.SendRequest(context.Background(), "GET", fmt.Sprintf("resources/%s", resourceID), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch resource data: %w", err)
	}

	actionsData := resourceData["actions"].(map[string]interface{})
	for _, action := range actions {
		if _, ok := actionsData[action]; !ok {
			return nil, nil, fmt.Errorf("invalid permission action: %s", action)
		}
	}

	return actions, permissionsData, nil
}

func (r *RoleMutationResolver) prepareInputMap(input models.CreateRoleInput, actions []string, permissions []dto.MstPermission, assignableScopeRef *dto.Mst_ResourceTypes, tenantID *uuid.UUID, userID uuid.UUID) map[string]interface{} {
	inputMap := helpers.StructToMap(input)
	inputMap["actions"] = actions
	inputMap["AssignableScopeRef"] = assignableScopeRef
	inputMap["TenantID"] = tenantID
	inputMap["Permissions"] = permissions
	inputMap["createdBy"] = userID.String()
	inputMap["updatedBy"] = userID.String()
	inputMap["createdAt"] = time.Now().Format(time.RFC3339)
	inputMap["updatedAt"] = time.Now().Format(time.RFC3339)
	return inputMap
}

func (r *RoleMutationResolver) preparePermitMap(input models.CreateRoleInput, inputMap map[string]interface{}, actions []string) map[string]interface{} {
	permitMap := map[string]interface{}{
		"name":        input.Name,
		"key":         input.ID,
		"attributes":  inputMap,
		"permissions": actions,
	}
	if input.Description != nil {
		permitMap["description"] = *input.Description
	}
	return permitMap
}

func (r *RoleMutationResolver) createTenantResource(input models.CreateRoleInput, resourceTypeID uuid.UUID, tenantID *uuid.UUID, userID uuid.UUID) error {
	return r.DB.Create(&dto.TenantResource{
		ResourceID:     input.ID,
		ResourceTypeID: resourceTypeID,
		Name:           input.Name,
		RowStatus:      1,
		TenantID:       tenantID,
		CreatedBy:      userID,
		UpdatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}).Error
}

func (r *RoleMutationResolver) prepareRoleObject(input models.CreateRoleInput, userID uuid.UUID) dto.TNTRole {
	return dto.TNTRole{
		ResourceID:          input.ID,
		Name:                input.Name,
		Description:         *input.Description,
		RoleType:            dto.RoleTypeEnumCustom,
		Version:             input.Version,
		ScopeResourceTypeID: input.AssignableScopeRef,
		CreatedBy:           userID,
		UpdatedBy:           userID,
		RowStatus:           1,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

func (r *RoleMutationResolver) createRolePermissions(roleID uuid.UUID, permissions []string, userID uuid.UUID) error {
	for _, permissionID := range permissions {
		if err := r.DB.Create(&dto.TNTRolePermission{
			ID:           uuid.New(),
			RoleID:       roleID,
			PermissionID: uuid.MustParse(permissionID),
			RowStatus:    1,
			CreatedBy:    userID,
			UpdatedBy:    userID,
		}).Error; err != nil {
			return fmt.Errorf("failed to create role permission: %w", err)
		}
	}
	return nil
}

func (r *RoleMutationResolver) getRoleByID(roleID uuid.UUID) (*dto.TNTRole, error) {
	var role dto.TNTRole
	err := r.DB.Where("resource_id = ? AND row_status = 1", roleID).First(&role).Error
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

func (r *RoleMutationResolver) validateUpdateRoleInput(input models.UpdateRoleInput, role *dto.TNTRole) error {
	if input.RoleType == "" {
		return ErrInvalidRoleType
	}
	if input.RoleType == "DEFAULT" {
		return ErrInvalidRoleType
	}
	if role.ScopeResourceTypeID != input.AssignableScopeRef {
		return ErrInvalidAssignableScope
	}
	return nil
}

func (r *RoleMutationResolver) prepareInputMapForUpdate(input models.UpdateRoleInput, actions []string, permissions []dto.MstPermission, assignableScopeRef *dto.Mst_ResourceTypes, role *dto.TNTRole) map[string]interface{} {
	inputMap := helpers.StructToMap(input)
	inputMap["actions"] = actions
	inputMap["AssignableScopeRef"] = assignableScopeRef
	inputMap["Permissions"] = permissions
	inputMap["createdBy"] = role.CreatedBy.String()
	inputMap["updatedBy"] = role.UpdatedBy.String()
	inputMap["createdAt"] = role.CreatedAt.Format(time.RFC3339)
	inputMap["updatedAt"] = time.Now().Format(time.RFC3339)
	return inputMap
}

func (r *RoleMutationResolver) preparePermitMapForUpdate(input models.UpdateRoleInput, inputMap map[string]interface{}, actions []string) map[string]interface{} {
	permitMap := map[string]interface{}{
		"name":        input.Name,
		"attributes":  inputMap,
		"permissions": actions,
	}
	if input.Description != nil {
		permitMap["description"] = *input.Description
	}
	return permitMap
}

func (r *RoleMutationResolver) updateRoleDetails(role *dto.TNTRole, input models.UpdateRoleInput) error {
	updates := map[string]interface{}{
		"name":       input.Name,
		"version":    input.Version,
		"updated_by": role.UpdatedBy,
		"updated_at": time.Now(),
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	return r.DB.Model(role).Where("resource_id = ? AND row_status = 1", input.ID).Updates(updates).Error
}

func (r *RoleMutationResolver) updateRolePermissions(roleID uuid.UUID, permissions []string, createdBy, updatedBy uuid.UUID) error {
	var existingPermissions []dto.TNTRolePermission
	if err := r.DB.Where("role_id = ? AND row_status = 1", roleID).Find(&existingPermissions).Error; err != nil {
		return fmt.Errorf("failed to fetch existing permissions: %w", err)
	}

	// Add new permissions
	for _, permissionID := range permissions {
		exists := false
		for _, p := range existingPermissions {
			if p.PermissionID.String() == permissionID {
				exists = true
				break
			}
		}
		if !exists {
			if err := r.DB.Create(&dto.TNTRolePermission{
				ID:           uuid.New(),
				RoleID:       roleID,
				PermissionID: uuid.MustParse(permissionID),
				RowStatus:    1,
				CreatedBy:    createdBy,
				UpdatedBy:    updatedBy,
			}).Error; err != nil {
				return fmt.Errorf("failed to create role permission: %w", err)
			}
		}
	}

	// Remove old permissions
	for _, p := range existingPermissions {
		exists := false
		for _, permissionID := range permissions {
			if p.PermissionID.String() == permissionID {
				exists = true
				break
			}
		}
		if !exists {
			if err := r.DB.Model(&dto.TNTRolePermission{}).Where("id = ?", p.ID).Updates(utils.UpdateDeletedMap()).Error; err != nil {
				return fmt.Errorf("failed to delete role permission: %w", err)
			}
		}
	}

	return nil
}

func (r *RoleMutationResolver) deleteRoleResources(roleID uuid.UUID) error {
	if err := r.DB.Model(&dto.TenantResource{}).Where("resource_id = ?", roleID).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		return fmt.Errorf("failed to delete tenant resource: %w", err)
	}
	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ?", roleID).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_id = ?", roleID).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		return fmt.Errorf("failed to delete role permissions: %w", err)
	}
	return nil
}

func validateCreateRoleInput(input models.CreateRoleInput) error {
	if input.Name == "" || input.AssignableScopeRef == uuid.Nil || input.RoleType == "" || input.RoleType == "DEFAULT" {
		return ErrInvalidRoleType
	}
	if err := utils.ValidateName(input.Name); err != nil {
		return fmt.Errorf("invalid role name: %w", err)
	}
	if err := ValidateMstResType(input.AssignableScopeRef); err != nil {
		return fmt.Errorf("invalid assignable scope ref: %w", err)
	}
	return nil
}

func validatePermissions(permissions []string) error {
	for _, permissionID := range permissions {
		if err := ValidatePermissionID(permissionID); err != nil {
			return fmt.Errorf("invalid permission ID: %w", err)
		}
	}
	return nil
}

func ValidateMstResType(resourceID uuid.UUID) error {
	var count int64
	if err := config.DB.Model(&dto.Mst_ResourceTypes{}).Where("resource_type_id = ? AND row_status = 1", resourceID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate resource type: %w", err)
	}
	if count == 0 {
		return ErrInvalidAssignableScope
	}
	return nil
}

func ValidatePermissionID(permissionID string) error {
	var count int64
	if err := config.DB.Model(&dto.MstPermission{}).Where("permission_id = ? AND row_status = 1", permissionID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate permission ID: %w", err)
	}
	if count == 0 {
		return ErrInvalidPermissions
	}
	return nil
}
