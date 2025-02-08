package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_graphql/internal/dto"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type AccountMutationResolver struct {
	DB *gorm.DB
}

// CreateAccount resolver for adding a new Account
func (r *AccountMutationResolver) CreateAccount(ctx context.Context, input models.CreateAccountInput) (*models.Account, error) {
	// var inputRequest dto.CreateAccountInput
	// ginContext := ctx.Value("GinContextKey").(*gin.Context)
	// if err := ginContext.ShouldBindJSON(&inputRequest); err != nil {
	// 	ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	ginContext.Abort()
	// 	return nil, nil
	// }

	userID := uuid.New()
	TenantResourceTypeID := "ed113bda-bbda-11ef-87ea-c03c5946f955"
	parsedTenantResourceTypeID, err := uuid.Parse(TenantResourceTypeID)
	if err != nil {
		return nil, fmt.Errorf("error parsing resource type  UUID: %w", err)
	}
	accountResource := &dto.TenantResource{
		ResourceID:       uuid.New(), // Generate new UUID
		Name:             input.Name,
		ParentResourceID: &input.ParentID,
		TenantID:         &input.TenantID,
		RowStatus:        1,
		CreatedBy:        userID.String(),
	}
	resourceType, err := dao.GetResourceTypeByName("Account")
	if err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}
	accountResource.ResourceTypeID = resourceType.ResourceTypeID
	if input.ParentID != uuid.Nil {
		//check if ParentID is valid
		resourceData, err := dao.GetResourceDetails(map[string]interface{}{
			"resource_id": input.ParentID,
			"row_status":  1,
		})
		if err != nil {
			return nil, fmt.Errorf("parent details not found: %w", err)
		}
		accountResource.ParentResourceID = resourceData.ParentResourceID
	}

	if input.TenantID != uuid.Nil {
		//check if ParentID is valid
		resourceData, err := dao.GetResourceDetails(map[string]interface{}{
			"resource_id":      input.TenantID,
			"resource_type_id": parsedTenantResourceTypeID,
			"row_status":       1,
		})
		if err != nil {
			return nil, fmt.Errorf("tenant details not found: %w", err)
		}
		accountResource.TenantID = &resourceData.ResourceID
	}

	// Save accountResource to the database
	if err := r.DB.Create(&accountResource).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant resource: %w", err)
	}

	// Prepare metadata (ContactInfo)
	metadata := map[string]interface{}{
		"description": input.Description,
		"billingInfo": input.BillingInfo,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create a new TenantMetadata
	accountMetadata := &dto.TenantMetadata{
		ResourceID: accountResource.ResourceID.String(),
		Metadata:   metadataJSON,
		CreatedBy:  userID.String(),
		RowStatus:  1,
	}

	// Save accountMetadata to the database
	if err := r.DB.Create(&accountMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant metadata: %w", err)
	}

	// Return the created Tenant object
	return &models.Account{
		ID:          accountResource.ResourceID,
		Name:        accountResource.Name,
		Description: input.Description,
		BillingInfo: &models.BillingInfo{
			CreditCardNumber: input.BillingInfo.CreditCardNumber,
			CreditCardType:   input.BillingInfo.CreditCardType,
			ExpirationDate:   input.BillingInfo.ExpirationDate,
			Cvv:              input.BillingInfo.Cvv,
			ID:               accountResource.ResourceID,
			BillingAddress: &models.BillingAddress{
				Street:  input.BillingInfo.BillingAddress.Street,
				City:    input.BillingInfo.BillingAddress.City,
				State:   input.BillingInfo.BillingAddress.State,
				Zipcode: input.BillingInfo.BillingAddress.Zipcode,
				Country: input.BillingInfo.BillingAddress.Country,
			},
		},
		CreatedAt: accountResource.CreatedAt.String(),
	}, nil

}

// CreateAccount resolver for adding a new Account
func (r *AccountMutationResolver) UpdateAccount(ctx context.Context, input models.UpdateAccountInput) (*models.Account, error) {
	userID := uuid.New()
	var accountResource dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{ResourceID: input.ID}).First(&accountResource).Error; err != nil {
		return nil, fmt.Errorf("tenant resource not found: %w", err)
	}

	// Update TenantResource fields if provided
	if input.Name != nil && *input.Name != "" {
		accountResource.Name = *input.Name
	}
	if *input.ParentID != uuid.Nil {
		// Validate ParentOrgID
		var parentOrg dto.TenantResource
		if err := r.DB.Where(&dto.TenantResource{ResourceID: *input.ParentID}).First(&parentOrg).Error; err != nil {
			return nil, fmt.Errorf("parent organization not found: %w", err)
		}
		accountResource.ParentResourceID = &parentOrg.ResourceID
	}

	accountResource.UpdatedBy = userID.String()

	// Save updated TenantResource to the database
	if err := r.DB.Where(&dto.TenantResource{ResourceID: accountResource.ResourceID}).Updates(&accountResource).Error; err != nil {
		return nil, fmt.Errorf("failed to update tenant resource: %w", err)
	}

	// Fetch the existing TenantMetadata
	var accountMetadata dto.TenantMetadata
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: accountResource.ResourceID.String()}).First(&accountMetadata).Error; err != nil {
		return nil, fmt.Errorf("tenant metadata not found: %w", err)
	}

	// Unmarshal the existing metadata
	metadata := map[string]interface{}{}
	if err := json.Unmarshal(accountMetadata.Metadata, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Update metadata fields if provided
	if input.Description != nil && *input.Description != "" {
		metadata["description"] = *input.Description
	}
	if input.BillingInfo != nil {
		billingInfo := metadata["billingInfo"].(map[string]interface{})
		if input.BillingInfo.CreditCardNumber != "" {
			billingInfo["creditCardNumber"] = input.BillingInfo.CreditCardNumber
		}
		// if input.ContactInfo.PhoneNumber != nil && *input.ContactInfo.PhoneNumber != "" {
		// 	contactInfo["phoneNumber"] = *input.ContactInfo.PhoneNumber
		// }
		// if input.ContactInfo.Address != nil {
		// 	address := contactInfo["address"].(map[string]interface{})
		// 	if input.ContactInfo.Address.Street != nil && *input.ContactInfo.Address.Street != "" {
		// 		address["street"] = *input.ContactInfo.Address.Street
		// 	}
		// 	if input.ContactInfo.Address.City != nil && *input.ContactInfo.Address.City != "" {
		// 		address["city"] = *input.ContactInfo.Address.City
		// 	}
		// 	if input.ContactInfo.Address.State != nil && *input.ContactInfo.Address.State != "" {
		// 		address["state"] = *input.ContactInfo.Address.State
		// 	}
		// 	if input.ContactInfo.Address.ZipCode != nil && *input.ContactInfo.Address.ZipCode != "" {
		// 		address["zipCode"] = *input.ContactInfo.Address.ZipCode
		// 	}
		// 	if input.ContactInfo.Address.Country != nil && *input.ContactInfo.Address.Country != "" {
		// 		address["country"] = *input.ContactInfo.Address.Country
		// 	}
		// }
		metadata["billingInfo"] = billingInfo
	}

	// Marshal the updated metadata back to JSON
	updatedMetadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated metadata: %w", err)
	}
	accountMetadata.Metadata = updatedMetadataJSON
	accountMetadata.UpdatedBy = userID.String()

	// Save updated TenantMetadata to the database
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: accountResource.ResourceID.String()}).Updates(&accountMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to update tenant metadata: %w", err)
	}

	// Return the updated Tenant object
	return &models.Account{
		ID:          accountResource.ResourceID,
		Name:        accountResource.Name,
		Description: ptr.String(accountResource.Name),
		UpdatedAt:   ptr.String(accountResource.UpdatedAt.String()),
	}, nil

}

// CreateAccount resolver for adding a new Account
func (r *AccountMutationResolver) DeleteAccount(ctx context.Context, id uuid.UUID) (bool, error) {
	userID := uuid.New()
	tx := r.DB.Begin()

	var accountResource dto.TenantResource
	var accountMetadata dto.TenantMetadata
	if err := tx.Where(&dto.TenantResource{ResourceID: id}).First(&accountResource).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("account resource not found: %w", err)
		}
		return false, fmt.Errorf("failed to fetch account resource: %w", err)
	}

	// Delete associated TenantMetadata
	if err := tx.Model(&accountMetadata).Where(&dto.TenantMetadata{ResourceID: id.String()}).Updates(map[string]interface{}{"RowStatus": 0, "UpdatedBy": userID}).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to delete tenant metadata: %w", err)
	}

	// Delete the TenantResource
	if err := tx.Model(&accountResource).Where(&dto.TenantResource{ResourceID: id}).Updates(map[string]interface{}{"RowStatus": 0, "UpdatedBy": userID}).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to delete account resource: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return success
	return true, nil

}
