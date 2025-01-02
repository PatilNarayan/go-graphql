package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"time"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type TenantMutationResolver struct {
	DB *gorm.DB
}

// CreateTenant resolver for adding a new Tenant
func (r *TenantMutationResolver) CreateTenant(ctx context.Context, input models.CreateTenantInput) (*models.Tenant, error) {
	// Create a new TenantResource
	tenantResource := &dto.TenantResource{
		ResourceID: uuid.New(), // Generate new UUID
		Name:       input.Name,
		CreatedBy:  input.CreatedBy,
		UpdatedBy:  input.CreatedBy,
		CreatedAt:  time.Now(),
	}

	//get resource type by name
	resourceType := dto.Mst_ResourceType{}
	if err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}
	tenantResource.ResourceTypeID = resourceType.ResourceTypeID

	if input.ParentOrgID != nil {
		//check if parentOrgID is valid
		var parentOrg dto.TenantResource
		if err := r.DB.Where(&dto.TenantResource{ResourceID: *input.ParentOrgID}).First(&parentOrg).Error; err != nil {
			return nil, fmt.Errorf("parent organization not found: %w", err)
		}
		tenantResource.ParentResourceID = input.ParentOrgID
	}

	// pc := permit.NewPermitClient()
	// _, err := pc.APIExecute(ctx, "POST", "tenants", map[string]interface{}{
	// 	"name": input.Name,
	// 	"key":  tenantResource.ResourceID.String(),
	// })

	// if err != nil {
	// 	return nil, err
	// }

	// Save TenantResource to the database
	if err := r.DB.Create(&tenantResource).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant resource: %w", err)
	}

	// Prepare metadata (ContactInfo)
	metadata := map[string]interface{}{
		"description": input.Description,
		"contactInfo": input.ContactInfo,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create a new TenantMetadata
	tenantMetadata := &dto.TenantMetadata{
		ResourceID: tenantResource.ResourceID.String(),
		Metadata:   metadataJSON,
		CreatedBy:  input.CreatedBy,
		CreatedAt:  time.Now(),
	}

	// Save TenantMetadata to the database
	if err := r.DB.Create(&tenantMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant metadata: %w", err)
	}

	// Return the created Tenant object
	return &models.Tenant{
		ID:          tenantResource.ResourceID,
		Name:        tenantResource.Name,
		Description: input.Description,
		ContactInfo: &models.ContactInfo{
			Email:       input.ContactInfo.Email,
			PhoneNumber: input.ContactInfo.PhoneNumber,
			Address: &models.Address{
				Street:  input.ContactInfo.Address.Street,
				City:    input.ContactInfo.Address.City,
				State:   input.ContactInfo.Address.State,
				ZipCode: input.ContactInfo.Address.ZipCode,
				Country: input.ContactInfo.Address.Country,
			},
		},
		CreatedAt: tenantResource.CreatedAt.String(),
		CreatedBy: &tenantResource.CreatedBy,
	}, nil
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
	if input.ParentOrgID != nil && *input.ParentOrgID != "" {
		// Validate ParentOrgID
		var parentOrg dto.TenantResource
		if err := r.DB.Where(&dto.TenantResource{ResourceID: uuid.MustParse(*input.ParentOrgID)}).First(&parentOrg).Error; err != nil {
			return nil, fmt.Errorf("parent organization not found: %w", err)
		}
		parsedUUID := uuid.MustParse(*input.ParentOrgID)
		tenantResource.ParentResourceID = &parsedUUID
	}
	tenantResource.UpdatedBy = input.UpdatedBy
	tenantResource.UpdatedAt = time.Now()

	// Save updated TenantResource to the database
	if err := r.DB.Save(&tenantResource).Error; err != nil {
		return nil, fmt.Errorf("failed to update tenant resource: %w", err)
	}

	// Fetch the existing TenantMetadata
	var tenantMetadata dto.TenantMetadata
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: tenantResource.ResourceID.String()}).First(&tenantMetadata).Error; err != nil {
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
	if err := r.DB.Save(&tenantMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to update tenant metadata: %w", err)
	}

	// Return the updated Tenant object
	return &models.Tenant{
		ID:          tenantResource.ResourceID,
		Name:        tenantResource.Name,
		Description: ptr.String(metadata["description"].(string)),
		ContactInfo: &models.ContactInfo{
			Email:       ptr.String(metadata["contactInfo"].(map[string]interface{})["email"].(string)),
			PhoneNumber: ptr.String(metadata["contactInfo"].(map[string]interface{})["phoneNumber"].(string)),
			Address: &models.Address{
				Street:  ptr.String(metadata["contactInfo"].(map[string]interface{})["address"].(map[string]interface{})["street"].(string)),
				City:    ptr.String(metadata["contactInfo"].(map[string]interface{})["address"].(map[string]interface{})["city"].(string)),
				State:   ptr.String(metadata["contactInfo"].(map[string]interface{})["address"].(map[string]interface{})["state"].(string)),
				ZipCode: ptr.String(metadata["contactInfo"].(map[string]interface{})["address"].(map[string]interface{})["zipCode"].(string)),
				Country: ptr.String(metadata["contactInfo"].(map[string]interface{})["address"].(map[string]interface{})["country"].(string)),
			},
		},
		UpdatedAt: ptr.String(tenantResource.UpdatedAt.String()),
		UpdatedBy: &tenantResource.UpdatedBy,
	}, nil
}

func (r *TenantMutationResolver) DeleteTenant(ctx context.Context, id uuid.UUID) (bool, error) {
	// Start a database transaction
	tx := r.DB.Begin()

	// Fetch the TenantResource to ensure it exists
	var tenantResource dto.TenantResource
	if err := tx.Where(&dto.TenantResource{ResourceID: id}).First(&tenantResource).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("tenant resource not found: %w", err)
		}
		return false, fmt.Errorf("failed to fetch tenant resource: %w", err)
	}

	// Delete associated TenantMetadata
	if err := tx.Where(&dto.TenantMetadata{ResourceID: id.String()}).Delete(&dto.TenantMetadata{}).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to delete tenant metadata: %w", err)
	}

	// Delete the TenantResource
	if err := tx.Delete(&tenantResource).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to delete tenant resource: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return success
	return true, nil
}
