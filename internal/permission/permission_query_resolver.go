package permission

import (
	"context"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"go_graphql/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

// PermissionQueryResolver handles permission-related queries.
type PermissionQueryResolver struct {
	DB *gorm.DB
}

// Permissions resolves the list of all permissions.
func (r *PermissionQueryResolver) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	logger.Log.Info("Fetching all permissions")

	var permissions []dto.TNTPermission
	if err := r.DB.Find(&permissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch permissions from the database")
		return nil, err
	}

	var result []*models.Permission
	for _, permission := range permissions {
		result = append(result, &models.Permission{
			ID:        permission.PermissionID,
			Name:      permission.Name,
			ServiceID: &permission.ServiceID,
			// RoleID:    &permission.RoleID,
			Action:    &permission.Action,
			CreatedAt: ptr.String(permission.CreatedAt.Format(time.RFC3339)),
			// CreatedBy: permission.CreatedBy,
			UpdatedAt: ptr.String(permission.UpdatedAt.Format(time.RFC3339)),
			// UpdatedBy: &permission.UpdatedBy,
		})
	}

	logger.Log.Infof("Fetched %d permissions", len(result))
	return result, nil
}

// GetPermission resolves a single permission by ID.
func (r *PermissionQueryResolver) GetPermission(ctx context.Context, id *uuid.UUID) (*models.Permission, error) {
	logger.Log.Infof("Fetching permission with ID: %s", *id)

	var permission dto.TNTPermission
	if err := r.DB.First(&permission, "permission_id = ?", *id).Error; err != nil {
		logger.AddContext(err).Warnf("Permission with ID %s not found", *id)
		return nil, err
	}

	logger.Log.Infof("Permission with ID %s fetched successfully", *id)
	return &models.Permission{
		ID:        permission.PermissionID,
		Name:      permission.Name,
		ServiceID: &permission.ServiceID,
		Action:    &permission.Action,
		// RoleID:    &permission.RoleID,
		CreatedAt: ptr.String(permission.CreatedAt.Format(time.RFC3339)),
		CreatedBy: permission.CreatedBy,
		UpdatedAt: ptr.String(permission.UpdatedAt.Format(time.RFC3339)),
		UpdatedBy: &permission.UpdatedBy,
	}, nil
}
