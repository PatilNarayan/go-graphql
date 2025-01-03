package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"go_graphql/logger"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type TenantQueryResolver struct {
	DB *gorm.DB
}

// Tenants resolver for fetching all Tenants
func (r *TenantQueryResolver) AllTenants(ctx context.Context) ([]*models.Tenant, error) {
	log := logger.Log.WithField("operation", "AllTenants")
	log.Info("Fetching all tenants")

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		log.WithError(err).Error("Failed to find resource type")
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	var tenantResources []dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{ResourceTypeID: resourceType.ResourceTypeID}).Find(&tenantResources).Error; err != nil {
		log.WithError(err).Error("Failed to fetch tenants")
		return nil, fmt.Errorf("failed to fetch tenants: %w", err)
	}

	log.WithField("count", len(tenantResources)).Info("Found tenants")
	tenants := make([]*models.Tenant, 0, len(tenantResources))

	for _, tr := range tenantResources {
		tenantLog := log.WithField("tenantID", tr.ResourceID)
		var parentOrg *dto.TenantResource

		if tr.ParentResourceID != nil {
			if err := r.DB.Where(&dto.TenantResource{ResourceID: *tr.ParentResourceID}).First(&parentOrg).Error; err != nil {
				tenantLog.WithError(err).Error("Failed to fetch parent organization")
				return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
			}
			tenantLog.WithField("parentID", *tr.ParentResourceID).Info("Found parent organization")
		}

		tenant := convertTenantToGraphQL(&tr, parentOrg)
		var metadata dto.TenantMetadata
		if err := r.DB.Where(&dto.TenantMetadata{ResourceID: tr.ResourceID.String()}).First(&metadata).Error; err == nil {
			updateTenantWithMetadata(tenant, metadata)
			tenantLog.Info("Updated tenant with metadata")
		} else {
			tenantLog.WithError(err).Warn("No metadata found for tenant")
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func (r *TenantQueryResolver) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	log := logger.Log.WithField("operation", "GetTenant").WithField("tenantID", id)
	log.Info("Fetching tenant")

	if id == uuid.Nil {
		return nil, fmt.Errorf("tenant ID is required")
	}

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		log.WithError(err).Error("Failed to find resource type")
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	var tenantResource dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{ResourceID: id, ResourceTypeID: resourceType.ResourceTypeID}).First(&tenantResource).Error; err != nil {
		log.WithError(err).Error("Failed to fetch tenant")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch tenant: %w", err)
	}

	var parentOrg *dto.TenantResource
	if tenantResource.ParentResourceID != nil {
		if err := r.DB.Where(&dto.TenantResource{ResourceID: *tenantResource.ParentResourceID}).First(&parentOrg).Error; err != nil {
			log.WithError(err).Error("Failed to fetch parent organization")
			return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
		}
		log.WithField("parentID", *tenantResource.ParentResourceID).Info("Found parent organization")
	}

	tenant := convertTenantToGraphQL(&tenantResource, parentOrg)
	var metadata dto.TenantMetadata
	if err := r.DB.Where("resource_id = ?", id).First(&metadata).Error; err == nil {
		updateTenantWithMetadata(tenant, metadata)
		log.Info("Updated tenant with metadata")
	} else {
		log.WithError(err).Warn("No metadata found for tenant")
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
