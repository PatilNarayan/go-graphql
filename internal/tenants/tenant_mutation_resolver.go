package tenants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_graphql/gql/models"
	"go_graphql/internal/constants"
	"go_graphql/internal/dto"
	"go_graphql/internal/role"
	"go_graphql/logger"
	"go_graphql/permit"
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
	log := logger.Log.WithField("operation", "CreateTenant")
	log.Info("Starting tenant creation process")

	tenantResource := &dto.TenantResource{
		ResourceID: uuid.New(),
		Name:       input.Name,
		CreatedBy:  constants.DefaltCreatedBy, //input.CreatedBy,
		UpdatedBy:  constants.DefaltCreatedBy, //input.CreatedBy,
		CreatedAt:  time.Now(),
	}
	log.WithField("tenantID", tenantResource.ResourceID).Info("Generated new tenant resource")

	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
		log.WithError(err).Error("Failed to find resource type")
		return nil, fmt.Errorf("resource type not found: %w", err)
	}
	tenantResource.ResourceTypeID = resourceType.ResourceTypeID

	if input.ParentOrgID != nil {
		var parentOrg dto.TenantResource
		if err := r.DB.Where(&dto.TenantResource{ResourceID: *input.ParentOrgID}).First(&parentOrg).Error; err != nil {
			log.WithError(err).Error("Failed to find parent organization")
			return nil, fmt.Errorf("parent organization not found: %w", err)
		}
		tenantResource.ParentResourceID = input.ParentOrgID
		log.WithField("parentOrgID", input.ParentOrgID).Info("Parent organization validated")
	}

	pc := permit.NewPermitClient()
	_, err := pc.APIExecute(ctx, "POST", "tenants", map[string]interface{}{
		"name": input.Name,
		"key":  tenantResource.ResourceID.String(),
	})

	if err != nil {
		return nil, err
	}

	if err := r.DB.Create(&tenantResource).Error; err != nil {
		log.WithError(err).Error("Failed to create tenant resource")
		return nil, fmt.Errorf("failed to create tenant resource: %w", err)
	}
	log.Info("Tenant resource created successfully")

	metadata := map[string]interface{}{
		"description": input.Description,
		"contactInfo": input.ContactInfo,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		log.WithError(err).Error("Failed to marshal metadata")
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	tenantMetadata := &dto.TenantMetadata{
		ResourceID: tenantResource.ResourceID.String(),
		Metadata:   metadataJSON,
		CreatedBy:  constants.DefaltCreatedBy, //input.CreatedBy,
		CreatedAt:  time.Now(),
	}

	if err := r.DB.Create(&tenantMetadata).Error; err != nil {
		log.WithError(err).Error("Failed to create tenant metadata")
		return nil, fmt.Errorf("failed to create tenant metadata: %w", err)
	}
	log.Info("Tenant metadata created successfully")

	err = role.CreateMstRole(tenantResource.ResourceID)
	if err != nil {
		log.WithError(err).Error("Failed to create role")
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	log.Info("Role created successfully")

	result := &models.Tenant{
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
	}
	log.WithField("tenantID", result.ID).Info("Tenant creation completed successfully")
	return result, nil
}

// UpdateTenant resolver for updating a Tenant
func (r *TenantMutationResolver) UpdateTenant(ctx context.Context, input models.UpdateTenantInput) (*models.Tenant, error) {
	log := logger.Log.WithField("operation", "UpdateTenant")
	log.Info("Starting tenant update process")

	// Fetch the existing TenantResource
	var tenantResource dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{ResourceID: input.ID}).First(&tenantResource).Error; err != nil {
		log.WithError(err).Error("tenant resource not found")
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
			log.WithError(err).Error("parent organization not found")
			return nil, fmt.Errorf("parent organization not found: %w", err)
		}
		parsedUUID := uuid.MustParse(*input.ParentOrgID)
		tenantResource.ParentResourceID = &parsedUUID
	}
	tenantResource.UpdatedBy = constants.DefaltUpdatedBy //input.UpdatedBy
	tenantResource.UpdatedAt = time.Now()

	pc := permit.NewPermitClient()
	_, err := pc.APIExecute(ctx, "PATCH", "tenants/"+input.ID.String(), map[string]interface{}{
		"name": input.Name,
	})
	if err != nil {
		return nil, err
	}

	// Save updated TenantResource to the database
	if err := r.DB.Save(&tenantResource).Error; err != nil {
		log.WithError(err).Error("failed to update tenant resource")
		return nil, fmt.Errorf("failed to update tenant resource: %w", err)
	}

	// Fetch the existing TenantMetadata
	var tenantMetadata dto.TenantMetadata
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: tenantResource.ResourceID.String()}).First(&tenantMetadata).Error; err != nil {
		log.WithError(err).Error("tenant metadata not found")
		return nil, fmt.Errorf("tenant metadata not found: %w", err)
	}

	// Unmarshal the existing metadata
	metadata := map[string]interface{}{}
	if err := json.Unmarshal(tenantMetadata.Metadata, &metadata); err != nil {
		log.WithError(err).Error("failed to unmarshal metadata")
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
		log.WithError(err).Error("failed to marshal updated metadata")
		return nil, fmt.Errorf("failed to marshal updated metadata: %w", err)
	}
	tenantMetadata.Metadata = updatedMetadataJSON
	tenantMetadata.UpdatedBy = "system" // Replace with actual user
	tenantMetadata.UpdatedAt = time.Now()

	// Save updated TenantMetadata to the database
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: tenantResource.ResourceID.String()}).UpdateColumns(&tenantMetadata).Error; err != nil {
		log.WithError(err).Error("failed to update tenant metadata")
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
	tx := r.DB.Begin()
	log := logger.Log.WithField("operation", "DeleteTenant")

	var tenantResource dto.TenantResource
	if err := tx.Where(&dto.TenantResource{ResourceID: id}).First(&tenantResource).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to fetch tenant resource")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("tenant resource not found: %w", err)
		}
		return false, fmt.Errorf("failed to fetch tenant resource: %w", err)
	}

	log.WithField("tenantID", id).Info("Found tenant resource")

	// Update TenantMetadata with both DeletedAt and RowStatus
	updates := map[string]interface{}{
		// "deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}
	if err := tx.Model(&dto.TenantMetadata{}).Where("resource_id = ?", id.String()).Updates(updates).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to soft delete tenant metadata")
		return false, fmt.Errorf("failed to soft delete tenant metadata: %w", err)
	}

	log.Info("Deleted tenant metadata")

	pc := permit.NewPermitClient()
	_, err := pc.APIExecute(ctx, "DELETE", "tenants/"+id.String(), nil)
	if err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to delete tenant from PDP")
		return false, fmt.Errorf("failed to delete tenant from PDP: %w", err)
	}

	// Update TenantResource with both DeletedAt and RowStatus
	if err := tx.Model(&dto.TenantResource{}).Where("resource_id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to delete tenant resource")
		return false, fmt.Errorf("failed to delete tenant resource: %w", err)
	}

	log.Info("Deleted tenant resource")

	if err := tx.Commit().Error; err != nil {
		log.WithError(err).Error("failed to commit transaction")
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info("Deleted tenant successfully")

	return true, nil
}
