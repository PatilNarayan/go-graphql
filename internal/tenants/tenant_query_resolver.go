package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type TenantQueryResolver struct {
	DB *gorm.DB
}

// Tenants resolver for fetching all Tenants
func (r *TenantQueryResolver) AllTenants(ctx context.Context) ([]*models.Tenant, error) {
	var tenantResources []dto.TenantResource
	if err := r.DB.Find(&tenantResources).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tenants: %w", err)
	}

	tenants := make([]*models.Tenant, 0, len(tenantResources))
	for _, tenantResource := range tenantResources {
		var parentorg *dto.TenantResource
		if tenantResource.ParentResourceID != nil {
			if err := r.DB.Where(&dto.TenantResource{ResourceID: *tenantResource.ParentResourceID}).First(&parentorg).Error; err != nil {
				return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
			}
		}

		tenant := convertTenantToGraphQL(&tenantResource, parentorg)
		// Fetch associated metadata
		var tenantMetadata dto.TenantMetadata
		if err := r.DB.Where(&dto.TenantMetadata{ResourceID: tenantResource.ResourceID.String()}).First(&tenantMetadata).Error; err == nil {
			// Parse metadata and update tenant
			updateTenantWithMetadata(tenant, tenantMetadata)
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

// GetTenant resolver for fetching a single Tenant by ID
func (r *TenantQueryResolver) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	var tenantResource dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{ResourceID: id}).First(&tenantResource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch tenant: %w", err)
	}

	var parentOrg *dto.TenantResource
	if tenantResource.ParentResourceID != nil {
		if err := r.DB.Where(&dto.TenantResource{ResourceID: *tenantResource.ParentResourceID}).First(&parentOrg).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("parent organization not found: %w", err)
			}
			return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
		}
	}

	tenant := convertTenantToGraphQL(&tenantResource, parentOrg)

	// Fetch associated metadata
	var tenantMetadata dto.TenantMetadata
	if err := r.DB.Where("resource_id = ?", id).First(&tenantMetadata).Error; err == nil {
		// Parse metadata and update tenant
		updateTenantWithMetadata(tenant, tenantMetadata)
	}

	return tenant, nil
}

// Helper function to convert a TenantResource to a GraphQL Tenant model
func convertTenantToGraphQL(tenant *dto.TenantResource, parentOrg *dto.TenantResource) *models.Tenant {
	if tenant == nil {
		return nil
	}

	resp := &models.Tenant{
		ID:        tenant.ResourceID,
		Name:      tenant.Name,
		CreatedAt: tenant.CreatedAt.String(),
		CreatedBy: &tenant.CreatedBy,
		UpdatedAt: ptr.String(tenant.UpdatedAt.String()),
		UpdatedBy: &tenant.UpdatedBy,
	}

	if parentOrg != nil {
		resp.ParentOrg = &models.Root{
			ID:        parentOrg.ResourceID,
			Name:      parentOrg.Name,
			CreatedAt: parentOrg.CreatedAt.String(),
			UpdatedAt: ptr.String(parentOrg.UpdatedAt.String()),
			CreatedBy: &parentOrg.CreatedBy,
			UpdatedBy: &parentOrg.UpdatedBy,
		}
	}

	return resp
}

// Helper function to update Tenant with metadata
func updateTenantWithMetadata(tenant *models.Tenant, metadata dto.TenantMetadata) {
	if tenant == nil {
		return
	}

	var meta map[string]interface{}
	if err := json.Unmarshal(metadata.Metadata, &meta); err != nil {
		return
	}

	// Safely update fields from metadata
	if description, ok := meta["description"].(string); ok {
		tenant.Description = ptr.String(description)
	}

	if contactInfo, ok := meta["contactInfo"].(map[string]interface{}); ok {
		// Check if ContactInfo is non-nil
		if tenant.ContactInfo == nil {
			tenant.ContactInfo = &models.ContactInfo{}
		}

		if email, ok := contactInfo["email"].(string); ok {
			tenant.ContactInfo.Email = ptr.String(email)
		}

		if phoneNumber, ok := contactInfo["phoneNumber"].(string); ok {
			tenant.ContactInfo.PhoneNumber = ptr.String(phoneNumber)
		}

		// Handle Address
		if address, ok := contactInfo["address"].(map[string]interface{}); ok {
			// Initialize Address if it's nil
			if tenant.ContactInfo.Address == nil {
				tenant.ContactInfo.Address = &models.Address{}
			}

			if street, ok := address["street"].(string); ok {
				tenant.ContactInfo.Address.Street = ptr.String(street)
			}

			if city, ok := address["city"].(string); ok {
				tenant.ContactInfo.Address.City = ptr.String(city)
			}

			if state, ok := address["state"].(string); ok {
				tenant.ContactInfo.Address.State = ptr.String(state)
			}

			if zipCode, ok := address["zipCode"].(string); ok {
				tenant.ContactInfo.Address.ZipCode = ptr.String(zipCode)
			}

			if country, ok := address["country"].(string); ok {
				tenant.ContactInfo.Address.Country = ptr.String(country)
			}
		}
	}
}
