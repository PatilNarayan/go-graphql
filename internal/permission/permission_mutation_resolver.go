package permission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go_graphql/gql/models"
	"go_graphql/internal/constants"
	"go_graphql/internal/dto"
	"go_graphql/logger"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type PermissionMutationResolver struct {
	DB *gorm.DB
}

func (r *PermissionMutationResolver) DeletePermission(ctx context.Context, id uuid.UUID) (bool, error) {
	logger.Log.Infof("Attempting to delete permission with ID: %s", id)

	updates := map[string]interface{}{
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}
	result := r.DB.Model(&dto.TNTPermission{}).Where("permission_id = ?", id).Updates(updates)
	if result.Error != nil {
		logger.AddContext(result.Error).Errorf("Failed to delete permission with ID: %s", id)
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		logger.Log.Warnf("Permission with ID %s not found", id)
		return false, errors.New("permission not found")
	}

	logger.Log.Infof("Permission with ID %s deleted successfully", id)
	return true, nil
}

func (r *PermissionMutationResolver) CreatePermission(ctx context.Context, input *models.CreatePermission) (*models.Permission, error) {
	logger.Log.Info("Starting permission creation")

	if input == nil {
		logger.Log.Error("Input is required for creating permission")
		return nil, errors.New("input is required")
	}

	// if input.RoleID == nil {
	// 	logger.Log.Error("Role ID is required for creating permission")
	// 	return nil, errors.New("role ID is required")
	// } else {
	// 	if err := utils.ValidateRoleID(*input.RoleID); err != nil {
	// 		logger.AddContext(err).Error("Invalid role ID")
	// 		return nil, fmt.Errorf("invalid role ID: %v", err)
	// 	}
	// }

	if input.ServiceID == uuid.Nil {
		logger.Log.Error("Service ID is required for creating permission")
		return nil, errors.New("service ID is required")
	}

	if input.Action == "" {
		logger.Log.Error("Action is required for creating permission")
		return nil, errors.New("action is required")
	}

	if input.Name == "" {
		logger.Log.Error("Name is required for creating permission")
		return nil, errors.New("name is required")
	} else {
		if err := r.DB.Where("name = ?", input.Name).First(&dto.TNTPermission{}).Error; err == nil {
			logger.Log.Errorf("Permission with name %s already exists", input.Name)
			return nil, fmt.Errorf("permission with name %s already exists", input.Name)
		}
	}

	permission := &dto.TNTPermission{
		PermissionID: uuid.New(),
		Name:         input.Name,
		ServiceID:    input.ServiceID.String(),
		Action:       input.Action,
		RowStatus:    1,
		// RoleID:       *input.RoleID,
		CreatedBy: constants.DefaltCreatedBy,
		UpdatedBy: constants.DefaltUpdatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.DB.Create(permission).Error; err != nil {
		logger.AddContext(err).Error("Failed to create permission")
		return nil, err
	}

	logger.Log.Infof("Permission with ID %s created successfully", permission.PermissionID)

	return convertToGraphQLPermission(permission), nil
}

func (r *PermissionMutationResolver) UpdatePermission(ctx context.Context, input *models.UpdatePermission) (*models.Permission, error) {
	permissionID := input.ID
	logger.Log.Infof("Starting update for permission with ID: %s", permissionID)

	if input == nil {
		logger.Log.Error("Input is required for updating permission")
		return nil, errors.New("input is required")
	}

	if input.Name == "" {
		logger.Log.Error("Name is required for updating permission")
		return nil, errors.New("name is required")
	} else {
		if err := r.DB.Where("name = ?", input.Name).First(&dto.TNTPermission{}).Error; err == nil {
			logger.Log.Errorf("Permission with name %s already exists", input.Name)
			return nil, fmt.Errorf("permission with name %s already exists", input.Name)
		}
	}

	if input.ServiceID == nil {
		logger.Log.Error("Service ID is required for updating permission")
		return nil, errors.New("service ID is required")
	}

	if input.Action == nil {
		logger.Log.Error("Action is required for updating permission")
		return nil, errors.New("action is required")
	}

	var permission dto.TNTPermission
	if err := r.DB.First(&permission, "permission_id = ?", permissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("Permission with ID %s not found", permissionID)
			return nil, errors.New("permission not found")
		}
		logger.AddContext(err).Errorf("Failed to find permission with ID: %s", permissionID)
		return nil, err
	}

	// if input.RoleID == nil {
	// 	logger.Log.Error("Role ID is required for updating permission")
	// 	return nil, errors.New("role ID is required")
	// } else if err := utils.ValidateRoleID(*input.RoleID); err != nil {
	// 	logger.AddContext(err).Error("Invalid role ID")
	// 	return nil, fmt.Errorf("invalid role ID: %v", err)
	// }

	permission.Name = input.Name
	permission.ServiceID = input.ServiceID.String()
	permission.Action = *input.Action
	// permission.RoleID = *input.RoleID
	permission.UpdatedAt = time.Now()

	if err := r.DB.Save(&permission).Error; err != nil {
		logger.AddContext(err).Errorf("Failed to update permission with ID: %s", permissionID)
		return nil, err
	}

	logger.Log.Infof("Permission with ID %s updated successfully", permissionID)
	return convertToGraphQLPermission(&permission), nil
}

// Helper function to convert DTO to GraphQL model
func convertToGraphQLPermission(p *dto.TNTPermission) *models.Permission {
	if p == nil {
		return nil
	}

	return &models.Permission{
		ID:        p.PermissionID,
		Name:      p.Name,
		ServiceID: &p.ServiceID,
		Action:    &p.Action,
		CreatedAt: ptr.String(p.CreatedAt.String()),
		CreatedBy: p.CreatedBy,
		UpdatedAt: ptr.String(p.UpdatedAt.String()),
		UpdatedBy: ptr.String(p.UpdatedBy),
	}
}
