package roles

import (
	"context"
	"iam_services_main_v1/gql/models"

	"gorm.io/gorm"
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	DB *gorm.DB
}

func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.CreateRoleInput) (models.OperationResult, error) {
	return nil, nil
}

func (r *RoleMutationResolver) UpdateRole(ctx context.Context, input models.UpdateRoleInput) (models.OperationResult, error) {
	return nil, nil
}

func (r *RoleMutationResolver) DeleteRole(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	return nil, nil
}
