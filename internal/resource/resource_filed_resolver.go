package resource

import (
	"context"
	"go_graphql/gql/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ResourceFieldResolver struct {
	DB *gorm.DB
}

type ResourceQueryResolver struct {
	DB *gorm.DB
}

func (r *ResourceQueryResolver) GetResource(ctx context.Context, id uuid.UUID) (models.Resource, error)
func (r *ResourceQueryResolver) AllResources(ctx context.Context) ([]models.Resource, error)
