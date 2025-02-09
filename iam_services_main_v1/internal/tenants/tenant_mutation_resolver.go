package tenants

import (
	"context"
	"encoding/json"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/constants"
	"iam_services_main_v1/internal/dao"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/roles"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantMutationResolver struct {
	DB           *gorm.DB
	PermitClient *permit.PermitClient
}

func (t *TenantMutationResolver) validateParentOrg(parentOrgID string) (*uuid.UUID, error) {
	if parentOrgID == "" {
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

func (t *TenantMutationResolver) createTenantResource(name string, resourceID, parentID uuid.UUID) (*dto.TenantResources, error) {
	var resourceType dto.Mst_ResourceTypes
	if err := t.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	tenant := &dto.TenantResources{
		ResourceID:       resourceID,
		Name:             name,
		CreatedBy:        constants.DefaltCreatedBy,
		UpdatedBy:        constants.DefaltCreatedBy,
		CreatedAt:        time.Now(),
		ResourceTypeID:   resourceType.ResourceTypeID,
		ParentResourceID: &parentID,
	}

	if err := t.DB.Create(tenant).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant resource: %w", err)
	}

	return tenant, nil
}

func (t *TenantMutationResolver) createTenantMetadata(resourceID uuid.UUID, description *string, contactInfo *models.ContactInfoInput) error {
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
		CreatedBy:  constants.DefaltCreatedBy,
		CreatedAt:  time.Now(),
	}

	if err := t.DB.Create(tenantMetadata).Error; err != nil {
		return fmt.Errorf("failed to create tenant metadata: %w", err)
	}

	return nil
}

func (t *TenantMutationResolver) updateMetadata(resourceID uuid.UUID, description *string, contactInfo *models.ContactInfoInput) error {
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
		"updated_by": constants.DefaltUpdatedBy,
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

// CreateTenant resolver for adding a new Tenant
func (t *TenantMutationResolver) CreateTenant(ctx context.Context, input models.CreateTenantInput) (*models.Tenant, error) {

	parentID, err := t.validateParentOrg(input.ParentOrgID)
	if err != nil {
		return nil, err
	}

	// Create tenant in permit
	var tenantData interface{}
	if tenantData, err = t.PermitClient.SendRequest(ctx, "POST", "tenants", map[string]interface{}{
		"name":       input.Name,
		"key":        uuid.NewString(),
		"attributes": input,
	}); err != nil {
		return nil, err
	}

	tenantID, ok := tenantData.(map[string]interface{})["id"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to create tenant in permit")
	}
	tenantUuid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tenant ID: %w", err)
	}

	// Create resource instance
	if _, err = t.PermitClient.SendRequest(ctx, "POST", "resource_instances", map[string]interface{}{
		"key":        uuid.NewString(),
		"resource":   "tenant",
		"tenant":     tenantID,
		"attributes": input,
	}); err != nil {
		return nil, err
	}

	// Create tenant resource
	tenantResource, err := t.createTenantResource(input.Name, tenantUuid, *parentID)
	if err != nil {
		return nil, err
	}

	// Create tenant metadata
	if err := t.createTenantMetadata(tenantResource.ResourceID, input.Description, input.ContactInfo); err != nil {
		return nil, err
	}

	// Create default role
	if err := roles.CreateMstRole(tenantResource.ResourceID); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	tq := &TenantQueryResolver{DB: t.DB, PermitClient: t.PermitClient}
	return tq.GetTenant(ctx, tenantResource.ResourceID)
}

// UpdateTenant resolver for updating a Tenant
func (t *TenantMutationResolver) UpdateTenant(ctx context.Context, input models.UpdateTenantInput) (*models.Tenant, error) {

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
		"updated_by": constants.DefaltUpdatedBy,
		"updated_at": time.Now(),
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.ParentOrgID != nil {
		parentID, err := t.validateParentOrg(*input.ParentOrgID)
		if err != nil {
			return nil, err
		}
		updates["parent_resource_id"] = parentID
	}

	if err := t.DB.Model(&dto.TenantResources{}).Where("resource_id = ?", input.ID).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update tenant resource: %w", err)
	}

	// Update metadata
	if err := t.updateMetadata(input.ID, input.Description, input.ContactInfo); err != nil {
		return nil, err
	}

	tq := &TenantQueryResolver{DB: t.DB, PermitClient: t.PermitClient}
	return tq.GetTenant(ctx, input.ID)
}

// DeleteTenant resolver for deleting a Tenant
func (t *TenantMutationResolver) DeleteTenant(ctx context.Context, id uuid.UUID) (bool, error) {
	tx := t.DB.Begin()

	updates := map[string]interface{}{
		"row_status": 0,
	}

	// Delete from permit
	if _, err := t.PermitClient.SendRequest(ctx, "DELETE", fmt.Sprintf("tenants/%s", id), nil); err != nil {
		tx.Rollback()
		return false, err
	}

	// Update metadata
	if err := tx.Model(&dto.TenantMetadata{}).Where("resource_id = ?", id).UpdateColumns(updates).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to soft delete tenant metadata: %w", err)
	}

	// Update resource
	if err := tx.Model(&dto.TenantResources{}).Where("resource_id = ?", id).UpdateColumns(updates).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to delete tenant resource: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return true, nil
}
