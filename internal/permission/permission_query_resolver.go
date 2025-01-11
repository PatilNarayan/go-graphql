package permission

import (
	"context"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
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
	var permissions []dto.TNTPermission
	if err := r.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}

	var result []*models.Permission
	for _, permission := range permissions {
		result = append(result, &models.Permission{
			ID:        permission.PermissionID,
			Name:      permission.Name,
			ServiceID: &permission.ServiceID,
			Action:    &permission.Action,
			CreatedAt: ptr.String(permission.CreatedAt.Format(time.RFC3339)),
			CreatedBy: permission.CreatedBy,
			UpdatedAt: ptr.String(permission.UpdatedAt.Format(time.RFC3339)),
			UpdatedBy: &permission.UpdatedBy,
		})
	}

	return result, nil
}

func (r *PermissionQueryResolver) GetPermission(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	var permission dto.TNTPermission
	if err := r.DB.First(&permission, "permission_id = ?", id).Error; err != nil {
		return nil, err
	}

	return &models.Permission{
		ID:        permission.PermissionID,
		Name:      permission.Name,
		ServiceID: &permission.ServiceID,
		Action:    &permission.Action,
		CreatedAt: ptr.String(permission.CreatedAt.Format(time.RFC3339)),
		CreatedBy: permission.CreatedBy,
		UpdatedAt: ptr.String(permission.UpdatedAt.Format(time.RFC3339)),
		UpdatedBy: &permission.UpdatedBy,
	}, nil
}
