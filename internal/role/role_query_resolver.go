package role

import (
	"context"
	"errors"

	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

// RoleQueryResolver handles role-related queries.
type RoleQueryResolver struct {
	DB *gorm.DB
}

// Roles resolves the list of all roles.
func (r *RoleQueryResolver) Roles(ctx context.Context) ([]*models.Role, error) {
	var roles []dto.Role
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
func (r *RoleQueryResolver) GetRole(ctx context.Context, id string) (*models.Role, error) {
	var role dto.Role
	if err := r.DB.First(&role, "role_id = ?", id).Error; err != nil {
		return nil, errors.New("role not found")
	}
	return convertRoleToGraphQL(&role), nil
}

// Helper function to convert database Role to GraphQL Role models.
func convertRoleToGraphQL(role *dto.Role) *models.Role {
	return &models.Role{
		ID:          role.RoleID,
		Name:        role.Name,
		Description: ptr.String(role.Description),
		RoleType:    models.RoleTypeEnum(role.RoleType),
		Version:     &role.Version,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   ptr.String(role.UpdatedAt.String()),
	}
}

// Helper function to convert permissions.
func convertPermissionsToGraphQL(permissions []dto.Permission) []string {
	var result []string
	for _, permission := range permissions {
		result = append(result, permission.Name)
	}
	return result
}
