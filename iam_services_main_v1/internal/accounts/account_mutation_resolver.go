package accounts

import (
	"context"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
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

	// userID := uuid.New()
	// TenantResourceTypeID := "ed113bda-bbda-11ef-87ea-c03c5946f955"
	// parsedTenantResourceTypeID, err := uuid.Parse(TenantResourceTypeID)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing resource type  UUID: %w", err)
	// }
	// resourceID := uuid.New()
	// accountResource := &dto.TenantResources{
	// 	ResourceID:       resourceID,
	// 	Name:             input.Name,
	// 	ParentResourceID: &input.ParentID,
	// 	TenantID:         &input.TenantID,
	// 	RowStatus:        1,
	// 	CreatedBy:        userID.String(),
	// }
	// resourceType, err := dao.GetResourceTypeByName("Account")
	// if err != nil {
	// 	return nil, fmt.Errorf("resource type not found: %w", err)
	// }
	// accountResource.ResourceTypeID = resourceType.ResourceTypeID
	// if input.ParentID != uuid.Nil {
	// 	//check if ParentID is valid
	// 	resourceData, err := dao.GetResourceItem(r.DB, map[string]interface{}{
	// 		"resource_id": input.ParentID,
	// 		"row_status":  1,
	// 	})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("parent details not found: %w", err)
	// 	}
	// 	accountResource.ParentResourceID = resourceData.ParentResourceID
	// }

	// if input.TenantID != uuid.Nil {
	// 	//check if ParentID is valid
	// 	resourceData, err := dao.GetResourceItem(r.DB, map[string]interface{}{
	// 		"resource_id":      input.TenantID,
	// 		"resource_type_id": parsedTenantResourceTypeID,
	// 		"row_status":       1,
	// 	})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("tenant details not found: %w", err)
	// 	}
	// 	accountResource.TenantID = &resourceData.ResourceID
	// }

	// // Prepare metadata (ContactInfo)
	// metadata := map[string]interface{}{
	// 	"name":        input.Name,
	// 	"description": input.Description,
	// 	"billingInfo": map[string]interface{}{
	// 		"creditCardNumber": input.BillingInfo.CreditCardNumber,
	// 		"creditCardType":   input.BillingInfo.CreditCardType,
	// 		"expirationDate":   input.BillingInfo.ExpirationDate,
	// 		"cvv":              input.BillingInfo.Cvv,
	// 		"billingAddress": map[string]interface{}{
	// 			"street":  input.BillingInfo.BillingAddress.Street,
	// 			"city":    input.BillingInfo.BillingAddress.City,
	// 			"state":   input.BillingInfo.BillingAddress.State,
	// 			"country": input.BillingInfo.BillingAddress.Country,
	// 			"zipcode": input.BillingInfo.BillingAddress.Zipcode,
	// 		},
	// 	},
	// }
	// metadataJSON, err := json.Marshal(metadata)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	// }

	// // Create a new TenantMetadata
	// accountMetadata := &dto.TenantMetadata{
	// 	ResourceID: resourceID.String(),
	// 	Metadata:   metadataJSON,
	// 	CreatedBy:  userID.String(),
	// 	RowStatus:  1,
	// }
	// pc := permit.NewPermitClient()
	// _, err = pc.APIExecute(ctx, "POST", "resource_instances", map[string]interface{}{
	// 	"key":        resourceID,
	// 	"resource":   resourceType.ResourceTypeID,
	// 	"tenant":     "1588d0ab-a634-4264-a0fb-7a0a8cf6a97a",
	// 	"attributes": metadata,
	// })

	// if err != nil {
	// 	return nil, err
	// }
	// err = r.DB.Transaction(func(tx *gorm.DB) error {
	// 	if err := tx.Create(&accountResource).Error; err != nil {
	// 		return fmt.Errorf("failed to create account  :%w", err)
	// 	}

	// 	if err := tx.Create(&accountMetadata).Error; err != nil {
	// 		return fmt.Errorf("failed to create account metadata :%w", err)
	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create account :%w", err)
	// }

	// // Return the created Tenant object
	// return &models.Account{
	// 	ID:          accountResource.ResourceID,
	// 	Name:        accountResource.Name,
	// 	Description: input.Description,
	// 	CreatedAt:   accountResource.CreatedAt.String(),
	// }, nil

	return &models.Account{}, nil
}

// CreateAccount resolver for adding a new Account
func (r *AccountMutationResolver) UpdateAccount(ctx context.Context, input models.UpdateAccountInput) (*models.Account, error) {
	// userID := uuid.New()
	// var accountResource dto.TenantResources
	// if err := r.DB.Where(&dto.TenantResources{ResourceID: input.ID}).First(&accountResource).Error; err != nil {
	// 	return nil, fmt.Errorf("tenant resource not found: %w", err)
	// }

	// // Update TenantResource fields if provided
	// if input.Name != nil && *input.Name != "" {
	// 	accountResource.Name = *input.Name
	// }
	// /*if *input.ParentID != uuid.Nil {
	// 	// Validate ParentOrgID
	// 	var parentOrg dto.TenantResources
	// 	if err := r.DB.Where(&dto.TenantResources{ResourceID: *input.ParentID}).First(&parentOrg).Error; err != nil {
	// 		return nil, fmt.Errorf("parent organization not found: %w", err)
	// 	}
	// 	accountResource.ParentResourceID = &parentOrg.ResourceID
	// }*/

	// accountResource.UpdatedBy = userID.String()

	// // Save updated TenantResource to the database
	// if err := r.DB.Where(&dto.TenantResources{ResourceID: accountResource.ResourceID}).Updates(&accountResource).Error; err != nil {
	// 	return nil, fmt.Errorf("failed to update tenant resource: %w", err)
	// }

	// if input.Description != nil || input.BillingInfo != nil {
	// 	// Fetch the existing TenantMetadata
	// 	var accountMetadata dto.TenantMetadata
	// 	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: accountResource.ResourceID.String()}).First(&accountMetadata).Error; err != nil {
	// 		return nil, fmt.Errorf("tenant metadata not found: %w", err)
	// 	}

	// 	// Unmarshal the existing metadata
	// 	metadata := map[string]interface{}{}
	// 	if err := json.Unmarshal(accountMetadata.Metadata, &metadata); err != nil {
	// 		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	// 	}
	// 	newMetadata := make(map[string]interface{})

	// 	if input.Description != nil {
	// 		newMetadata["description"] = input.Description
	// 	}

	// 	if input.BillingInfo != nil {
	// 		billingInfoMap := helpers.StructToMap(input.BillingInfo)
	// 		if len(billingInfoMap) > 0 {
	// 			newMetadata["billingInfo"] = billingInfoMap
	// 		}
	// 	}
	// 	mergedMetadata := helpers.MergeMaps(metadata, newMetadata)
	// 	serializedMetadata, err := json.Marshal(mergedMetadata)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to serialize merged metadata: %w", err)
	// 	}
	// 	accountMetadata.Metadata = serializedMetadata
	// 	accountMetadata.UpdatedBy = userID.String()
	// 	// Save updated TenantMetadata to the database
	// 	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: accountResource.ResourceID.String()}).Updates(&accountMetadata).Error; err != nil {
	// 		return nil, fmt.Errorf("failed to update tenant metadata: %w", err)
	// 	}
	// }

	// // Return the updated Tenant object
	// return &models.Account{
	// 	ID:          accountResource.ResourceID,
	// 	Name:        accountResource.Name,
	// 	Description: ptr.String(accountResource.Name),
	// 	UpdatedAt:   ptr.String(accountResource.UpdatedAt.String()),
	// }, nil

	return &models.Account{}, nil
}

// CreateAccount resolver for adding a new Account
func (r *AccountMutationResolver) DeleteAccount(ctx context.Context, id uuid.UUID) (bool, error) {
	// userID := uuid.New()
	// tx := r.DB.Begin()

	// var accountResource dto.TenantResources
	// var accountMetadata dto.TenantMetadata
	// if err := tx.Where(&dto.TenantResources{ResourceID: id}).First(&accountResource).Error; err != nil {
	// 	tx.Rollback()
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return false, fmt.Errorf("account resource not found: %w", err)
	// 	}
	// 	return false, fmt.Errorf("failed to fetch account resource: %w", err)
	// }

	// // Delete associated TenantMetadata
	// if err := tx.Model(&accountMetadata).Where(&dto.TenantMetadata{ResourceID: id.String()}).Updates(map[string]interface{}{"RowStatus": 0, "UpdatedBy": userID}).Error; err != nil {
	// 	tx.Rollback()
	// 	return false, fmt.Errorf("failed to delete tenant metadata: %w", err)
	// }

	// // Delete the TenantResource
	// if err := tx.Model(&accountResource).Where(&dto.TenantResources{ResourceID: id}).Updates(map[string]interface{}{"RowStatus": 0, "UpdatedBy": userID}).Error; err != nil {
	// 	tx.Rollback()
	// 	return false, fmt.Errorf("failed to delete account resource: %w", err)
	// }

	// // Commit the transaction
	// if err := tx.Commit().Error; err != nil {
	// 	return false, fmt.Errorf("failed to commit transaction: %w", err)
	// }

	// // Return success
	// return true, nil

	return true, nil
}
