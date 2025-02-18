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
	"iam_services_main_v1/internal/utils"
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
	if err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResourceTypeNotFound, err)
	}
	return &resourceType, nil
}

// AllTenants retrieves all tenants with their associated metadata
func (r *TenantQueryResolver) Tenants(ctx context.Context) (models.OperationResult, error) {
	// Get all tenants from permit
	allTenants, err := r.PC.SendRequest(ctx, "GET", "tenants", nil)
	if err != nil {
		errMsg := dto.CustomError{
			ErrorMessage: "Failed to retrieve tenants from permit system",
			ErrorCode:    "500",
			ErrorDetails: err.Error(),
		}
		return utils.FormatError(&errMsg), nil
	}

	fmt.Println(r.extractTenants(allTenants.([]interface{})))
	response, err := r.extractTenants(allTenants.([]interface{}))
	if err != nil {
		errMsg := dto.CustomError{
			ErrorMessage: "Failed to process tenant data",
			ErrorCode:    "500",
			ErrorDetails: err.Error(),
		}
		return utils.FormatError(&errMsg), nil
	}

	return response, nil
}

// GetTenant retrieves a single tenant by ID with its metadata
func (r *TenantQueryResolver) Tenant(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	if id == uuid.Nil {
		errMsg := dto.CustomError{
			ErrorMessage: "Tenant ID is required",
			ErrorCode:    "400",
			ErrorDetails: ErrTenantIDRequired.Error(),
		}
		return utils.FormatError(&errMsg), nil
	}

	// Get tenant from permit
	tenant, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants/%s", id), nil)
	if err != nil {
		errMsg := dto.CustomError{
			ErrorMessage: "Failed to retrieve tenant from permit system",
			ErrorCode:    "500",
			ErrorDetails: err.Error(),
		}
		return utils.FormatError(&errMsg), nil
	}

	data, err := r.extractTenantAttributes(tenant.(map[string]interface{}))
	if err != nil {
		errMsg := dto.CustomError{
			ErrorMessage: "Failed to process tenant attributes",
			ErrorCode:    "500",
			ErrorDetails: err.Error(),
		}
		return utils.FormatError(&errMsg), nil
	}

	// Return success response with tenant
	return utils.FormatSuccess(data)
}

// processTenantResources processes a slice of tenant resources and returns GraphQL tenant models
func (r *TenantQueryResolver) processTenantResources(resources []dto.TenantResource) ([]*models.Tenant, error) {
	tenants := make([]*models.Tenant, 0, len(resources))

	for _, tr := range resources {
		var parentOrg *dto.TenantResource
		if tr.ParentResourceID != nil {
			if err := r.DB.Where(&dto.TenantResource{
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

	return nil, nil
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

func (r *TenantQueryResolver) extractTenants(rawTenants []interface{}) (models.OperationResult, error) {
	var tenants []models.Data

	for _, rawTenant := range rawTenants {
		tenantMap, ok := rawTenant.(map[string]interface{})
		if !ok {
			errMsg := dto.CustomError{
				ErrorMessage: "Failed to parse tenant data",
				ErrorCode:    "400", // Changed from PARSING_ERROR to 400 for consistency
				ErrorDetails: "Invalid tenant format received",
			}
			return utils.FormatError(&errMsg), nil
		}

		tenant, err := r.extractTenantAttributes(tenantMap)
		if err != nil {
			errMsg := dto.CustomError{
				ErrorMessage: "Failed to extract tenant attributes",
				ErrorCode:    "500",
				ErrorDetails: err.Error(),
			}
			return utils.FormatError(&errMsg), nil
		}
		// Add to tenants slice
		tenants = append(tenants, tenant)
	}

	// Return success response with tenants
	return utils.FormatSuccess(tenants)
}

func (r *TenantQueryResolver) ETenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	// Get tenant from permit
	tenant, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tenant from permit: %w", err)
	}

	data, err := r.extractTenantAttributes(tenant.(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("failed to extract tenant attributes: %w", err)
	}

	return data, err
}
