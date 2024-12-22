package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type TenantQueryResolver struct {
	DB *gorm.DB
}

// Tenants resolver for fetching all Tenants
func (r *TenantQueryResolver) Tenants(ctx context.Context) ([]*models.Tenant, error) {
	var Tenants []*dto.Tenant
	if err := r.DB.Find(&Tenants).Error; err != nil {
		return nil, err
	}
	var TenantsGraphQL []*models.Tenant
	for _, tenant := range Tenants {
		tenantGraphQL := convertTenantToGraphQL(tenant)
		TenantsGraphQL = append(TenantsGraphQL, tenantGraphQL)
	}
	return TenantsGraphQL, nil
}

// GetTenant resolver for fetching a single Tenant by ID
func (r *TenantQueryResolver) GetTenant(ctx context.Context, id string) (*models.Tenant, error) {
	if id == "" {
		return nil, errors.New("id cannot be nil")
	}

	var Tenant dto.Tenant
	if err := r.DB.Where(&dto.Tenant{ID: id}).First(&Tenant).Error; err != nil {
		return nil, err
	}
	return convertTenantToGraphQL(&Tenant), nil
}

func convertTenantToGraphQL(tenant *dto.Tenant) *models.Tenant {
	var tenantGraphQL models.Tenant

	// Ensure ID is always set
	tenantGraphQL.ID = tenant.ID

	// Ensure Name is always set
	tenantGraphQL.Name = tenant.Name

	// Ensure Description is always set (use empty string if nil)
	if tenant.Description == "" {
		tenantGraphQL.Description = ptr.String("") // Default to empty string if nil
	} else {
		tenantGraphQL.Description = &tenant.Description
	}

	// Ensure ParentOrgID is always set (use empty string if nil)
	if tenant.ParentOrgID == "" {
		tenantGraphQL.ParentOrgID = "" // Default to empty string if nil
	} else {
		tenantGraphQL.ParentOrgID = tenant.ParentOrgID
	}

	// Ensure Metadata is always set (use "{}" if nil or empty)
	if tenant.Metadata == nil {
		tenantGraphQL.Metadata = ptr.String("{}") // Default to "{}" if nil or empty
	} else {
		var temp interface{}
		err := json.Unmarshal([]byte(tenant.Metadata), &temp)
		if err != nil {
			tenantGraphQL.Metadata = ptr.String("{}") // Default to "{}" if invalid JSON
		} else {
			metaDataJson, err := json.Marshal(temp)
			if err != nil {
				tenantGraphQL.Metadata = ptr.String("{}") // Default to "{}" if re-marshalling fails
			} else {
				tenantGraphQL.Metadata = ptr.String(string(metaDataJson))
			}
		}
	}

	// Ensure ParentTenantID is always set (use empty string if nil)
	if tenant.ParentTenantID == "" {
		tenantGraphQL.ParentTenantID = nil // Default to nil if empty
	} else {
		tenantGraphQL.ParentTenantID = &tenant.ParentTenantID
	}

	// Ensure ResourceID is always set (use empty string if nil)
	if tenant.ResourceID == "" {
		tenantGraphQL.ResourceID = nil // Default to nil if empty
	} else {
		tenantGraphQL.ResourceID = &tenant.ResourceID
	}

	// Ensure CreatedAt is always set
	tenantGraphQL.CreatedAt = tenant.CreatedAt.String()

	// Ensure UpdatedAt is always set (use nil if empty)
	if tenant.UpdatedAt.IsZero() {
		tenantGraphQL.UpdatedAt = nil // Default to nil if empty
	} else {
		tenantGraphQL.UpdatedAt = ptr.String(tenant.UpdatedAt.String())
	}

	// Ensure UpdatedBy is always set (use empty string if nil)
	if tenant.UpdatedBy == "" {
		tenantGraphQL.UpdatedBy = nil // Default to nil if empty
	} else {
		tenantGraphQL.UpdatedBy = &tenant.UpdatedBy
	}

	// Ensure CreatedBy is always set (use empty string if nil)
	if tenant.CreatedBy == "" {
		tenantGraphQL.CreatedBy = nil // Default to nil if empty
	} else {
		tenantGraphQL.CreatedBy = &tenant.CreatedBy
	}

	return &tenantGraphQL
}
