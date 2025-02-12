package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"

	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type TenantFieldResolver struct {
	DB *gorm.DB
}

// ContactInfo resolves the contactInfo field for a tenant
func (t *TenantFieldResolver) ContactInfo(ctx context.Context, obj *models.Tenant) (*models.ContactInfo, error) {
	if obj == nil {
		return nil, errors.New("tenant object is nil")
	}

	var metadata dto.TenantMetadata
	if err := t.DB.Where("resource_id = ?", obj.ID).First(&metadata).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no metadata exists
		}
		return nil, fmt.Errorf("failed to fetch tenant metadata: %w", err)
	}

	var meta map[string]interface{}
	if err := json.Unmarshal(metadata.Metadata, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	contactInfo, ok := meta["contactInfo"].(map[string]interface{})
	if !ok {
		return nil, nil // Return nil if no contact info exists
	}

	return buildContactInfo(contactInfo), nil
}

// File: shared_utils.go

// buildContactInfo creates a ContactInfo model from raw contact data
func buildContactInfo(data map[string]interface{}) *models.ContactInfo {
	info := &models.ContactInfo{}

	if email, ok := data["email"].(string); ok {
		info.Email = ptr.String(email)
	}
	if phone, ok := data["phoneNumber"].(string); ok {
		info.PhoneNumber = ptr.String(phone)
	}

	if addressData, ok := data["address"].(map[string]interface{}); ok {
		info.Address = buildAddress(addressData)
	}

	return info
}

// buildAddress creates an Address model from raw address data
func buildAddress(data map[string]interface{}) *models.Address {
	addr := &models.Address{}

	if street, ok := data["street"].(string); ok {
		addr.Street = ptr.String(street)
	}
	if city, ok := data["city"].(string); ok {
		addr.City = ptr.String(city)
	}
	if state, ok := data["state"].(string); ok {
		addr.State = ptr.String(state)
	}
	if zipCode, ok := data["zipCode"].(string); ok {
		addr.ZipCode = ptr.String(zipCode)
	}
	if country, ok := data["country"].(string); ok {
		addr.Country = ptr.String(country)
	}

	return addr
}

// Helper functions (can be in either file or separate utils file)
func convertTenantToGraphQL(tenant *dto.TenantResources, parentOrg *dto.TenantResources) *models.Tenant {
	if tenant == nil {
		return nil
	}

	resp := &models.Tenant{
		ID:        tenant.ResourceID,
		Name:      tenant.Name,
		CreatedAt: tenant.CreatedAt.String(),
		CreatedBy: tenant.CreatedBy,
		UpdatedAt: tenant.UpdatedAt.String(),
		UpdatedBy: tenant.UpdatedBy,
	}

	if parentOrg != nil {
		resp.ParentOrg = &models.Root{
			ID:        parentOrg.ResourceID,
			Name:      parentOrg.Name,
			CreatedAt: parentOrg.CreatedAt.String(),
			UpdatedAt: parentOrg.UpdatedAt.String(),
			CreatedBy: parentOrg.CreatedBy,
			UpdatedBy: parentOrg.UpdatedBy,
		}
	}

	return resp
}
