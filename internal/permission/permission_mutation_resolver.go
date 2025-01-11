package permission

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
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type PermissionMutationResolver struct {
	DB *gorm.DB
}

func (r *PermissionMutationResolver) DeletePermission(ctx context.Context, id uuid.UUID) (bool, error) {
	updates := map[string]interface{}{
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}
	result := r.DB.Model(&dto.TNTPermission{}).Where("permission_id = ?", id).Updates(updates)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, errors.New("permission not found")
	}

	return true, nil
}

func (r *PermissionMutationResolver) CreatePermission(ctx context.Context, input *models.CreatePermission) (*models.Permission, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}

	if input.RoleID == nil {
		return nil, errors.New("role ID is required")
	} else {
		if err := utils.ValidateRoleID(*input.RoleID); err != nil {
			return nil, fmt.Errorf("invalid role ID: %v", err)
		}
	}

	if input.ServiceID == nil {
		return nil, errors.New("service ID is required")
	}

	if input.Action == nil {
		return nil, errors.New("action is required")
	}

	permission := &dto.TNTPermission{
		PermissionID: uuid.New(),
		Name:         input.Name,
		ServiceID:    *input.ServiceID,
		Action:       *input.Action,
		RowStatus:    1,
		RoleID:       *input.RoleID,
		CreatedBy:    constants.DefaltCreatedBy, //input.UpdatedBy,
		UpdatedBy:    constants.DefaltUpdatedBy, //input.UpdatedBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := r.DB.Create(permission).Error; err != nil {
		return nil, err
	}

	// Convert to GraphQL model
	return &models.Permission{
		ID:        permission.PermissionID,
		Name:      permission.Name,
		ServiceID: &permission.ServiceID,
		Action:    &permission.Action,
		CreatedAt: ptr.String(permission.CreatedAt.String()),
		CreatedBy: permission.CreatedBy,
		UpdatedAt: ptr.String(permission.UpdatedAt.String()),
		UpdatedBy: ptr.String(permission.UpdatedBy),
	}, nil
}

func (r *PermissionMutationResolver) UpdatePermission(ctx context.Context, permissionID uuid.UUID, input *models.UpdatePermission) (*models.Permission, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.ServiceID == nil {
		return nil, errors.New("service ID is required")
	}

	if input.Action == nil {
		return nil, errors.New("action is required")
	}

	// Find existing permission
	var permission dto.TNTPermission
	if err := r.DB.First(&permission, "permission_id = ?", permissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	// if input.UpdatedBy == "" {
	// 	return nil, errors.New("updatedBy is required")
	// }

	if input.RoleID == nil {
		return nil, errors.New("role ID is required")
	} else if err := utils.ValidateRoleID(*input.RoleID); err != nil {
		return nil, fmt.Errorf("invalid role ID: %v", err)
	}

	// Update fields
	permission.Name = input.Name
	permission.ServiceID = *input.ServiceID
	permission.Action = *input.Action
	permission.RoleID = *input.RoleID
	// permission.UpdatedBy = input.UpdatedBy
	permission.UpdatedAt = time.Now()

	// Save updates
	if err := r.DB.Save(&permission).Error; err != nil {
		return nil, err
	}

	// Convert to GraphQL model
	return &models.Permission{
		ID:        permission.PermissionID,
		Name:      permission.Name,
		ServiceID: &permission.ServiceID,
		Action:    &permission.Action,
		CreatedAt: ptr.String(permission.CreatedAt.String()),
		CreatedBy: permission.CreatedBy,
		UpdatedAt: ptr.String(permission.UpdatedAt.String()),
		UpdatedBy: ptr.String(permission.UpdatedBy),
	}, nil
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
