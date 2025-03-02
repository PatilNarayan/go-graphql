package roles

import (
	"context"
	"errors"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoleIDRequired    = errors.New("role ID is required")
	ErrInvalidUUIDFormat = errors.New("invalid UUID format")
	ErrInvalidRoleData   = errors.New("invalid role data")
	ErrInvalidPermission = errors.New("invalid permission data")
)

// RoleQueryResolver handles role-related queries.
type RoleQueryResolver struct {
	DB *gorm.DB
}

// Role retrieves a single role by ID.
func (r *RoleQueryResolver) Role(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	if id == uuid.Nil {
		return r.handleError("400", "Role ID is required", ErrRoleIDRequired)
	}

	// to get role we require role scope resource id without that we will not get able to create url endpoints
	role, err := r.getRoleFromDB(id)
	if err != nil {
		return r.handleError("400", "Role not found", err)
	}

	pc := permit.NewPermitClient()
	roleData, err := pc.SendRequest(ctx, "GET", fmt.Sprintf("resources/%s/roles/%s", role.ScopeResourceTypeID, id), nil)
	if err != nil {
		return r.handleError("400", "Error retrieving role from permit system", err)
	}

	mappedRole, err := r.mapToRole(roleData)
	if err != nil {
		return r.handleError("400", "Error mapping role data", err)
	}

	return utils.FormatSuccess([]models.Data{*mappedRole})
}

// Roles retrieves all roles.
func (r *RoleQueryResolver) Roles(ctx context.Context) (models.OperationResult, error) {
	pc := permit.NewPermitClient()
	data, err := pc.SendRequest(ctx, "GET", "resources?include_total_count=true", nil)
	if err != nil {
		return r.handleError("400", "Error retrieving roles from permit system", err)
	}

	roles, err := r.extractRolesFromData(data)
	if err != nil {
		return r.handleError("400", "Error extracting roles from data", err)
	}

	return utils.FormatSuccess(roles)
}

// Helper Functions

func (r *RoleQueryResolver) handleError(code, message string, err error) (models.OperationResult, error) {
	em := fmt.Sprintf("%s: %v", message, err)
	logger.LogError(em)
	return utils.FormatError(utils.FormatErrorStruct(code, message, em)), nil
}

func (r *RoleQueryResolver) getRoleFromDB(id uuid.UUID) (*dto.TNTRole, error) {
	var role dto.TNTRole
	err := r.DB.Where("resource_id = ? AND row_status = 1", id).First(&role).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRoleNotFound, err)
	}
	return &role, nil
}

func (r *RoleQueryResolver) extractRolesFromData(data map[string]interface{}) ([]models.Data, error) {
	var roles []models.Data

	for _, v := range data["data"].([]interface{}) {
		v := v.(map[string]interface{})
		if _, ok := v["roles"]; !ok {
			continue
		}

		for _, role := range v["roles"].(map[string]interface{}) {
			mappedRole, err := r.mapToRole(role.(map[string]interface{}))
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrInvalidRoleData, err)
			}
			roles = append(roles, *mappedRole)
		}
	}

	return roles, nil
}

func (r *RoleQueryResolver) mapToRole(roleData map[string]interface{}) (*models.Role, error) {
	attributes, ok := roleData["attributes"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: missing attributes", ErrInvalidRoleData)
	}

	role := &models.Role{}

	// Parse required fields
	if err := r.parseRoleFields(attributes, role); err != nil {
		return nil, err
	}

	// Parse optional fields
	if err := r.parseOptionalRoleFields(attributes, role); err != nil {
		return nil, err
	}

	// Parse permissions
	if err := r.parsePermissions(attributes, role); err != nil {
		return nil, err
	}

	// Parse assignable scope
	if err := r.parseAssignableScope(attributes, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (r *RoleQueryResolver) parseRoleFields(data map[string]interface{}, role *models.Role) error {
	// Parse ID
	idStr, ok := data["ID"].(string)
	if !ok {
		return fmt.Errorf("%w: missing or invalid ID", ErrInvalidRoleData)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
	}
	role.ID = id

	// Parse Name
	name, ok := data["Name"].(string)
	if !ok {
		return fmt.Errorf("%w: missing or invalid name", ErrInvalidRoleData)
	}
	role.Name = name

	// Parse CreatedAt
	createdAt, ok := data["createdAt"].(string)
	if !ok {
		return fmt.Errorf("%w: missing or invalid createdAt", ErrInvalidRoleData)
	}
	role.CreatedAt = createdAt

	// Parse UpdatedAt
	updatedAt, ok := data["updatedAt"].(string)
	if !ok {
		return fmt.Errorf("%w: missing or invalid updatedAt", ErrInvalidRoleData)
	}
	role.UpdatedAt = updatedAt

	// Parse CreatedBy
	createdByStr, ok := data["createdBy"].(string)
	if !ok {
		return fmt.Errorf("%w: missing or invalid createdBy", ErrInvalidRoleData)
	}
	createdBy, err := uuid.Parse(createdByStr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
	}
	role.CreatedBy = createdBy

	// Parse UpdatedBy
	updatedByStr, ok := data["updatedBy"].(string)
	if !ok {
		return fmt.Errorf("%w: missing or invalid updatedBy", ErrInvalidRoleData)
	}
	updatedBy, err := uuid.Parse(updatedByStr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
	}
	role.UpdatedBy = updatedBy

	// Parse RoleType
	roleType, ok := data["RoleType"].(string)
	if !ok {
		return fmt.Errorf("%w: missing or invalid roleType", ErrInvalidRoleData)
	}
	role.RoleType = models.RoleTypeEnum(roleType)

	// Parse Version
	version, ok := data["Version"].(string)
	if !ok {
		return fmt.Errorf("%w: missing or invalid version", ErrInvalidRoleData)
	}
	role.Version = version

	return nil
}

func (r *RoleQueryResolver) parseOptionalRoleFields(data map[string]interface{}, role *models.Role) error {
	// Parse Description (optional)
	if desc, ok := data["Description"].(string); ok {
		role.Description = &desc
	}
	return nil
}

func (r *RoleQueryResolver) parsePermissions(data map[string]interface{}, role *models.Role) error {
	perms, ok := data["Permissions"].([]interface{})
	if !ok {
		return nil // Permissions are optional
	}

	var permissions []*models.Permission
	for _, p := range perms {
		permMap, ok := p.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid permission format", ErrInvalidPermission)
		}

		permission, err := r.parsePermission(permMap)
		if err != nil {
			return err
		}
		permissions = append(permissions, permission)
	}

	role.Permissions = permissions
	return nil
}

func (r *RoleQueryResolver) parsePermission(data map[string]interface{}) (*models.Permission, error) {
	permission := &models.Permission{}

	// Parse ID
	idStr, ok := data["permissionId"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: missing or invalid permissionId", ErrInvalidPermission)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
	}
	permission.ID = id

	// Parse Name and Action
	if name, ok := data["name"].(string); ok {
		permission.Name = name
		permission.Action = name
	}

	// Parse AssignableScope
	if resourcetypeId, ok := data["resourcetypeId"].(string); ok {
		permission.AssignableScope = resourcetypeId
	}

	// Parse CreatedAt and UpdatedAt
	if createdAt, ok := data["createdAt"].(string); ok {
		permission.CreatedAt = createdAt
	}
	if updatedAt, ok := data["updatedAt"].(string); ok {
		permission.UpdatedAt = updatedAt
	}

	// Parse CreatedBy and UpdatedBy
	if createdByStr, ok := data["createdBy"].(string); ok {
		createdBy, err := uuid.Parse(createdByStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
		}
		permission.CreatedBy = createdBy
	}
	if updatedByStr, ok := data["updatedBy"].(string); ok {
		updatedBy, err := uuid.Parse(updatedByStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
		}
		permission.UpdatedBy = updatedBy
	}

	return permission, nil
}

func (r *RoleQueryResolver) parseAssignableScope(data map[string]interface{}, role *models.Role) error {
	scopeMap, ok := data["AssignableScopeRef"].(map[string]interface{})
	if !ok {
		return nil // Assignable scope is optional
	}

	assignableScope, err := r.parseRoot(scopeMap)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidAssignableScope, err)
	}
	role.AssignableScope = assignableScope
	return nil
}

func (r *RoleQueryResolver) parseRoot(data map[string]interface{}) (*models.Root, error) {
	root := &models.Root{}

	// Parse ID
	idStr, ok := data["resource_type_id"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: missing or invalid resource_type_id", ErrInvalidAssignableScope)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
	}
	root.ID = id

	// Parse Name
	if name, ok := data["name"].(string); ok {
		root.Name = name
	}

	// Parse CreatedAt and UpdatedAt
	if createdAt, ok := data["created_at"].(string); ok {
		root.CreatedAt = createdAt
	}
	if updatedAt, ok := data["updated_at"].(string); ok {
		root.UpdatedAt = updatedAt
	}

	// Parse CreatedBy and UpdatedBy
	if createdByStr, ok := data["created_by"].(string); ok {
		createdBy, err := uuid.Parse(createdByStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
		}
		root.CreatedBy = createdBy
	}
	if updatedByStr, ok := data["updated_by"].(string); ok {
		updatedBy, err := uuid.Parse(updatedByStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidUUIDFormat, err)
		}
		root.UpdatedBy = updatedBy
	}

	return root, nil
}
