package clientorganizationunits

import (
	"context"
	"iam_services_main_v1/gql/models"

	"gorm.io/gorm"
)

type ClientOrganizationUnitFieldResolver struct {
	DB *gorm.DB
}

func (c *ClientOrganizationUnitFieldResolver) Tenant(ctx context.Context, obj *models.ClientOrganizationUnit) (*models.Tenant, error) {
	return nil, nil
}
func (c *ClientOrganizationUnitFieldResolver) ParentOrg(ctx context.Context, obj *models.ClientOrganizationUnit) (models.Organization, error) {
	return nil, nil
}
