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
	// Logger *gormlogger.GORMLogger\
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

func (r *TenantQueryResolver) Tenants(ctx context.Context) (models.OperationResult, error) {
	var tenants []models.Data
	page := 1
	per_page := 100
	logger.LogInfo("Fetching tenants from permit system")
	for page <= per_page {
		response, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants?page=%d&per_page=%d&include_total_count=true", page, per_page), nil)
		if err != nil {
			em := fmt.Sprintf("Error retrieving tenants from permit system: %v", err)
			logger.LogWarn(em)
			return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving tenants from permit system", em)), nil
		}

		pageData, ok := response["data"].([]interface{})
		if !ok {
			em := fmt.Sprintf("Error retrieving tenants from permit system: %v", err)
			logger.LogWarn(em)
			return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving tenants from permit system", em)), nil
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
			per_page = int(count)
		}
		page++
	}

	return utils.FormatSuccess(tenants)
}

// Tenant retrieves a single tenant by ID with its metadata
func (r *TenantQueryResolver) Tenant(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	if id == uuid.Nil {
		em := fmt.Sprint(ErrTenantIDRequired)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "Tenant ID is required", em)), nil
	}

	tenant, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants/%s", id), nil)
	if err != nil {
		em := fmt.Sprintf("Error retrieving tenant from permit system: %v", err)
		logger.LogWarn(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving tenant from permit system", em)), nil
	}

	data, err := r.extractTenantAttributes(tenant)
	if err != nil {
		em := fmt.Sprintf("Error retrieving tenant from permit system: %v", err)
		logger.LogWarn(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving tenant from permit system", em)), nil
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

// extractTenants processes raw tenant data
func (r *TenantQueryResolver) extractTenants(rawTenants map[string]interface{}) (models.OperationResult, error) {
	var tenants []models.Data

	for _, rawTenant := range rawTenants["data"].([]interface{}) {
		tenantMap, ok := rawTenant.(map[string]interface{})
		if !ok {
			em := "Error retrieving tenants from permit system"
			logger.LogWarn(em)
			return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving tenants from permit system", em)), nil
		}

		tenant, err := r.extractTenantAttributes(tenantMap)
		if err != nil {
			em := fmt.Sprintf("Error retrieving tenants from permit system: %v", err)
			logger.LogWarn(em)
			return utils.FormatError(utils.FormatErrorStruct("400", "Error retrieving tenants from permit system", em)), nil
		}

		tenants = append(tenants, tenant)
	}

	return utils.FormatSuccess(tenants)
}

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

	parentOrgID := uuid.Nil
	if attributes, ok := data["attributes"].(map[string]interface{}); ok {
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

		if parentOrgIDStr, ok := attributes["ParentID"].(string); ok {
			parentOrgID = uuid.MustParse(parentOrgIDStr)
		}
	}
	var parentOrg *dto.TenantResource

	if err := r.DB.Where(&dto.TenantResource{
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

func (r *TenantQueryResolver) ETenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	// Get tenant from permit
	tenant, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tenant from permit: %w", err)
	}

	data, err := r.extractTenantAttributes(tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tenant attributes: %w", err)
	}

	return data, err
}
