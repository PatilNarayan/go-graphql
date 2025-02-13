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

	parentID, err := t.validateParentOrg(input.ParentID)
	if err != nil {
		return nil, err
	}

	// Extract gin.Context from GraphQL context
	ginCtx, ok := ctx.Value(middlewares.GinContextKey).(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unable to get gin context")
	}

	UserID, err := helpers.GetUserID(ginCtx)
	if err != nil {
		return nil, err
	}

	userUUID := uuid.MustParse(UserID)

	inputMap, err := helpers.StructToMap(input)
	if err != nil {
		return nil, err
	}
	inputMap["created_by"] = UserID
	inputMap["updated_by"] = UserID
	// Create tenant in permit

	if _, err = t.PermitClient.SendRequest(ctx, "POST", "tenants", map[string]interface{}{
		"name":       input.Name,
		"key":        input.ID,
		"attributes": inputMap,
	}); err != nil {
		return nil, err
	}

	// Create resource instance
	if _, err = t.PermitClient.SendRequest(ctx, "POST", "resource_instances", map[string]interface{}{
		"key":        input.ID,
		"resource":   "tenant",
		"tenant":     input.ID,
		"attributes": input,
	}); err != nil {
		return nil, err
	}

	// Create tenant resource
	tenantResource, err := t.createTenantResource(input.Name, input.ID, *parentID, userUUID)
	if err != nil {
		return nil, err
	}

	// Create tenant metadata
	if err := t.createTenantMetadata(tenantResource.ResourceID, input.Description, input.ContactInfo, userUUID); err != nil {
		return nil, err
	}

	// // Create default role
	// if err := roles.CreateMstRole(tenantResource.ResourceID); err != nil {
	// 	return nil, fmt.Errorf("failed to create role: %w", err)
	// }

	tq := &TenantQueryResolver{DB: t.DB, PermitClient: t.PermitClient}
	return tq.GetTenant(ctx, tenantResource.ResourceID)
}

// UpdateTenant resolver for updating a Tenant
func (t *TenantMutationResolver) UpdateTenant(ctx context.Context, input models.UpdateTenantInput) (models.OperationResult, error) {

	// Extract gin.Context from GraphQL context
	ginCtx, ok := ctx.Value(middlewares.GinContextKey).(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unable to get gin context")
	}

	UserID, err := helpers.GetUserID(ginCtx)
	if err != nil {
		return nil, err
	}

	// Update tenant in permit
	if _, err := t.PermitClient.SendRequest(ctx, "PATCH", fmt.Sprintf("tenants/%s", input.ID), map[string]interface{}{
		"name":       input.Name,
		"attributes": input,
	}); err != nil {
		return nil, err
	}

	// Update resource instance
	if _, err := t.PermitClient.SendRequest(ctx, "PATCH", fmt.Sprintf("resource_instances/%s", input.ID), map[string]interface{}{
		"attributes": input,
	}); err != nil {
		return nil, err
	}

	// Update tenant resource
	updates := map[string]interface{}{
		"updated_by": UserID,
		"updated_at": time.Now(),
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.ParentID != uuid.Nil {
		parentID, err := t.validateParentOrg(input.ParentID)
		if err != nil {
			return nil, err
		}
		updates["parent_resource_id"] = parentID
	}

	if err := t.DB.Model(&dto.TenantResources{}).Where("resource_id = ?", input.ID).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update tenant resource: %w", err)
	}

	// Update metadata
	if err := t.updateMetadata(input.ID, input.Description, input.ContactInfo, UserID); err != nil {
		return nil, err
	}

	tq := &TenantQueryResolver{DB: t.DB, PermitClient: t.PermitClient}
	return tq.GetTenant(ctx, input.ID)
}

// DeleteTenant resolver for deleting a Tenant
func (t *TenantMutationResolver) DeleteTenant(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	tx := t.DB.Begin()

	updates := map[string]interface{}{
		"row_status": 0,
	}

	// Delete from permit
	if _, err := t.PermitClient.SendRequest(ctx, "DELETE", fmt.Sprintf("tenants/%s", id), nil); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update metadata
	if err := tx.Model(&dto.TenantMetadata{}).Where("resource_id = ?", id).UpdateColumns(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to soft delete tenant metadata: %w", err)
	}

	// Update resource
	if err := tx.Model(&dto.TenantResources{}).Where("resource_id = ?", id).UpdateColumns(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete tenant resource: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.SuccessResponse{
		Success: true,
		Message: "Tenant deleted successfully",
		Data:    []models.Data{},
	}, nil
}

func (t *TenantMutationResolver) validateParentOrg(parentOrgID uuid.UUID) (*uuid.UUID, error) {
	if parentOrgID == uuid.Nil {
		return nil, fmt.Errorf("parent organization ID is required")
	}

	resourceType, err := dao.GetResourceTypeByName("Root")
	if err != nil {
		return nil, fmt.Errorf("failed to get resource type IDs: %w", err)
	}

	var parentOrg dto.TenantResources
	if err := t.DB.Where(
		"resource_id = ? AND resource_type_id in (?) AND row_status = 1",
		parentOrgID, resourceType.ResourceTypeID,
	).First(&parentOrg).Error; err != nil {
		return nil, fmt.Errorf("parent organization not found: %w", err)
	}

	return &parentOrg.ResourceID, nil
}

func (t *TenantMutationResolver) createTenantResource(name string, resourceID, parentID uuid.UUID, UserID uuid.UUID) (*dto.TenantResources, error) {
	var resourceType dto.Mst_ResourceTypes
	if err := t.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	tenant := &dto.TenantResources{
		ResourceID:       resourceID,
		Name:             name,
		CreatedBy:        UserID,
		UpdatedBy:        UserID,
		CreatedAt:        time.Now(),
		ResourceTypeID:   resourceType.ResourceTypeID,
		ParentResourceID: &parentID,
	}

	if err := t.DB.Create(tenant).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant resource: %w", err)
	}

	return tenant, nil
}

func (t *TenantMutationResolver) createTenantMetadata(resourceID uuid.UUID, description *string, contactInfo *models.ContactInfoInput, UserID uuid.UUID) error {
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
		CreatedBy:  UserID,
		CreatedAt:  time.Now(),
	}

	if err := t.DB.Create(tenantMetadata).Error; err != nil {
		return fmt.Errorf("failed to create tenant metadata: %w", err)
	}

	return nil
}

func (t *TenantMutationResolver) updateMetadata(resourceID uuid.UUID, description *string, contactInfo *models.ContactInfoInput, userID string) error {
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
		"updated_by": userID,
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

func (t *TenantMutationResolver) updateAddress(contactInfo map[string]interface{}, address *models.CreateAddressInput) {
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
