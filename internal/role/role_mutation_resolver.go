package role

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	// Check permissionIds are valid
	if input.PermissionsIds != nil && len(input.PermissionsIds) > 0 {
		if err := validatePermissions(r.DB, input.PermissionsIds); err != nil {
			return nil, fmt.Errorf("invalid permissions: %v", err)
		}
	}

	// Convert PermissionsIDs to JSON
	permissionsJSON, err := json.Marshal(input.PermissionsIds)
	if err != nil {
		return nil, fmt.Errorf("failed to convert permissions to JSON: %v", err)
	}

	// Create a new role entity
	role := dto.Role{
		RoleID: uuid.New().String(),
	}
	if input.Name != "" {
		role.Name = input.Name
	}
	if input.Description != nil {
		role.Description = *input.Description
	}
	if input.ResourceID != nil {
		role.ResourceID = *input.ResourceID
	}
	if input.RoleType != "" {
		role.RoleType = string(input.RoleType)
	}
	if input.PermissionsIds != nil {
		role.PermissionsIDs = string(permissionsJSON)
	}
	if input.Version != nil {
		role.Version = *input.Version
	}
	role.CreatedAt = time.Now()
	role.CreatedBy = input.CreatedBy
	role.UpdatedBy = input.CreatedBy

	// Save to the database
	err = r.DB.Create(&role).Error
	if err != nil {
		return nil, err
	}

	// Create role assignments
	for _, permissionID := range input.PermissionsIds {
		permissionAssignment := dto.RoleAssignment{
			RoleAssignmentID: uuid.New().String(),
			RoleID:           role.RoleID,
			PermissionID:     permissionID,
			CreatedAt:        time.Now(),
			CreatedBy:        input.CreatedBy,
			UpdatedBy:        input.CreatedBy,
		}

		if err := r.DB.Create(&permissionAssignment).Error; err != nil {
			return nil, err
		}
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

	if input.ResourceID != nil {
		role.ResourceID = *input.ResourceID
	}

	if input.UpdatedBy != nil {
		role.UpdatedBy = *input.UpdatedBy
	} else {
		return nil, errors.New("updatedBy is required")
	}

	role.UpdatedAt = time.Now()

	// Update PermissionsIDs and RoleAssignments
	if input.PermissionsIds != nil {
		permissionsJSON, err := json.Marshal(input.PermissionsIds)
		if err != nil {
			return nil, fmt.Errorf("failed to convert permissions to JSON: %v", err)
		}
		role.PermissionsIDs = string(permissionsJSON)

		// Update role assignments
		if err := r.DB.Where("role_id = ?", id).Delete(&dto.RoleAssignment{}).Error; err != nil {
			return nil, fmt.Errorf("failed to clear old role assignments: %v", err)
		}
		for _, permissionID := range input.PermissionsIds {
			permissionAssignment := dto.RoleAssignment{
				RoleAssignmentID: uuid.New().String(),
				RoleID:           id,
				PermissionID:     permissionID,
				CreatedAt:        time.Now(),
				CreatedBy:        *input.UpdatedBy,
				UpdatedBy:        *input.UpdatedBy,
			}
			if err := r.DB.Create(&permissionAssignment).Error; err != nil {
				return nil, fmt.Errorf("failed to create role assignment: %v", err)
			}
		}
	}

	// Save changes explicitly using UpdateColumns
	updateData := map[string]interface{}{
		"name":            role.Name,
		"description":     role.Description,
		"role_type":       role.RoleType,
		"permissions_ids": role.PermissionsIDs,
		"version":         role.Version,
		"updated_by":      role.UpdatedBy,
		"updated_at":      role.UpdatedAt,
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
	// Delete role assignments
	if err := r.DB.Where("role_id = ?", id).Delete(&dto.RoleAssignment{}).Error; err != nil {
		return false, fmt.Errorf("failed to delete role assignments: %v", err)
	}

	// Attempt to delete the role
	if err := r.DB.Delete(&dto.Role{}, "role_id = ?", id).Error; err != nil {
		return false, err
	}
	return true, nil
}

func validatePermissions(db *gorm.DB, permissionIDs []string) error {
	for _, permissionID := range permissionIDs {
		var permission dto.Permission
		if err := db.First(&permission, "permission_id = ?", permissionID).Error; err != nil {
			return fmt.Errorf("permission %s not found", permissionID)
		}
	}
	return nil
}
