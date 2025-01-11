package role

import (
	"context"
	"errors"

	"go_graphql/config"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

// RoleQueryResolver handles role-related queries.
type RoleQueryResolver struct {
	DB *gorm.DB
}

// Roles resolves the list of all roles.
func (r *RoleQueryResolver) AllRoles(ctx context.Context) ([]*models.Role, error) {
	var roles []dto.TNTRole
	if err := r.DB.Find(&roles).Error; err != nil {
		return nil, err
	}

	var result []*models.Role
	for _, role := range roles {
		convertedRole := convertRoleToGraphQL(&role)
		result = append(result, convertedRole)
	}
	return result, nil
}

// GetRole resolves a single role by ID.
func (r *RoleQueryResolver) GetRole(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role dto.TNTRole
	if err := r.DB.First(&role, "resource_id = ?", id).Error; err != nil {
		return nil, errors.New("role not found")
	}
	return convertRoleToGraphQL(&role), nil
}

// Helper function to convert database Role to GraphQL Role models.
func convertRoleToGraphQL(role *dto.TNTRole) *models.Role {
	var permissions []dto.TNTPermission

	tx := config.DB.Where("role_id = ?", role.ResourceID).Find(&permissions)
	if tx.Error != nil {
		return nil
	}

	permissionsIds := []*string{}
	for _, permission := range permissions {
		permissionsIds = append(permissionsIds, ptr.String(permission.PermissionID.String()))
	}
	res := &models.Role{
		ID:             role.ResourceID,
		Name:           role.Name,
		Description:    ptr.String(role.Description),
		RoleType:       models.RoleTypeEnum(role.RoleType),
		Version:        &role.Version,
		CreatedAt:      role.CreatedAt.String(),
		UpdatedAt:      ptr.String(role.UpdatedAt.String()),
		UpdatedBy:      &role.UpdatedBy,
		CreatedBy:      &role.CreatedBy,
		PermissionsIds: permissionsIds,
	}

	var parentOrg dto.TenantResource
	if role.ParentResourceID != nil {
		if err := config.DB.Where(&dto.TenantResource{ResourceID: *role.ParentResourceID}).First(&parentOrg).Error; err != nil {
			return nil
		}
		res.ParentOrg = &models.Root{
			ID:        parentOrg.ResourceID,
			Name:      parentOrg.Name,
			CreatedAt: parentOrg.CreatedAt.String(),
			UpdatedAt: ptr.String(parentOrg.UpdatedAt.String()),
			CreatedBy: &parentOrg.CreatedBy,
			UpdatedBy: &parentOrg.UpdatedBy,
		}
	}
	return res
}
