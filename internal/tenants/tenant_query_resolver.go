package tenants

import (
	"context"
	"errors"
	"go_graphql/internal/dto"

	"gorm.io/gorm"
)

type TenantQueryResolver struct {
	DB *gorm.DB
}

// Tenants resolver for fetching all Tenants
func (r *TenantQueryResolver) Tenants(ctx context.Context) ([]*dto.Tenant, error) {
	var Tenants []*dto.Tenant
	if err := r.DB.Find(&Tenants).Error; err != nil {
		return nil, err
	}
	return Tenants, nil
}

// GetTenant resolver for fetching a single Tenant by ID
func (r *TenantQueryResolver) GetTenant(ctx context.Context, id *string) (*dto.Tenant, error) {
	if id == nil {
		return nil, errors.New("id cannot be nil")
	}

	var Tenant dto.Tenant
	if err := r.DB.Where(&dto.Tenant{ID: *id}).First(&Tenant).Error; err != nil {
		return nil, err
	}
	return &Tenant, nil
}
