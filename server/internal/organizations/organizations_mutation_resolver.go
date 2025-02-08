package organizations

import (
	"context"
	"go_graphql/internal/dto"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type OrganizationMutationResolver struct {
	DB *gorm.DB
}

// CreateOrganization resolver for adding a new organization
func (r *OrganizationMutationResolver) CreateOrganization(ctx context.Context, name string) (*dto.Organization, error) {
	organization := &dto.Organization{
		OrganizationID: uuid.Must(uuid.NewV4()),
		Name:           name,
	}
	if err := r.DB.Create(organization).Error; err != nil {
		return nil, err
	}
	return organization, nil
}
