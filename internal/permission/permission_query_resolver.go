package permission

import (
	"context"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"time"

	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

// PermissionQueryResolver handles permission-related queries.
type PermissionQueryResolver struct {
	DB *gorm.DB
}

// Permissions resolves the list of all permissions.
func (r *PermissionQueryResolver) GetPermission(ctx context.Context) ([]*models.Permission, error) {
	var permissions []dto.Permission
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
			CreatedAt: permission.CreatedDate.Format(time.RFC3339),
			CreatedBy: permission.CreatedBy,
			UpdatedAt: ptr.String(permission.UpdatedDate.Format(time.RFC3339)),
			UpdatedBy: &permission.UpdatedBy,
		})
	}

	return result, nil
}
