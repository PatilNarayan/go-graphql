package organizations

import (
	"context"
	"go_graphql/internal/dto"

	"gorm.io/gorm"
)

type OrganizationFieldResolver struct {
	DB *gorm.DB
}

type OrganizationQueryResolver struct {
	DB *gorm.DB
}

type OrganizationMutationResolver struct {
	DB *gorm.DB
}

// Organizations implements generated.QueryResolver.
func (q *OrganizationQueryResolver) Organizations(ctx context.Context) ([]dto.Organization, error) {
	panic("unimplemented")
}

func (q *OrganizationQueryResolver) GetOrganization(ctx context.Context, id string) (dto.Organization, error) {
	panic("unimplemented")
}

func (r *OrganizationMutationResolver) CreateOrganization(ctx context.Context, name string) (dto.Organization, error) {
	panic("unimplemented")
}
