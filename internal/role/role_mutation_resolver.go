package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go_graphql/gql/models"
	"go_graphql/internal/constants"
	"go_graphql/internal/dto"
	"go_graphql/internal/utils"

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

	if input.ParentOrgID == uuid.Nil {
		return nil, errors.New("resource ID is required")
	} else {
		if err := utils.ValidateResourceID(input.ParentOrgID); err != nil {
			return nil, fmt.Errorf("invalid ParentOrgID")
		}
	}

	role := dto.TNTRole{}
	if input.Name != "" {
		role.Name = input.Name
	}
	if input.Description != nil {
		role.Description = *input.Description
	}
	if input.ParentOrgID != uuid.Nil {
		role.ParentResourceID = &input.ParentOrgID
	}

	role.RoleType = "CUSTOM"

	if input.Version != "" {
		role.Version = input.Version
	}

	role.CreatedAt = time.Now()
	role.CreatedBy = constants.DefaltCreatedBy //input.CreatedBy
	role.UpdatedBy = constants.DefaltCreatedBy //input.CreatedBy
	role.RowStatus = 1

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ?", "Role").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	tenantResource := &dto.TenantResource{
		ResourceID:       uuid.New(),
		ParentResourceID: &input.ParentOrgID,
		ResourceTypeID:   resourceType.ResourceTypeID,
		Name:             input.Name,
		RowStatus:        1,
		CreatedBy:        constants.DefaltCreatedBy, // input.CreatedBy,
		UpdatedBy:        constants.DefaltCreatedBy, // input.CreatedBy,
		CreatedAt:        time.Now(),
	}

	tenantResource.TenantID = tenantResource.ResourceID

	if err := r.DB.Create(&tenantResource).Error; err != nil {
		return nil, err
	}

	// Save to the database
	role.ResourceID = tenantResource.ResourceID
	err := r.DB.Create(&role).Error
	if err != nil {
		return nil, err
	}

	return convertRoleToGraphQL(&role), nil
}

// UpdateRole handles updating an existing role.
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, id uuid.UUID, input models.RoleInput) (*models.Role, error) {
	// Fetch the existing role
	var role dto.TNTRole
	if err := r.DB.First(&role, "resource_id = ?", id).Error; err != nil {
		return nil, errors.New("role not found")
	}

	if input.ParentOrgID == uuid.Nil {
		return nil, errors.New("resource ID is required")
	} else {
		if err := utils.ValidateResourceID(input.ParentOrgID); err != nil {
			return nil, fmt.Errorf("invalid resource ID: %v", err)
		}
	}

	// Validate and update fields
	if input.Name != "" {
		role.Name = input.Name
	}
	// if input.Description != nil {
	// 	role.Description = *input.Description
	// }
	// if input.RoleType != "" {
	// 	role.RoleType = string(input.RoleType)
	// }
	if input.Version != "" {
		role.Version = input.Version
	}

	if input.ParentOrgID != uuid.Nil {
		role.ParentResourceID = &input.ParentOrgID
	}

	// if input.UpdatedBy != nil {
	// 	role.UpdatedBy = *input.UpdatedBy
	// } else {
	// 	return nil, errors.New("updatedBy is required")
	// }

	role.UpdatedBy = constants.DefaltCreatedBy

	role.UpdatedAt = time.Now()

	// Save changes explicitly using UpdateColumns
	updateData := map[string]interface{}{
		"name": role.Name,
		// "description":     role.Description,
		"parent_resource_id": role.ParentResourceID,
		// "permissions_ids": role.PermissionsIDs,
		"version":    role.Version,
		"updated_by": role.UpdatedBy,
		"updated_at": role.UpdatedAt,
	}

	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ?", id).Updates(updateData).Error; err != nil {
		return nil, err
	}

	var updatedData dto.TNTRole
	if err := r.DB.First(&updatedData, "resource_id = ?", id).Error; err != nil {
		return nil, errors.New("role not found")
	}

	return convertRoleToGraphQL(&updatedData), nil
}

// DeleteRole handles deleting a role by ID.
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, id uuid.UUID) (bool, error) {

	var roleDB dto.TNTRole
	if err := r.DB.First(&roleDB, "resource_id = ?", id).Error; err != nil {
		return false, errors.New("role not found")
	}

	updates := map[string]interface{}{
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}

	// Attempt to delete the role
	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ?", id).Updates(updates).Error; err != nil {
		return false, err
	}

	if err := r.DB.Model(&dto.TenantResource{}).Where("resource_id = ?", id).Updates(updates).Error; err != nil {
		return false, err
	}

	if err := r.DB.Model(&dto.TNTPermission{}).Where("role_id = ?", id).Updates(updates).Error; err != nil {
		return false, err
	}

	return true, nil
}
