package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"iam_services_main_v1/gormlogger"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"

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
	DB     *gorm.DB
	PC     *permit.PermitClient
	Logger *gormlogger.GORMLogger
}

// getTenantResourceType retrieves the resource type for tenants
func (r *TenantQueryResolver) getTenantResourceType() (*dto.Mst_ResourceTypes, error) {
	var resourceType dto.Mst_ResourceTypes
	// startTime := time.Now()
	err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error
	if err != nil {
		// r.Logger.Trace(context.Background(), startTime, func() (string, int64) {
		// 	return "SELECT * FROM mst_resource_types WHERE name = 'Tenant'", 0
		// }, err)
		return nil, fmt.Errorf("%w: %v", ErrResourceTypeNotFound, err)
	}

	// r.Logger.Trace(context.Background(), startTime, func() (string, int64) {
	// 	return "SELECT * FROM mst_resource_types WHERE name = 'Tenant'", 1
	// }, nil)
	return &resourceType, nil
}

func (r *TenantQueryResolver) Tenants(ctx context.Context) (models.OperationResult, error) {
	var tenants []models.Data
	page := 1
	pageCount := 1
	// r.Logger.Info(ctx, "Fetching tenants with pagination")

	for page <= pageCount {
		response, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants?page=%d", page), nil)
		if err != nil {
			// r.Logger.Error(ctx, "Failed to retrieve tenants from permit system", err)
			return utils.FormatError(&dto.CustomError{
				ErrorMessage: "Failed to retrieve tenants from permit system",
				ErrorCode:    "500",
				ErrorDetails: err.Error(),
			}), nil
		}

		pageData, ok := response["data"].([]interface{})
		if !ok {
			// r.Logger.Warn(ctx, "Invalid tenant data format received")
			return utils.FormatError(&dto.CustomError{
				ErrorMessage: "Invalid tenant data format",
				ErrorCode:    "400",
				ErrorDetails: "Failed to parse tenant data",
			}), nil
		}

		for _, rawTenant := range pageData {
			tenantMap, ok := rawTenant.(map[string]interface{})
			if !ok {
				// r.Logger.Warn(ctx, "Skipping invalid tenant format")
				continue
			}

			tenant, err := r.extractTenantAttributes(tenantMap)
			if err != nil {
				r.Logger.Error(ctx, "Failed to extract tenant attributes", err)
				continue
			}
			tenants = append(tenants, tenant)
		}

		if count, ok := response["page_count"].(float64); ok {
			pageCount = int(count)
		}
		page++
	}

	return utils.FormatSuccess(tenants)
}

// Tenant retrieves a single tenant by ID with its metadata
func (r *TenantQueryResolver) Tenant(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	if id == uuid.Nil {
		// r.Logger.Warn(ctx, "Tenant ID is required")
		return utils.FormatError(&dto.CustomError{
			ErrorMessage: "Tenant ID is required",
			ErrorCode:    "400",
			ErrorDetails: ErrTenantIDRequired.Error(),
		}), nil
	}

	tenant, err := r.PC.SendRequest(ctx, "GET", fmt.Sprintf("tenants/%s", id), nil)
	if err != nil {
		// r.Logger.Error(ctx, "Failed to retrieve tenant from permit system", err)
		return utils.FormatError(&dto.CustomError{
			ErrorMessage: "Failed to retrieve tenant from permit system",
			ErrorCode:    "500",
			ErrorDetails: err.Error(),
		}), nil
	}

	data, err := r.extractTenantAttributes(tenant)
	if err != nil {
		r.Logger.Error(ctx, "Failed to process tenant attributes", err)
		return utils.FormatError(&dto.CustomError{
			ErrorMessage: "Failed to process tenant attributes",
			ErrorCode:    "500",
			ErrorDetails: err.Error(),
		}), nil
	}

	// r.Logger.Info(ctx, "Successfully retrieved tenant data", data)
	return utils.FormatSuccess(data)
}

// enrichTenantWithMetadata fetches additional metadata for a tenant
func (r *TenantQueryResolver) enrichTenantWithMetadata(tenant *models.Tenant) error {
	if tenant == nil {
		return nil
	}

	// startTime := time.Now()
	var metadata dto.TenantMetadata
	err := r.DB.Where("resource_id = ?", tenant.ID).First(&metadata).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		// r.Logger.Trace(context.Background(), startTime, func() (string, int64) {
		// 	return "SELECT * FROM tenant_metadata WHERE resource_id = ?", 0
		// }, err)
		return fmt.Errorf("failed to fetch tenant metadata: %w", err)
	}

	// r.Logger.Trace(context.Background(), startTime, func() (string, int64) {
	// 	return "SELECT * FROM tenant_metadata WHERE resource_id = ?", 1
	// }, nil)

	var meta map[string]interface{}
	if err := json.Unmarshal(metadata.Metadata, &meta); err != nil {
		// r.Logger.Error(context.Background(), "Failed to unmarshal metadata", err)
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
			// r.Logger.Warn(context.Background(), "Invalid tenant format received")
			return utils.FormatError(&dto.CustomError{
				ErrorMessage: "Failed to parse tenant data",
				ErrorCode:    "400",
				ErrorDetails: "Invalid tenant format received",
			}), nil
		}

		tenant, err := r.extractTenantAttributes(tenantMap)
		if err != nil {
			// r.Logger.Error(context.Background(), "Failed to extract tenant attributes", err)
			return utils.FormatError(&dto.CustomError{
				ErrorMessage: "Failed to extract tenant attributes",
				ErrorCode:    "500",
				ErrorDetails: err.Error(),
			}), nil
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
