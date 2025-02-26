package tenants

import (
	"context"
	"encoding/json"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/helpers"
	"iam_services_main_v1/internal/dao"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/utils"
	"iam_services_main_v1/internal/validations"
	"iam_services_main_v1/pkg/logger"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantMutationResolver struct {
	DB           *gorm.DB
	PermitClient *permit.PermitClient
}

// CreateTenant resolver for adding a new Tenant
func (t *TenantMutationResolver) CreateTenant(ctx context.Context, input models.CreateTenantInput) (models.OperationResult, error) {

	parentID, err := t.validateParentOrg(*input.ParentID)
	if err != nil {
		em := fmt.Sprint(err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("404", "Invalid parent organization", em)), nil
	}
	newtenantID := uuid.New()
	// Extract gin.Context from GraphQL context
	//ginCtx, ok := ctx.Value(middlewares.GinContextKey).(*gin.Context)
	// if !ok {
	// 	return nil, fmt.Errorf("unable to get gin context")
	// }
	//UserID := ginCtx.MustGet("userID").(string)
	//userUUID := uuid.MustParse(UserID)
	userUUID := uuid.New()
	inputMap := helpers.StructToMap(input)
	if err != nil {
		em := fmt.Sprintf("Error validating create tenant input: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("400", "Invalid input", em)), nil
	}
	var resourceTypeRoot dto.Mst_ResourceTypes
	if err := t.DB.Where("name = ?", "Root").First(&resourceTypeRoot).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	inputMap["created_by"] = userUUID
	inputMap["updated_by"] = userUUID
	// Create tenant in permit

	if _, err = t.PermitClient.SendRequest(ctx, "POST", "tenants", map[string]interface{}{
		"name":       input.Name,
		"key":        newtenantID, //generate uuid
		"attributes": inputMap,
	}); err != nil {
		em := fmt.Sprintf("Error creating tenant in permit system: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error creating tenant in permit system", em)), nil
	}

	resourceType, err := dao.GetResourceTypeByName("Tenant")
	if err != nil {
		em := fmt.Sprintf("Error getting resource type: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error getting resource type", em)), nil
	}

	// Create resource instance
	if _, err = t.PermitClient.SendRequest(ctx, "POST", "resource_instances", map[string]interface{}{
		"key":        input.ID,
		"resource":   resourceType.ResourceTypeID,
		"tenant":     newtenantID,
		"attributes": input,
	}); err != nil {
		em := fmt.Sprintf("Error creating resource instance of tenant in permit system: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error creating resource instance of tenant in permit system", em)), nil
	}

	//set realtionship between tenant and Root
	if _, err = t.PermitClient.SendRequest(ctx, "POST", "relationship_tuples", map[string]interface{}{
		"subject":  fmt.Sprint(resourceTypeRoot.ResourceTypeID, ":", newtenantID),
		"relation": "parent",
		"object":   fmt.Sprint(resourceType.ResourceTypeID, ":", newtenantID),
		"tenant":   newtenantID,
	}); err != nil {
		em := fmt.Sprintf("Error creating resource relationship of tenant in permit system: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error creating resource relationship of tenant in permit system", em)), nil
	}

	// Create tenant resource
	tenantResource, err := t.createTenantResource(input.Name, newtenantID, *parentID, userUUID, uuid.Nil)
	if err != nil {
		em := fmt.Sprintf("Error creating tenant resource: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error creating tenant resource", em)), nil
	}

	// Create tenant metadata
	if err := t.createTenantMetadata(tenantResource.ResourceID, input.Description, input.ContactInfo, userUUID); err != nil {
		em := fmt.Sprintf("Error creating tenant metadata: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error creating tenant metadata", em)), nil
	}
	// Create tenant resource
	tenantResource, err = t.createTenantResource(input.Name, input.ID, *parentID, userUUID, newtenantID)
	if err != nil {
		em := fmt.Sprintf("Error creating tenant resource: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error creating tenant resource", em)), nil
	}

	// Create tenant metadata
	if err := t.createTenantMetadata(tenantResource.ResourceID, input.Description, input.ContactInfo, userUUID); err != nil {
		em := fmt.Sprintf("Error creating tenant metadata: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error creating tenant metadata", em)), nil
	}

	tq := &TenantQueryResolver{DB: t.DB, PC: t.PermitClient}
	res, err := tq.ETenant(ctx, newtenantID)
	if err != nil {
		em := fmt.Sprintf("Error getting tenant: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error getting tenant", em)), nil
	}

	var data []models.Data
	data = append(data, res)
	return utils.FormatSuccess(data)
}

// UpdateTenant resolver for updating a Tenant
func (t *TenantMutationResolver) UpdateTenant(ctx context.Context, input models.UpdateTenantInput) (models.OperationResult, error) {

	resourceType, err := dao.GetResourceTypeByName("Tenant")
	if err != nil {
		em := fmt.Sprintf("Error getting resource type: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error getting resource type", em)), nil
	}

	var parentOrg *dto.TenantResource

	if err := t.DB.Where(&dto.TenantResource{
		TenantID:       &input.ID,
		ResourceTypeID: resourceType.ResourceTypeID,
	}).First(&parentOrg).Error; err != nil {
		em := fmt.Sprintf("Error getting parent org: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error getting parent org", em)), nil
	}

	inputMap := helpers.StructToMap(input)

	inputMap["created_by"] = parentOrg.CreatedBy
	inputMap["updated_by"] = parentOrg.UpdatedBy

	// Update tenant in permit
	if _, err := t.PermitClient.SendRequest(ctx, "PATCH", fmt.Sprintf("tenants/%s", input.ID), map[string]interface{}{
		"name":       input.Name,
		"attributes": inputMap,
	}); err != nil {
		em := fmt.Sprintf("Error updating tenant in permit system: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error updating tenant in permit system", em)), nil
	}

	// Update tenant resource
	updates := map[string]interface{}{
		"updated_by": parentOrg.CreatedBy,
		"updated_at": time.Now(),
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}

	if input.ParentID != nil && *input.ParentID != uuid.Nil {
		parentID, err := t.validateParentOrg(*input.ParentID)
		if err != nil {
			em := fmt.Sprintf("Error getting parent org: %v", err)
			logger.LogError(em)
			return utils.FormatError(utils.FormatErrorStruct("500", "Error getting parent org", em)), nil
		}
		updates["parent_resource_id"] = parentID
	}

	if err := t.DB.Model(&dto.TenantResource{}).Where("resource_id = ?", input.ID).Updates(updates).Error; err != nil {
		em := fmt.Sprintf("Error updating tenant resource: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error updating tenant resource", em)), nil
	}

	// Update metadata
	if err := t.updateMetadata(input.ID, input.Description, input.ContactInfo, parentOrg.CreatedBy); err != nil {
		em := fmt.Sprintf("Error updating tenant metadata: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error updating tenant metadata", em)), nil
	}

	tq := &TenantQueryResolver{DB: t.DB, PC: t.PermitClient}
	res, err := tq.ETenant(ctx, input.ID)
	if err != nil {
		em := fmt.Sprintf("Error getting tenant: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error getting tenant", em)), nil
	}

	var data []models.Data
	data = append(data, res)

	// Return success response with tenants
	return utils.FormatSuccess(data)
}

// DeleteTenant resolver for deleting a Tenant
func (t *TenantMutationResolver) DeleteTenant(ctx context.Context, input models.DeleteInput) (models.OperationResult, error) {
	tx := t.DB.Begin()

	// Delete from permit
	if _, err := t.PermitClient.SendRequest(ctx, "DELETE", fmt.Sprintf("tenants/%s", input.ID), nil); err != nil {
		tx.Rollback()
		em := fmt.Sprintf("Error deleting tenant in permit system: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error deleting tenant in permit system", em)), nil
	}

	// Update metadata
	if err := tx.Model(&dto.TenantMetadata{}).Where("resource_id = ?", input.ID).UpdateColumns(validations.UpdateDeletedMap()).Error; err != nil {
		tx.Rollback()
		em := fmt.Sprintf("Error updating tenant metadata: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error updating tenant metadata", em)), nil
	}

	// Update resource
	if err := tx.Model(&dto.TenantResource{}).Where("resource_id= ?", input.ID).UpdateColumns(validations.UpdateDeletedMap()).Error; err != nil {
		tx.Rollback()
		em := fmt.Sprintf("Error updating tenant resource: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error updating tenant resource", em)), nil
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		em := fmt.Sprintf("Error committing transaction: %v", err)
		logger.LogError(em)
		return utils.FormatError(utils.FormatErrorStruct("500", "Error committing transaction", em)), nil
	}

	return utils.FormatSuccess([]models.Data{})
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

func (t *TenantMutationResolver) createTenantResource(name string, resourceID, parentID uuid.UUID, UserID, tenantID uuid.UUID) (*dto.TenantResource, error) {
	var resourceType dto.Mst_ResourceTypes
	if err := t.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	tenant := &dto.TenantResource{
		ResourceID:       resourceID,
		Name:             name,
		CreatedBy:        UserID,
		UpdatedBy:        UserID,
		CreatedAt:        time.Now(),
		ResourceTypeID:   resourceType.ResourceTypeID,
		ParentResourceID: &parentID,
		//TenantID:         &tenantID,
	}
	if tenantID != uuid.Nil {
		tenant.TenantID = &tenantID
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
		UpdatedBy:  UserID,
		UpdatedAt:  time.Now(),
	}

	if err := t.DB.Create(tenantMetadata).Error; err != nil {
		return fmt.Errorf("failed to create tenant metadata: %w", err)
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
		"metadata": updatedMetadataJSON,
		//"updated_by": userID,
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
