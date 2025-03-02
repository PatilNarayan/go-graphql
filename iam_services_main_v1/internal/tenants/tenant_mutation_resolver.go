package tenants

import (
	"context"
	"encoding/json"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/dao"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/middlewares"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/internal/validations"
	"iam_services_main_v1/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantMutationResolver struct {
	DB           *gorm.DB
	PermitClient *permit.PermitClient
}

// CreateTenant resolver for adding a new Tenant
func (t *TenantMutationResolver) CreateTenant(ctx context.Context, input models.CreateTenantInput) (models.OperationResult, error) {

	newTenantID := uuid.New()
	// Extract gin.Context from GraphQL context
	ginCtx, ok := ctx.Value(middlewares.GinContextKey).(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unable to get gin context")
	}
	UserID := ginCtx.MustGet("userID").(string)
	userUUID := uuid.MustParse(UserID)

	inputMap := helpers.StructToMap(input)

	if err := t.createTenantInPermit(ctx, input.Name, newTenantID, inputMap); err != nil {
		return t.handleError("500", "Error creating tenant in permit system", err)
	}

	resourceType, err := dao.GetResourceTypeByName("Tenant")
	if err != nil {
		return t.handleError("500", "Error getting resource type", err)
	}

	if err := t.createResourceInstanceInPermit(ctx, input.ID, resourceType.ResourceTypeID, newTenantID, input); err != nil {
		return t.handleError("500", "Error creating resource instance of tenant in permit system", err)
	}

	tenantResource, err := t.createTenantResource(input.Name, newTenantID, *input.ParentID, userUUID, uuid.Nil)
	if err != nil {
		return t.handleError("500", "Error creating tenant resource", err)
	}

	if err := t.createTenantMetadata(tenantResource.ResourceID, input.Description, input.ContactInfo, userUUID); err != nil {
		return t.handleError("500", "Error creating tenant metadata", err)
	}

	return t.getTenantResponse(ctx, newTenantID)
}

// UpdateTenant resolver for updating a Tenant
func (t *TenantMutationResolver) UpdateTenant(ctx context.Context, input models.UpdateTenantInput) (models.OperationResult, error) {

	tenant, err := TenantDataPermit(ctx, &TenantQueryResolver{DB: t.DB, PC: t.PermitClient}, input.ID)
	if err != nil {
		return t.handleError("500", "Error retrieving tenant from permit system", err)
	}
	inputMap := helpers.StructToMap(input)
	inputMap["created_by"] = tenant.CreatedBy
	inputMap["updated_by"] = tenant.UpdatedBy

	if err := t.updateTenantInPermit(ctx, input.ID, *input.Name, inputMap); err != nil {
		return t.handleError("500", "Error updating tenant in permit system", err)
	}

	if err := t.updateTenantResource(input.ID, input.Name, input.ParentID, tenant.UpdatedBy); err != nil {
		return t.handleError("500", "Error updating tenant resource", err)
	}

	if err := t.updateMetadata(input.ID, input.Description, input.ContactInfo, tenant.UpdatedBy); err != nil {
		return t.handleError("500", "Error updating tenant metadata", err)
	}

	return t.getTenantResponse(ctx, input.ID)
}

// DeleteTenant resolver for deleting a Tenant
func (t *TenantMutationResolver) DeleteTenant(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := t.deleteTenantInPermit(ctx, input.ID); err != nil {
		tx.Rollback()
		return t.handleError("500", "Error deleting tenant in permit system", err)
	}

	if err := t.updateTenantMetadata(tx, input.ID); err != nil {
		tx.Rollback()
		return t.handleError("500", "Error updating tenant metadata", err)
	}

	if err := t.updateTenantResourceStatus(tx, input.ID); err != nil {
		tx.Rollback()
		return t.handleError("500", "Error updating tenant resource", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return t.handleError("500", "Error committing transaction", err)
	}

	return utils.FormatSuccess([]models.Data{})
}

// Helper functions

func (t *TenantMutationResolver) handleError(code, message string, err error) (models.OperationResult, error) {
	em := fmt.Sprintf("%s: %v", message, err)
	logger.LogError(em)
	return utils.FormatError(utils.FormatErrorStruct(code, message, em)), nil
}

func (t *TenantMutationResolver) validateParentOrg(parentOrgID uuid.UUID) (*uuid.UUID, error) {
	if parentOrgID == uuid.Nil {
		return nil, fmt.Errorf("parent organization ID is required")
	}

	resourceType, err := dao.GetResourceTypeByName("Root")
	if err != nil {
		return nil, fmt.Errorf("failed to get resource type IDs: %w", err)
	}

	var parentOrg dto.TenantResource
	if err := t.DB.Where(
		"resource_id = ? AND resource_type_id in (?) AND row_status = 1",
		parentOrgID, resourceType.ResourceTypeID,
	).First(&parentOrg).Error; err != nil {
		return nil, fmt.Errorf("parent organization not found: %w", err)
	}

	return &parentOrg.ResourceID, nil
}

func (t *TenantMutationResolver) createTenantInPermit(ctx context.Context, name string, tenantID uuid.UUID, attributes map[string]interface{}) error {
	_, err := t.PermitClient.SendRequest(ctx, "POST", "tenants", map[string]interface{}{
		"name":       name,
		"key":        tenantID,
		"attributes": attributes,
	})
	return err
}

func (t *TenantMutationResolver) createResourceInstanceInPermit(ctx context.Context, resourceID, resourceTypeID uuid.UUID, tenantID uuid.UUID, input models.CreateTenantInput) error {
	_, err := t.PermitClient.SendRequest(ctx, "POST", "resource_instances", map[string]interface{}{
		"key":        resourceID,
		"resource":   resourceTypeID,
		"tenant":     tenantID,
		"attributes": input,
	})
	return err
}

func (t *TenantMutationResolver) createTenantResource(name string, resourceID, parentID uuid.UUID, userID, tenantID uuid.UUID) (*dto.TenantResource, error) {
	var resourceType dto.Mst_ResourceTypes
	if err := t.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	tenant := &dto.TenantResource{
		ResourceID:       resourceID,
		Name:             name,
		CreatedBy:        userID,
		UpdatedBy:        userID,
		CreatedAt:        time.Now(),
		ResourceTypeID:   resourceType.ResourceTypeID,
		ParentResourceID: &parentID,
	}
	if tenantID != uuid.Nil {
		tenant.TenantID = &tenantID
	}
	if err := t.DB.Create(tenant).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant resource: %w", err)
	}

	return tenant, nil
}

func (t *TenantMutationResolver) createTenantMetadata(resourceID uuid.UUID, description *string, contactInfo *models.ContactInfoInput, userID uuid.UUID) error {
	metadata := map[string]interface{}{
		"description": description,
		"contactInfo": contactInfo,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	tenantMetadata := &dto.TenantMetadata{
		ResourceID: resourceID,
		Metadata:   metadataJSON,
		CreatedBy:  userID,
		CreatedAt:  time.Now(),
		UpdatedBy:  userID,
		UpdatedAt:  time.Now(),
	}

	if err := t.DB.Create(tenantMetadata).Error; err != nil {
		return fmt.Errorf("failed to create tenant metadata: %w", err)
	}

	return nil
}

func (t *TenantMutationResolver) updateTenantInPermit(ctx context.Context, tenantID uuid.UUID, name string, attributes map[string]interface{}) error {
	_, err := t.PermitClient.SendRequest(ctx, "PATCH", fmt.Sprintf("tenants/%s", tenantID), map[string]interface{}{
		"name":       name,
		"attributes": attributes,
	})
	return err
}

func (t *TenantMutationResolver) updateTenantResource(tenantID uuid.UUID, name *string, parentID *uuid.UUID, userID uuid.UUID) error {
	updates := map[string]interface{}{
		"updated_by": userID,
		"updated_at": time.Now(),
	}
	if name != nil {
		updates["name"] = *name
	}

	if parentID != nil && *parentID != uuid.Nil {
		parentResourceID, err := t.validateParentOrg(*parentID)
		if err != nil {
			return fmt.Errorf("error getting parent org: %w", err)
		}
		updates["parent_resource_id"] = parentResourceID
	}

	if err := t.DB.Model(&dto.TenantResource{}).Where("resource_id = ?", tenantID).Updates(updates).Error; err != nil {
		return fmt.Errorf("error updating tenant resource: %w", err)
	}

	return nil
}

func (t *TenantMutationResolver) updateMetadata(resourceID uuid.UUID, description *string, contactInfo *models.ContactInfoInput, userID uuid.UUID) error {
	var tenantMetadata dto.TenantMetadata
	if err := t.DB.Where("resource_id = ?", resourceID).First(&tenantMetadata).Error; err != nil {
		return fmt.Errorf("tenant metadata not found: %w", err)
	}

	metadata := make(map[string]interface{})
	if err := json.Unmarshal(tenantMetadata.Metadata, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if description != nil {
		metadata["description"] = *description
	}

	if contactInfo != nil {
		t.updateContactInfo(metadata, contactInfo)
	}

	updatedMetadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal updated metadata: %w", err)
	}

	updates := map[string]interface{}{
		"metadata":   updatedMetadataJSON,
		"updated_at": time.Now(),
	}

	if err := t.DB.Model(&dto.TenantMetadata{}).Where("resource_id = ?", resourceID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update tenant metadata: %w", err)
	}

	return nil
}

func (t *TenantMutationResolver) updateContactInfo(metadata map[string]interface{}, input *models.ContactInfoInput) {
	contactInfo, ok := metadata["contactInfo"].(map[string]interface{})
	if !ok {
		contactInfo = make(map[string]interface{})
	}

	if input.Email != nil {
		contactInfo["email"] = *input.Email
	}
	if input.PhoneNumber != nil {
		contactInfo["phoneNumber"] = *input.PhoneNumber
	}
	if input.Address != nil {
		t.updateAddress(contactInfo, input.Address)
	}

	metadata["contactInfo"] = contactInfo
}

func (t *TenantMutationResolver) updateAddress(contactInfo map[string]interface{}, address *models.AddressInput) {
	addressMap, ok := contactInfo["address"].(map[string]interface{})
	if !ok {
		addressMap = make(map[string]interface{})
	}

	if address.Street != nil {
		addressMap["street"] = *address.Street
	}
	if address.City != nil {
		addressMap["city"] = *address.City
	}
	if address.State != nil {
		addressMap["state"] = *address.State
	}
	if address.ZipCode != nil {
		addressMap["zipCode"] = *address.ZipCode
	}
	if address.Country != nil {
		addressMap["country"] = *address.Country
	}

	contactInfo["address"] = addressMap
}

func (t *TenantMutationResolver) getTenantResponse(ctx context.Context, tenantID uuid.UUID) (models.OperationResult, error) {
	tq := &TenantQueryResolver{DB: t.DB, PC: t.PermitClient}
	return tq.Tenant(ctx, tenantID)
}

func (t *TenantMutationResolver) getParentOrg(tenantID, resourceTypeID uuid.UUID) (*dto.TenantResource, error) {
	var parentOrg dto.TenantResource
	if err := t.DB.Where(&dto.TenantResource{
		TenantID:       &tenantID,
		ResourceTypeID: resourceTypeID,
	}).First(&parentOrg).Error; err != nil {
		return nil, fmt.Errorf("error getting parent org: %w", err)
	}
	return &parentOrg, nil
}

func (t *TenantMutationResolver) getTenantResource(tx *gorm.DB, tenantID uuid.UUID) (*dto.TenantResource, error) {
	var tenant dto.TenantResource
	if err := tx.Where("resource_id = ? AND row_status = 1", tenantID).First(&tenant).Error; err != nil {
		return nil, fmt.Errorf("error getting tenant: %w", err)
	}
	return &tenant, nil
}

func (t *TenantMutationResolver) deleteTenantInPermit(ctx context.Context, tenantID uuid.UUID) error {
	_, err := t.PermitClient.SendRequest(ctx, "DELETE", fmt.Sprintf("tenants/%s", tenantID), nil)
	return err
}

func (t *TenantMutationResolver) updateTenantMetadata(tx *gorm.DB, tenantID uuid.UUID) error {
	if err := tx.Model(&dto.TenantMetadata{}).Where("resource_id = ?", tenantID).UpdateColumns(validations.UpdateDeletedMap()).Error; err != nil {
		return fmt.Errorf("error updating tenant metadata: %w", err)
	}
	return nil
}

func (t *TenantMutationResolver) updateTenantResourceStatus(tx *gorm.DB, resourceID uuid.UUID) error {
	if err := tx.Model(&dto.TenantResource{}).Where("resource_id= ?", resourceID).UpdateColumns(validations.UpdateDeletedMap()).Error; err != nil {
		return fmt.Errorf("error updating tenant resource: %w", err)
	}
	return nil
}
