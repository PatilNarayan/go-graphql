package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_graphql/config"
	"go_graphql/gql/models"
	"go_graphql/internal/constants"
	"go_graphql/internal/dto"
	"go_graphql/internal/role"
	"go_graphql/permit"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantMutationResolver struct {
	DB *gorm.DB
}

// CreateTenant resolver for adding a new Tenant
func (r *TenantMutationResolver) CreateTenant(ctx context.Context, input models.CreateTenantInput) (*models.Tenant, error) {

	tenantResource := &dto.TenantResource{
		ResourceID: uuid.New(),
		Name:       input.Name,
		CreatedBy:  constants.DefaltCreatedBy, //input.CreatedBy,
		UpdatedBy:  constants.DefaltCreatedBy, //input.CreatedBy,
		CreatedAt:  time.Now(),
	}

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}
	tenantResource.ResourceTypeID = resourceType.ResourceTypeID

	if input.ParentOrgID == "" {
		return nil, errors.New("parent organization ID is required")
	}
	// Validate ParentOrgID
	resourceTypeId, err := GetResourceTypeIDs([]string{"Root"})
	if err != nil {
		return nil, fmt.Errorf("failed to get resource type IDs: %w", err)
	}
	var parentOrg dto.TenantResource
	if err := r.DB.Where(
		"resource_id = ? AND resource_type_id in (?) AND row_status = 1",
		input.ParentOrgID, resourceTypeId,
	).First(&parentOrg).Error; err != nil {
		return nil, fmt.Errorf("parent organization not found: %w", err)
	}
	pid := uuid.MustParse(input.ParentOrgID)
	tenantResource.ParentResourceID = &pid

	pc := permit.NewPermitClient()
	_, err = pc.SendRequest(ctx, "POST", "tenants", map[string]interface{}{
		"name": input.Name,
		"key":  tenantResource.ResourceID.String(),
	})

	if err != nil {
		return nil, err
	}

	if err := r.DB.Create(&tenantResource).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant resource: %w", err)
	}

	metadata := map[string]interface{}{
		"description": input.Description,
		"contactInfo": input.ContactInfo,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	tenantMetadata := &dto.TenantMetadata{
		ResourceID: tenantResource.ResourceID,
		Metadata:   metadataJSON,
		CreatedBy:  constants.DefaltCreatedBy, //input.CreatedBy,
		CreatedAt:  time.Now(),
	}

	if err := r.DB.Create(&tenantMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant metadata: %w", err)
	}

	err = role.CreateMstRole(tenantResource.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	tq := &TenantQueryResolver{DB: r.DB}
	return tq.GetTenant(ctx, tenantResource.ResourceID)

}

// UpdateTenant resolver for updating a Tenant
func (r *TenantMutationResolver) UpdateTenant(ctx context.Context, input models.UpdateTenantInput) (*models.Tenant, error) {

	// Fetch the existing TenantResource
	var tenantResource dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{ResourceID: input.ID}).First(&tenantResource).Error; err != nil {
		return nil, fmt.Errorf("tenant resource not found: %w", err)
	}

	// Update TenantResource fields if provided
	if input.Name != nil && *input.Name != "" {
		tenantResource.Name = *input.Name
	}
	if input.ParentOrgID != nil {
		resourceTypeId, err := GetResourceTypeIDs([]string{"Root"})
		if err != nil {
			return nil, fmt.Errorf("failed to get resource type IDs: %w", err)
		}
		var parentOrg dto.TenantResource
		if err := r.DB.Where(
			"resource_id = ? AND resource_type_id in (?) AND row_status = 1",
			input.ParentOrgID, resourceTypeId,
		).First(&parentOrg).Error; err != nil {
			return nil, fmt.Errorf("parent organization not found: %w", err)
		}
		parsedUUID := parentOrg.ResourceID
		tenantResource.ParentResourceID = &parsedUUID
	} else {
		return nil, errors.New("parent organization ID is required")
	}
	tenantResource.UpdatedBy = constants.DefaltUpdatedBy //input.UpdatedBy
	tenantResource.UpdatedAt = time.Now()

	pc := permit.NewPermitClient()
	_, err := pc.SendRequest(ctx, "PATCH", "tenants/"+input.ID.String(), map[string]interface{}{
		"name": input.Name,
	})
	if err != nil {
		return nil, err
	}

	// Save updated TenantResource to the database
	if err := r.DB.Save(&tenantResource).Error; err != nil {
		return nil, fmt.Errorf("failed to update tenant resource: %w", err)
	}

	// Fetch the existing TenantMetadata
	var tenantMetadata dto.TenantMetadata
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: tenantResource.ResourceID}).First(&tenantMetadata).Error; err != nil {
		return nil, fmt.Errorf("tenant metadata not found: %w", err)
	}

	// Unmarshal the existing metadata
	metadata := map[string]interface{}{}
	if err := json.Unmarshal(tenantMetadata.Metadata, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Update metadata fields if provided
	if input.Description != nil && *input.Description != "" {
		metadata["description"] = *input.Description
	}
	if input.ContactInfo != nil {
		contactInfo := metadata["contactInfo"].(map[string]interface{})
		if input.ContactInfo.Email != nil && *input.ContactInfo.Email != "" {
			contactInfo["email"] = *input.ContactInfo.Email
		}
		if input.ContactInfo.PhoneNumber != nil && *input.ContactInfo.PhoneNumber != "" {
			contactInfo["phoneNumber"] = *input.ContactInfo.PhoneNumber
		}
		if input.ContactInfo.Address != nil {
			address := contactInfo["address"].(map[string]interface{})
			if input.ContactInfo.Address.Street != nil && *input.ContactInfo.Address.Street != "" {
				address["street"] = *input.ContactInfo.Address.Street
			}
			if input.ContactInfo.Address.City != nil && *input.ContactInfo.Address.City != "" {
				address["city"] = *input.ContactInfo.Address.City
			}
			if input.ContactInfo.Address.State != nil && *input.ContactInfo.Address.State != "" {
				address["state"] = *input.ContactInfo.Address.State
			}
			if input.ContactInfo.Address.ZipCode != nil && *input.ContactInfo.Address.ZipCode != "" {
				address["zipCode"] = *input.ContactInfo.Address.ZipCode
			}
			if input.ContactInfo.Address.Country != nil && *input.ContactInfo.Address.Country != "" {
				address["country"] = *input.ContactInfo.Address.Country
			}
		}
		metadata["contactInfo"] = contactInfo
	}

	// Marshal the updated metadata back to JSON
	updatedMetadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated metadata: %w", err)
	}
	tenantMetadata.Metadata = updatedMetadataJSON
	tenantMetadata.UpdatedBy = "system" // Replace with actual user
	tenantMetadata.UpdatedAt = time.Now()

	// Save updated TenantMetadata to the database
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: tenantResource.ResourceID}).UpdateColumns(&tenantMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to update tenant metadata: %w", err)
	}

	tq := &TenantQueryResolver{DB: r.DB}
	return tq.GetTenant(ctx, tenantResource.ResourceID)
}

func (r *TenantMutationResolver) DeleteTenant(ctx context.Context, id uuid.UUID) (bool, error) {
	tx := r.DB.Begin()

	var tenantResource dto.TenantResource
	if err := tx.Where(&dto.TenantResource{ResourceID: id}).First(&tenantResource).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("tenant resource not found: %w", err)
		}
		return false, fmt.Errorf("failed to fetch tenant resource: %w", err)
	}

	// Update TenantMetadata with both DeletedAt and RowStatus
	updates := map[string]interface{}{
		// "deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}
	if err := tx.Model(&dto.TenantMetadata{}).Where("resource_id = ?", id.String()).Updates(updates).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to soft delete tenant metadata: %w", err)
	}

	pc := permit.NewPermitClient()
	_, err := pc.SendRequest(ctx, "DELETE", "tenants/"+id.String(), nil)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to delete tenant from PDP: %w", err)
	}

	// Update TenantResource with both DeletedAt and RowStatus
	if err := tx.Model(&dto.TenantResource{}).Where("resource_id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to delete tenant resource: %w", err)
	}

	if err := role.DeleteDefaultRole(id); err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to delete default role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return true, nil
}

func GetResourceTypeIDs(resourceName []string) ([]string, error) {
	resourceType := []dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name in (?) AND row_status = 1", resourceName).Find(&resourceType).Error; err != nil {
		return nil, err
	}
	var resourceIds []string
	for _, resource := range resourceType {
		resourceIds = append(resourceIds, resource.ResourceTypeID.String())
	}

	return resourceIds, nil
}
