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
	role := &dto.Role{
		RoleID:      uuid.New().String(),
		Name:        input.Name,
		Description: *input.Description,
		RoleType:    string(input.RoleType),
		Version:     input.Version,
		CreatedAt:   time.Now(),
	}

	// Save to the database
	err := r.DB.Create(role).Error
	if err != nil {
		return nil, err
	}

	return convertRoleToGraphQL(role), nil
}

// UpdateRole handles updating an existing role.
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, id string, input models.RoleInput) (*models.Role, error) {
	// Fetch the existing role
	var role dto.Role
	if err := r.DB.First(&role, "role_id = ?", id).Error; err != nil {
		return nil, errors.New("role not found")
	}

	// Update fields
	if input.Name != "" {
		role.Name = input.Name
	}
	if input.Description != nil {
		role.Description = *input.Description
	}
	role.RoleType = string(input.RoleType)
	role.Version = input.Version
	role.UpdatedAt = time.Now()

	// Save changes to the database
	if err := r.DB.Save(&role).Error; err != nil {
		return nil, err
	}

	return convertRoleToGraphQL(&role), nil
}

// DeleteRole handles deleting a role by ID.
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, id string) (bool, error) {
	// Attempt to delete the role
	if err := r.DB.Delete(&dto.Role{}, "role_id = ?", id).Error; err != nil {
		return false, err
	}
	return true, nil
}
