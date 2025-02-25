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

// RoleQueryResolver handles role-related queries.
type RoleQueryResolver struct {
	DB *gorm.DB
}

func (r *RoleQueryResolver) Role(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {

	if id == uuid.Nil {
		em := fmt.Sprint("Role ID is required")
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "Role ID is required", em)), nil
	}

	var role dto.TNTRole
	if err := r.DB.Where("resource_id = ? AND row_status = 1", id).First(&role).Error; err != nil {
		em := fmt.Sprintf("role not found: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "role not found", em)), nil
	}

	pc := permit.NewPermitClient()
	data, err := pc.SendRequest(ctx, "GET", fmt.Sprintf("resources/%s/roles/%s", role.ScopeResourceTypeID, id), nil)
	if err != nil {
		em := fmt.Sprintf("Error retrieving role from permit system: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving role from permit system", em)), nil
	}
	res, err := MapToRole(data)
	if err != nil {
		em := fmt.Sprintf("Error retrieving role from permit system: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving role from permit system", em)), nil
	}
	var resdata []models.Data
	resdata = append(resdata, res)
	return utils.FormatSuccess(resdata)

}
func (r *RoleQueryResolver) Roles(ctx context.Context) (models.OperationResult, error) {

	pc := permit.NewPermitClient()
	data, err := pc.SendRequest(ctx, "GET", "resources?include_total_count=true", nil)
	if err != nil {
		em := fmt.Sprintf("Error retrieving roles from permit system: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving roles from permit system", em)), nil
	}
	fmt.Println(data)
	var roles []models.Data
	for _, v := range data["data"].([]interface{}) {
		v := v.(map[string]interface{})
		if _, ok := v["roles"]; !ok {
			continue
		}
		for _, role := range v["roles"].(map[string]interface{}) {
			data, err := MapToRole(role.(map[string]interface{}))
			if err != nil {
				em := fmt.Sprintf("Error retrieving roles from permit system: %v", err)
				logger.LogError(em)
				return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving roles from permit system", em)), nil
			}
			roles = append(roles, *data)
		}
	}

	return utils.FormatSuccess(roles)
}

// Convert map[string]interface{} to Role struct
func MapToRole(roleData map[string]interface{}) (*models.Role, error) {
	fmt.Println(roleData)
	var role models.Role

	data := roleData["attributes"].(map[string]interface{})
	// Parse required string fields
	if idStr, ok := data["ID"].(string); ok {
		roleID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for id")
		}
		role.ID = roleID
	} else {
		return nil, errors.New("missing or invalid id")
	}

	if name, ok := data["Name"].(string); ok {
		role.Name = name
	} else {
		return nil, errors.New("missing or invalid name")
	}

	if createdAt, ok := data["createdAt"].(string); ok {
		role.CreatedAt = createdAt
	} else {
		return nil, errors.New("missing or invalid createdAt")
	}

	if updatedAt, ok := data["updatedAt"].(string); ok {
		role.UpdatedAt = updatedAt
	} else {
		return nil, errors.New("missing or invalid updatedAt")
	}

	// Parse optional description
	if desc, ok := data["Description"].(string); ok {
		role.Description = &desc
	}

	// Parse UUID fields
	if createdByStr, ok := data["createdBy"].(string); ok {
		createdBy, err := uuid.Parse(createdByStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for createdBy")
		}
		role.CreatedBy = createdBy
	} else {
		return nil, errors.New("missing or invalid createdBy")
	}

	if updatedByStr, ok := data["updatedBy"].(string); ok {
		updatedBy, err := uuid.Parse(updatedByStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for updatedBy")
		}
		role.UpdatedBy = updatedBy
	} else {
		return nil, errors.New("missing or invalid updatedBy")
	}

	// Parse role type
	if roleType, ok := data["RoleType"].(string); ok {
		role.RoleType = models.RoleTypeEnum(roleType)
	} else {
		return nil, errors.New("missing or invalid roleType")
	}

	// Parse version
	if version, ok := data["Version"].(string); ok {
		role.Version = version
	} else {
		return nil, errors.New("missing or invalid version")
	}

	// Parse permissions
	if perms, ok := data["Permissions"].([]interface{}); ok {
		var permissions []*models.Permission
		for _, p := range perms {
			if permMap, ok := p.(map[string]interface{}); ok {
				permission, err := parsePermission(permMap)
				if err != nil {
					return nil, err
				}
				permissions = append(permissions, permission)
			}
		}
		role.Permissions = permissions
	}

	// Parse assignable scope (Root)
	if scopeMap, ok := data["AssignableScopeRef"].(map[string]interface{}); ok {
		assignableScope, err := parseRoot(scopeMap)
		if err != nil {
			return nil, err
		}
		role.AssignableScope = assignableScope
	}

	return &role, nil
}

// Parse permission struct
func parsePermission(data map[string]interface{}) (*models.Permission, error) {
	var permission models.Permission

	if idStr, ok := data["permissionId"].(string); ok {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for permissionId")
		}
		permission.ID = id
	} else {
		return nil, errors.New("missing or invalid permissionId")
	}

	if action, ok := data["action"].(string); ok {
		permission.Action = action
	}

	if name, ok := data["name"].(string); ok {
		permission.Name = name
		permission.Action = name
	}

	if resourcetypeId, ok := data["resourcetypeId"].(string); ok {
		permission.AssignableScope = resourcetypeId
	}

	if createdAt, ok := data["createdAt"].(string); ok {
		permission.CreatedAt = createdAt
	}

	if updatedAt, ok := data["updatedAt"].(string); ok {
		permission.UpdatedAt = updatedAt
	}

	if createdByStr, ok := data["createdBy"].(string); ok {
		createdBy, err := uuid.Parse(createdByStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for createdBy")
		}
		permission.CreatedBy = createdBy
	}

	if updatedByStr, ok := data["updatedBy"].(string); ok {
		updatedBy, err := uuid.Parse(updatedByStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for updatedBy")
		}
		permission.UpdatedBy = updatedBy
	}

	return &permission, nil
}

// Parse Root struct (AssignableScope)
func parseRoot(data map[string]interface{}) (*models.Root, error) {
	var root models.Root

	if idStr, ok := data["resource_type_id"].(string); ok {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for resource_type_id")
		}
		root.ID = id
	} else {
		return nil, errors.New("missing or invalid resource_type_id")
	}

	if name, ok := data["name"].(string); ok {
		root.Name = name
	}

	if createdAt, ok := data["created_at"].(string); ok {
		root.CreatedAt = createdAt
	}

	if updatedAt, ok := data["updated_at"].(string); ok {
		root.UpdatedAt = updatedAt
	}

	if createdByStr, ok := data["created_by"].(string); ok {
		createdBy, err := uuid.Parse(createdByStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for created_by")
		}
		root.CreatedBy = createdBy
	}

	if updatedByStr, ok := data["updated_by"].(string); ok {
		updatedBy, err := uuid.Parse(updatedByStr)
		if err != nil {
			return nil, errors.New("invalid UUID format for updated_by")
		}
		root.UpdatedBy = updatedBy
	}

	return &root, nil
}
