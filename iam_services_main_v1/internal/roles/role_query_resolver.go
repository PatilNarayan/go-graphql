package roles

import (
	"context"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleQueryResolver handles role-related queries.
type RoleQueryResolver struct {
	DB *gorm.DB
}

func (r *RoleQueryResolver) Role(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	return nil, nil
}
func (r *RoleQueryResolver) Roles(ctx context.Context) (models.OperationResult, error) {
	return nil, nil
}
