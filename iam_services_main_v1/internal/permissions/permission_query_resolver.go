package permissions

import (
	"context"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/permit"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PermissionQueryResolver handles permission-related queries.
type PermissionQueryResolver struct {
	DB     *gorm.DB
	Permit *permit.PermitClient
}

// Permissions resolves the list of all permissions.
func (r *PermissionQueryResolver) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	//logger.Log.Info("Fetching all permissions")

	var permissions []dto.MstPermission
	if err := r.DB.Find(&permissions).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch permissions from the database")
		return nil, err
	}

	// pc := permit.NewPermitClient()

	// update := map[string]interface{}{
	// 	"permissions": permissionAction,
	// }

	// _, err := pc.SendRequest(ctx, "POST", fmt.Sprintf("resources/%s/roles/%s/permissions", assinableScopeName, roleName), update)
	// if err != nil {
	// 	return err
	// }
	var result []*models.Permission
	for _, permission := range permissions {
		result = append(result, &models.Permission{
			ID:   permission.PermissionID,
			Name: permission.Name,
			// ServiceID: &permission.ServiceID,
			// RoleID:    &permission.RoleID,
			// Action:    &permission.Action,
			CreatedAt: permission.CreatedAt.Format(time.RFC3339),
			// CreatedBy: permission.CreatedBy,
			UpdatedAt: permission.UpdatedAt.Format(time.RFC3339),
			// UpdatedBy: &permission.UpdatedBy,
		})
	}
	// // https://api.permit.io/v2/schema/{proj_id}/{env_id}/resources/{resource_id}/actions
	// res, err := r.Permit.SendRequest(ctx,"GET", fmt.Sprintf("resources/%s/actions", input.Resource))
	// fmt.Println(res, err)
	// if err != nil {
	// 	return nil, err
	// }

	//logger.Log.Infof("Fetched %d permissions", len(result))
	return result, nil
}

// GetPermission resolves a single permission by ID.
func (r *PermissionQueryResolver) GetPermission(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	//logger.Log.Infof("Fetching permission with ID: %s", id)

	var permission dto.MstPermission
	if err := r.DB.First(&permission, "permission_id = ?", id).Error; err != nil {
		//logger.AddContext(err).Warnf("Permission with ID %s not found", id)
		return nil, err
	}

	//logger.Log.Infof("Permission with ID %s fetched successfully", id)
	return &models.Permission{
		ID:   permission.PermissionID,
		Name: permission.Name,
		// ServiceID: &permission.ServiceID,
		// Action:    &permission.Action,
		// RoleID:    &permission.RoleID,
		CreatedAt: permission.CreatedAt.Format(time.RFC3339),
		// CreatedBy: permission.CreatedBy,
		UpdatedAt: permission.UpdatedAt.Format(time.RFC3339),
		// UpdatedBy: permission.UpdatedBy,
	}, nil
}
