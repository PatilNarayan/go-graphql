package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type AccountQueryResolver struct {
	DB *gorm.DB
}

// Tenants resolver for fetching all Tenants
func (r *AccountQueryResolver) AllAccounts(ctx context.Context) ([]*models.Account, error) {
	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("gin context not found in the request")
	}
	resourceType, err := dao.GetResourceTypeByName("Account")
	if err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	var accountResources []dto.TenantResource
	tenantID, exists := ginCtx.Get("tenantID")
	if !exists {
		if err := r.DB.Where(&dto.TenantResource{ResourceTypeID: resourceType.ResourceTypeID, RowStatus: 1}).Find(&accountResources).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch tenants: %w", err)
		}
	} else {
		tenantIDStr, ok := tenantID.(string)
		if !ok {
			return nil, fmt.Errorf("tenantID is not a string")
		}
		parsedTenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing resource type  UUID: %w", err)
		}
		if err := r.DB.Where(&dto.TenantResource{ResourceTypeID: resourceType.ResourceTypeID, RowStatus: 1, TenantID: &parsedTenantID}).Find(&accountResources).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch tenants: %w", err)
		}
	}

	accounts := make([]*models.Account, 0, len(accountResources))
	for _, accountResource := range accountResources {
		// Fetch associated metadata
		var accountMetadata dto.TenantMetadata
		if err := r.DB.Where(&dto.TenantMetadata{ResourceID: accountResource.ResourceID.String(), RowStatus: 1}).First(&accountMetadata).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
		}
		var resourceMetadata models.Account
		if err := json.Unmarshal(accountMetadata.Metadata, &resourceMetadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshall the data: %w", err)
		}
		account := &models.Account{
			ID:          accountResource.ResourceID,
			Name:        accountResource.Name,
			Description: resourceMetadata.Description,
			CreatedAt:   accountResource.CreatedAt.String(),
			CreatedBy:   &accountResource.CreatedBy,
			UpdatedAt:   ptr.String(accountResource.UpdatedAt.String()),
			UpdatedBy:   &accountResource.UpdatedBy,
		}
		accounts = append(accounts, account)
	}

	return accounts, nil

}

// GetTenant resolver for fetching a single Tenant by ID
func (r *AccountQueryResolver) GetAccount(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("gin context not found in the request")
	}
	resourceType, err := dao.GetResourceTypeByName("Account")
	if err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	var accountResource dto.TenantResource
	tenantID, exists := ginCtx.Get("tenantID")
	if !exists {
		if err := r.DB.Where(&dto.TenantResource{ResourceTypeID: resourceType.ResourceTypeID, RowStatus: 1}).Find(&accountResource).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch tenants: %w", err)
		}
	} else {
		tenantIDStr, ok := tenantID.(string)
		if !ok {
			return nil, fmt.Errorf("tenantID is not a string")
		}
		parsedTenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing resource type  UUID: %w", err)
		}
		if err := r.DB.Where(&dto.TenantResource{ResourceTypeID: resourceType.ResourceTypeID, RowStatus: 1, TenantID: &parsedTenantID}).Find(&accountResource).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch tenants: %w", err)
		}
	}

	var parentorg *dto.TenantResource
	if accountResource.ParentResourceID != &uuid.Nil {
		if err := r.DB.Where(&dto.TenantResource{ResourceID: *accountResource.ParentResourceID}).First(&parentorg).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
		}
	}
	//account := convertAccountToGraphQL(&accountResources, parentorg)
	// Fetch associated metadata
	var accountMetadata dto.TenantMetadata
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: accountResource.ResourceID.String(), RowStatus: 1}).First(&accountMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
	}
	var resourceMetadata models.Account
	if err := json.Unmarshal(accountMetadata.Metadata, &resourceMetadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshall the data: %w", err)
	}

	return &models.Account{
		ID:          accountResource.ResourceID,
		Name:        accountResource.Name,
		Description: resourceMetadata.Description,
		CreatedAt:   accountResource.CreatedAt.String(),
		BillingInfo: &models.BillingInfo{
			CreditCardNumber: resourceMetadata.BillingInfo.CreditCardNumber,
			CreditCardType:   resourceMetadata.BillingInfo.CreditCardType,
			Cvv:              resourceMetadata.BillingInfo.Cvv,
			ExpirationDate:   resourceMetadata.BillingInfo.ExpirationDate,
			BillingAddress: &models.BillingAddress{
				Street:  resourceMetadata.BillingInfo.BillingAddress.Street,
				City:    resourceMetadata.BillingInfo.BillingAddress.City,
				State:   resourceMetadata.BillingInfo.BillingAddress.State,
				Zipcode: resourceMetadata.BillingInfo.BillingAddress.Zipcode,
				Country: resourceMetadata.BillingInfo.BillingAddress.Country,
			},
		},
		//CreatedBy:   &account.CreatedBy,
		UpdatedAt: ptr.String(accountResource.UpdatedAt.String()),
		//UpdatedBy:   &account.UpdatedBy,
	}, nil
	//return account, nil
}
