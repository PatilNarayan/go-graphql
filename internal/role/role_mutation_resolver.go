package role

import (
	"context"
	"errors"
	"time"

	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	DB *gorm.DB
}

// CreateRole handles creating a new role.
func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.RoleInput) (*models.Role, error) {
	// Validate input
	if input.Name == "" {
		return nil, errors.New("role name is required")
	}

	// Create a new role entity
	role := dto.Role{
		RoleID:      uuid.New().String(),
		Name:        input.Name,
		Description: *input.Description,
		RoleType:    string(input.RoleType),
		Version:     *input.Version,
		CreatedAt:   time.Now(),
		CreatedBy:   input.CreatedBy,
		UpdatedBy:   *input.UpdatedBy,
	}

	// Save to the database
	err := r.DB.Create(&role).Error
	if err != nil {
		return nil, err
	}

	return convertRoleToGraphQL(&role), nil
}

// UpdateRole handles updating an existing role.
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, id string, input models.RoleInput) (*models.Role, error) {
	// Fetch the existing role
	var role dto.Role
	if err := r.DB.First(&role, "role_id = ?", id).Error; err != nil {
		return nil, errors.New("role not found")
	}

	// Validate and update fields
	if input.Name != "" {
		role.Name = input.Name
	}
	if input.Description != nil {
		role.Description = *input.Description
	}
	if input.RoleType != "" {
		role.RoleType = string(input.RoleType)
	}
	if input.Version != nil {
		role.Version = *input.Version
	}
	if input.UpdatedBy != nil {
		role.UpdatedBy = *input.UpdatedBy
	} else {
		return nil, errors.New("updatedBy is required")
	}

	role.UpdatedAt = time.Now()

	// Save changes explicitly using UpdateColumns
	updateData := map[string]interface{}{
		"name":        role.Name,
		"description": role.Description,
		"role_type":   role.RoleType,
		"version":     role.Version,
		"updated_by":  role.UpdatedBy,
		"updated_at":  role.UpdatedAt,
	}

	if err := r.DB.Model(&dto.Role{}).Where("role_id = ?", id).Updates(updateData).Error; err != nil {
		return nil, err
	}

	var updatedData dto.Role
	if err := r.DB.First(&updatedData, "role_id = ?", id).Error; err != nil {
		return nil, errors.New("role not found")
	}

	return convertRoleToGraphQL(&updatedData), nil
}

// DeleteRole handles deleting a role by ID.
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, id string) (bool, error) {
	// Attempt to delete the role
	if err := r.DB.Delete(&dto.Role{}, "role_id = ?", id).Error; err != nil {
		return false, err
	}
	return true, nil
}
