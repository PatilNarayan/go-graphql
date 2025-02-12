package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"

	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/permit"
)

var (
	ErrTenantIDRequired     = errors.New("tenant ID is required")
	ErrResourceTypeNotFound = errors.New("resource type not found")
	ErrTenantNotFound       = errors.New("tenant not found")
	ErrParentOrgNotFound    = errors.New("failed to fetch parent organization")
)

// TenantQueryResolver handles tenant-related GraphQL queries
type TenantQueryResolver struct {
	DB           *gorm.DB
	PermitClient *permit.PermitClient
}

// getTenantResourceType retrieves the resource type for tenants
func (r *TenantQueryResolver) getTenantResourceType() (*dto.Mst_ResourceTypes, error) {
	var resourceType dto.Mst_ResourceTypes
	if err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResourceTypeNotFound, err)
	}
	return &resourceType, nil
}

// AllTenants retrieves all tenants with their associated metadata
func (r *TenantQueryResolver) AllTenants(ctx context.Context) ([]*models.Tenant, error) {

	// get all tenants from permit
	allTenants, err := r.PermitClient.SendRequest(ctx, "GET", "tenants", nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(allTenants)

	return r.extractTenants(allTenants.([]interface{}))

	// resourceType, err := r.getTenantResourceType()
	// if err != nil {
	// 	return nil, err
	// }

	// var tenantResources []dto.TenantResources
	// if err := r.DB.Where(&dto.TenantResources{
	// 	ResourceTypeID: resourceType.ResourceTypeID,
	// }).Find(&tenantResources).Error; err != nil {
	// 	return nil, fmt.Errorf("failed to fetch tenants: %w", err)
	// }

	// return r.processTenantResources(tenantResources)
}

// GetTenant retrieves a single tenant by ID with its metadata
func (r *TenantQueryResolver) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	if id == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	// get tenant from permit
	tenant, err := r.PermitClient.SendRequest(ctx, "GET", fmt.Sprintf("tenants/%s", id), nil)
	if err != nil {
		return nil, err
	}

	return r.extractTenantAttributes(tenant.(map[string]interface{}))

	// resourceType, err := r.getTenantResourceType()
	// if err != nil {
	// 	return nil, err
	// }

	// var tenantResource dto.TenantResources
	// if err := r.DB.Where(&dto.TenantResources{
	// 	ResourceID:     id,
	// 	ResourceTypeID: resourceType.ResourceTypeID,
	// }).First(&tenantResource).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, ErrTenantNotFound
	// 	}
	// 	return nil, fmt.Errorf("failed to fetch tenant: %w", err)
	// }

	// tenants, err := r.processTenantResources([]dto.TenantResources{tenantResource})
	// if err != nil {
	// 	return nil, err
	// }
	// if len(tenants) == 0 {
	// 	return nil, ErrTenantNotFound
	// }
	// return tenants[0], nil
}

// processTenantResources processes a slice of tenant resources and returns GraphQL tenant models
func (r *TenantQueryResolver) processTenantResources(resources []dto.TenantResources) ([]*models.Tenant, error) {
	tenants := make([]*models.Tenant, 0, len(resources))

	for _, tr := range resources {
		var parentOrg *dto.TenantResources
		if tr.ParentResourceID != nil {
			if err := r.DB.Where(&dto.TenantResources{
				ResourceID: *tr.ParentResourceID,
			}).First(&parentOrg).Error; err != nil {
				return nil, fmt.Errorf("%w: %v", ErrParentOrgNotFound, err)
			}
		}

		tenant := convertTenantToGraphQL(&tr, parentOrg)
		if err := r.enrichTenantWithMetadata(tenant); err != nil {
			// Log error but continue processing other tenants
			continue
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func (r *TenantQueryResolver) enrichTenantWithMetadata(tenant *models.Tenant) error {
	if tenant == nil {
		return nil
	}

	var metadata dto.TenantMetadata
	if err := r.DB.Where("resource_id = ?", tenant.ID).First(&metadata).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Not finding metadata is not an error
		}
		return fmt.Errorf("failed to fetch tenant metadata: %w", err)
	}

	var meta map[string]interface{}
	if err := json.Unmarshal(metadata.Metadata, &meta); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Update description
	if description, ok := meta["description"].(string); ok {
		tenant.Description = ptr.String(description)
	}

	// Update contact info
	if contactData, ok := meta["contactInfo"].(map[string]interface{}); ok {
		tenant.ContactInfo = buildContactInfo(contactData)
	}

	return nil
}

func (r *TenantQueryResolver) extractTenantAttributes(data map[string]interface{}) (*models.Tenant, error) {
	tenant := &models.Tenant{}

	if id, ok := data["id"].(string); ok {
		tenant.ID = uuid.MustParse(id)
	}

	if name, ok := data["name"].(string); ok {
		tenant.Name = name
	}

	if createdAt, ok := data["created_at"].(string); ok {
		tenant.CreatedAt = createdAt
	}

	if updatedAt, ok := data["updated_at"].(string); ok {
		tenant.UpdatedAt = updatedAt
	}

	parentOrgID := uuid.Nil
	if attributes, ok := data["attributes"].(map[string]interface{}); ok {
		if attrName, ok := attributes["name"].(string); ok {
			tenant.Name = attrName
		}

		if description, ok := attributes["description"].(string); ok {
			tenant.Description = &description
		}

		if createdBy, ok := attributes["createdBy"].(string); ok {
			tenant.CreatedBy = uuid.MustParse(createdBy)
		}

		if updatedBy, ok := attributes["updatedBy"].(string); ok {
			tenant.UpdatedBy = uuid.MustParse(updatedBy)
		}

		if contactInfo, ok := attributes["contactInfo"].(map[string]interface{}); ok {
			tenant.ContactInfo = buildContactInfo(contactInfo)
		}

		if parentOrgIDStr, ok := attributes["parentOrgId"].(string); ok {
			parentOrgID = uuid.MustParse(parentOrgIDStr)
		}
	}
	var parentOrg *dto.TenantResources

	if err := r.DB.Where(&dto.TenantResources{
		ResourceID: parentOrgID,
	}).First(&parentOrg).Error; err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParentOrgNotFound, err)
	}

	if parentOrg != nil {
		tenant.ParentOrg = &models.Root{
			ID:        parentOrg.ResourceID,
			Name:      parentOrg.Name,
			CreatedAt: parentOrg.CreatedAt.String(),
			UpdatedAt: parentOrg.UpdatedAt.String(),
			CreatedBy: parentOrg.CreatedBy,
			UpdatedBy: parentOrg.UpdatedBy,
		}
	}

	return tenant, nil
}

func (r *TenantQueryResolver) extractTenants(data []interface{}) ([]*models.Tenant, error) {
	var tenants []*models.Tenant

	for _, item := range data {
		tenantMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid tenant data format")
		}

		tenant, err := r.extractTenantAttributes(tenantMap)
		if err != nil {
			return nil, err
		}

		tenants = append(tenants, tenant)
	}

	return tenants, nil
}
