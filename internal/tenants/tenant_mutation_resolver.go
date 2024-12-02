package tenants

import (
	"context"
	"fmt"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"go_graphql/permit"

	"gorm.io/gorm"
)

type TenantMutationResolver struct {
	DB *gorm.DB
}

// CreateTenant resolver for adding a new Tenant
func (r *TenantMutationResolver) CreateTenant(ctx context.Context, input models.TenantInput) (*dto.Tenant, error) {

	tenant, err := permit.CreateTenant(input.Name)
	if err != nil {
		return nil, err
	}

	Tenant := &dto.Tenant{Name: input.Name, ParentOrgID: input.ParentOrgID, RemoteTenantID: tenant.Key}
	if input.Description != nil {
		Tenant.Description = *input.Description
	}
	if input.ContactInfoID != "" {
		Tenant.ContactInfoID = input.ContactInfoID
	}
	if err := r.DB.Create(Tenant).Error; err != nil {
		return nil, err
	}
	return Tenant, nil
}

// UpdateTenant resolver for updating a Tenant
func (r *TenantMutationResolver) UpdateTenant(ctx context.Context, id string, input models.TenantInput) (*dto.Tenant, error) {
	var Tenant *dto.Tenant
	if err := r.DB.First(&Tenant, id).Error; err != nil {
		return nil, err
	}

	if Tenant != nil {
		_, err := permit.UpdateTenant(Tenant.RemoteTenantID, input.Name)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Tenant not found")
	}

	Tenant.Name = input.Name
	Tenant.ParentOrgID = input.ParentOrgID

	if err := r.DB.Save(&Tenant).Error; err != nil {
		return nil, err
	}
	return Tenant, nil
}

// DeleteTenant resolver for deleting a Tenant
func (r *TenantMutationResolver) DeleteTenant(ctx context.Context, id string) (bool, error) {
	var tenant *dto.Tenant
	if err := r.DB.First(&tenant, id).Error; err != nil {
		return false, err
	}
	if tenant == nil {
		return false, fmt.Errorf("Tenant not found")
	}
	_, err := permit.DeleteTenant(tenant.RemoteTenantID)
	if err != nil {
		return false, err
	}

	if err := r.DB.Delete(&tenant).Error; err != nil {
		return false, err
	}
	return true, nil
}
