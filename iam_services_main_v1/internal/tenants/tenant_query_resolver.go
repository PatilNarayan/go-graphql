package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

var (
	ErrTenantIDRequired     = errors.New("tenant ID is required")
	ErrResourceTypeNotFound = errors.New("resource type not found")
	ErrTenantNotFound       = errors.New("tenant not found")
	ErrParentOrgNotFound    = errors.New("failed to fetch parent organization")
)

// TenantQueryResolver handles tenant-related GraphQL queries
type TenantQueryResolver struct {
	DB *gorm.DB
	PC *permit.PermitClient
}

// getTenantResourceType retrieves the resource type for tenants
func (r *TenantQueryResolver) getTenantResourceType() (*dto.Mst_ResourceTypes, error) {
	var resourceType dto.Mst_ResourceTypes
	err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResourceTypeNotFound, err)
	}
	return &resourceType, nil
}

// Tenants retrieves a list of tenants from the permit system
func (r *TenantQueryResolver) Tenants(ctx context.Context) (models.OperationResult, error) {
	var tenants []models.Data
	page := 1
	perPage := 100

	logger.LogInfo("Fetching tenants from permit system")

	for page <= perPage {
		response, err := r.fetchTenantsFromPermit(ctx, page, perPage)
		if err != nil {
			return r.handleError("400", "Error retrieving tenants from permit system", err)
		}

		pageData, ok := response["data"].([]interface{})
		if !ok {
			return r.handleError("400", "Error retrieving tenants from permit system", errors.New("invalid data format"))
		}

		for _, rawTenant := range pageData {
			tenantMap, ok := rawTenant.(map[string]interface{})
			if !ok {
				continue
			}

			tenant, err := r.extractTenantAttributes(tenantMap)
			if err != nil {
				continue
			}
			tenants = append(tenants, tenant)
		}

		if count, ok := response["page_count"].(float64); ok {
			perPage = int(count)
		}
		page++
	}

	return utils.FormatSuccess(tenants)
}

// Tenant retrieves a single tenant by ID with its metadata
func (r *TenantQueryResolver) Tenant(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	if id == uuid.Nil {
		return r.handleError("400", "Tenant ID is required", ErrTenantIDRequired)
	}

	tenant, err := r.fetchTenantFromPermit(ctx, id)
	if err != nil {
		return r.handleError("400", "Error retrieving tenant from permit system", err)
	}

	data, err := r.extractTenantAttributes(tenant)
	if err != nil {
		return r.handleError("400", "Error retrieving tenant from permit system", err)
	}

	var modelsData []models.Data
	modelsData = append(modelsData, *data)
	return utils.FormatSuccess(modelsData)
}

// enrichTenantWithMetadata fetches additional metadata for a tenant
func (r *TenantQueryResolver) enrichTenantWithMetadata(tenant *models.Tenant) error {
	if tenant == nil {
		return nil
	}

	var metadata dto.TenantMetadata
	err := r.DB.Where("resource_id = ?", tenant.ID).First(&metadata).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("failed to fetch tenant metadata: %w", err)
	}

	var meta map[string]interface{}
	if err := json.Unmarshal(metadata.Metadata, &meta); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if description, ok := meta["description"].(string); ok {
		tenant.Description = ptr.String(description)
	}

	return nil
}

// extractTenantAttributes processes raw tenant data into a Tenant model
func (r *TenantQueryResolver) extractTenantAttributes(data map[string]interface{}) (*models.Tenant, error) {
	tenant := &models.Tenant{}

	if id, ok := data["key"].(string); ok {
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

	if attributes, ok := data["attributes"].(map[string]interface{}); ok {
		tenant = r.extractAttributesFromMap(tenant, attributes)
	}

	parentOrgID := uuid.Nil
	if tenant.ParentOrg != nil {
		parentOrgID = tenant.ID
	}

	parentOrg, err := r.fetchParentOrg(parentOrgID)
	if err != nil {
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

// extractAttributesFromMap extracts attributes from a map and populates the tenant model
func (r *TenantQueryResolver) extractAttributesFromMap(tenant *models.Tenant, attributes map[string]interface{}) *models.Tenant {
	if attrName, ok := attributes["Name"].(string); ok {
		tenant.Name = attrName
	}

	if description, ok := attributes["Description"].(string); ok {
		tenant.Description = &description
	}

	if createdBy, ok := attributes["created_by"].(string); ok {
		tenant.CreatedBy = uuid.MustParse(createdBy)
	}

	if updatedBy, ok := attributes["updated_by"].(string); ok {
		tenant.UpdatedBy = uuid.MustParse(updatedBy)
	}

	if contactInfo, ok := attributes["ContactInfo"].(map[string]interface{}); ok {
		tenant.ContactInfo = buildContactInfo(contactInfo)
	}

	return tenant
}

// fetchTenantsFromPermit fetches tenants from the permit system
func (r *TenantQueryResolver) fetchTenantsFromPermit(ctx context.Context, page, perPage int) (map[string]interface{}, error) {
	response, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants?page=%d&per_page=%d&include_total_count=true", page, perPage), nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tenants from permit system: %w", err)
	}
	return response, nil
}

// fetchTenantFromPermit fetches a single tenant from the permit system
func (r *TenantQueryResolver) fetchTenantFromPermit(ctx context.Context, id uuid.UUID) (map[string]interface{}, error) {
	tenant, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tenant from permit system: %w", err)
	}
	return tenant, nil
}

// fetchParentOrg fetches the parent organization for a tenant
func (r *TenantQueryResolver) fetchParentOrg(parentOrgID uuid.UUID) (*dto.TenantResource, error) {
	var parentOrg dto.TenantResource
	err := r.DB.Where(&dto.TenantResource{
		ResourceID: parentOrgID,
	}).First(&parentOrg).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
	}
	return &parentOrg, nil
}

// handleError logs and formats an error response
func (r *TenantQueryResolver) handleError(code, message string, err error) (models.OperationResult, error) {
	em := fmt.Sprintf("%s: %v", message, err)
	logger.LogError(em)
	return utils.FormatError(utils.FormatErrorStruct(code, message, em)), nil
}

func TenantDataPermit(ctx context.Context, r *TenantQueryResolver, id uuid.UUID) (*models.Tenant, error) {
	tenant, err := r.fetchTenantFromPermit(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tenant from permit system: %w", err)
	}

	data, err := r.extractTenantAttributes(tenant)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tenant from permit system: %w", err)
	}

	return data, nil
}
