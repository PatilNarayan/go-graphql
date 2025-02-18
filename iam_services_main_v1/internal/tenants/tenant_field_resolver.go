package tenants

import (
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"

	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type TenantFieldResolver struct {
	DB *gorm.DB
}

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
func convertTenantToGraphQL(tenant *dto.TenantResource, parentOrg *dto.TenantResource) *models.Tenant {
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
