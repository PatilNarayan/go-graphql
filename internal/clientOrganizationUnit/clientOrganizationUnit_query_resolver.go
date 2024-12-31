package clientorganizationunit

import (
	"context"
	"go_graphql/gql/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ClientOrganizationUnitQueryResolver struct {
	DB *gorm.DB
}

func (r *ClientOrganizationUnitQueryResolver) GetClientOrganizationUnit(ctx context.Context, id uuid.UUID) (*models.ClientOrganizationUnit, error) {
	return nil, nil
}
func (r *ClientOrganizationUnitQueryResolver) AllClientOrganizationUnits(ctx context.Context) ([]*models.ClientOrganizationUnit, error) {
	return nil, nil
}
