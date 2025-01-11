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
	"go_graphql/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	DB *gorm.DB
}

// CreateRole handles creating a new role.
func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.RoleInput) (*models.Role, error) {
	logger.Log.Info("Starting CreateRole")

	// Validate input
	if input.Name == "" {
		logger.Log.Warn("Role name is required")
		return nil, errors.New("role name is required")
	}

	if input.ParentOrgID == uuid.Nil {
		logger.Log.Warn("ParentOrgID is required")
		return nil, errors.New("resource ID is required")
	} else {
		if err := utils.ValidateResourceID(input.ParentOrgID); err != nil {
			logger.AddContext(err).Error("Invalid ParentOrgID")
			return nil, fmt.Errorf("invalid ParentOrgID")
		}
	}

	role := dto.TNTRole{}
	role.Name = input.Name
	if input.Description != nil {
		role.Description = *input.Description
	}
	role.ParentResourceID = &input.ParentOrgID
	role.RoleType = "CUSTOM"
	role.Version = input.Version
	role.CreatedAt = time.Now()
	role.CreatedBy = constants.DefaltCreatedBy
	role.UpdatedBy = constants.DefaltCreatedBy
	role.RowStatus = 1

	logger.Log.Info("Fetching resource type")
	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ?", "Role").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	logger.Log.Info("Creating tenant resource")
	tenantResource := &dto.TenantResource{
		ResourceID:       uuid.New(),
		ParentResourceID: &input.ParentOrgID,
		ResourceTypeID:   resourceType.ResourceTypeID,
		Name:             input.Name,
		RowStatus:        1,
		CreatedBy:        constants.DefaltCreatedBy,
		UpdatedBy:        constants.DefaltCreatedBy,
		CreatedAt:        time.Now(),
		TenantID:         uuid.New(),
	}

	if err := r.DB.Create(&tenantResource).Error; err != nil {
		logger.AddContext(err).Error("Failed to create tenant resource")
		return nil, err
	}

	role.ResourceID = tenantResource.ResourceID
	logger.Log.Info("Saving role to database")
	if err := r.DB.Create(&role).Error; err != nil {
		logger.AddContext(err).Error("Failed to save role")
		return nil, err
	}

	logger.Log.Info("Role created successfully")
	return convertRoleToGraphQL(&role), nil
}

// UpdateRole handles updating an existing role.
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, id uuid.UUID, input models.RoleInput) (*models.Role, error) {
	logger.Log.Infof("Starting UpdateRole for ID: %s", id)

	var role dto.TNTRole
	if err := r.DB.First(&role, "resource_id = ?", id).Error; err != nil {
		logger.AddContext(err).Warn("Role not found")
		return nil, errors.New("role not found")
	}

	// Update fields
	if input.Name != "" {
		role.Name = input.Name
	}
	role.Version = input.Version
	role.ParentResourceID = &input.ParentOrgID
	role.UpdatedBy = constants.DefaltCreatedBy
	role.UpdatedAt = time.Now()

	logger.Log.Infof("Updating role in database for ID: %s", id)
	updateData := map[string]interface{}{
		"name":               role.Name,
		"parent_resource_id": role.ParentResourceID,
		"version":            role.Version,
		"updated_by":         role.UpdatedBy,
		"updated_at":         role.UpdatedAt,
	}

	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ?", id).Updates(updateData).Error; err != nil {
		logger.AddContext(err).Error("Failed to update role")
		return nil, err
	}

	logger.Log.Infof("Role updated successfully for ID: %s", id)
	return convertRoleToGraphQL(&role), nil
}

// DeleteRole handles deleting a role by ID.
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, id uuid.UUID) (bool, error) {
	logger.Log.Infof("Starting DeleteRole for ID: %s", id)

	var roleDB dto.TNTRole
	if err := r.DB.First(&roleDB, "resource_id = ?", id).Error; err != nil {
		logger.AddContext(err).Warn("Role not found")
		return false, errors.New("role not found")
	}

	updates := map[string]interface{}{
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}

	logger.Log.Infof("Marking role as deleted for ID: %s", id)
	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ?", id).Updates(updates).Error; err != nil {
		logger.AddContext(err).Error("Failed to delete role")
		return false, err
	}

	logger.Log.Infof("Role deleted successfully for ID: %s", id)
	return true, nil
}
