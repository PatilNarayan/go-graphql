package tenants

import (
	"context"
	"go_graphql/internal/dto"
	"time"

	"gorm.io/gorm"
)

type TenantFieldResolver struct {
	DB *gorm.DB
}

// ContactInfo implements generated.TenantResolver.
func (r *TenantFieldResolver) ContactInfo(ctx context.Context, obj *dto.Tenant) (*string, error) {
	return nil, nil
}

// ParentOrg implements generated.TenantResolver.
func (r *TenantFieldResolver) ParentOrg(ctx context.Context, obj *dto.Tenant) (dto.Organization, error) {
	return nil, nil
}

// ID resolves the ID field on the Article type
func (r *TenantFieldResolver) ID(ctx context.Context, obj *dto.Tenant) (string, error) {
	return "obj", nil
}

// Implement the CreatedAt method for Tenant
func (r *TenantFieldResolver) CreatedAt(ctx context.Context, obj *dto.Tenant) (string, error) {
	// Assuming the "createdAt" field is a time.Time object, format it as a string
	createdAtStr := obj.CreatedAt.Format(time.RFC3339)
	return createdAtStr, nil
}

// Implement the UpdatedAt method for Tenant
func (r *TenantFieldResolver) UpdatedAt(ctx context.Context, obj *dto.Tenant) (*string, error) {
	// Assuming the "createdAt" field is a time.Time object, format it as a string
	UpdatedAtStr := obj.UpdatedAt.Format(time.RFC3339)
	return &UpdatedAtStr, nil
}

func (r *TenantFieldResolver) Organization(ctx context.Context, obj *dto.Tenant) (*dto.Organization, error) {
	var organization dto.Organization
	if err := r.DB.First(&organization, obj.ParentOrgID).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}
