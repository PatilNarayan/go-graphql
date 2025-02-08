package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"go_graphql/internal/dto"

	"gorm.io/gorm"
)

type AccountFieldResolver struct {
	DB *gorm.DB
}

// BillingInfo resolves the BillingInfo field on the Account type
func (r *AccountFieldResolver) BillingInfo(ctx context.Context, obj *models.Account) (*models.BillingInfo, error) {
	var accountMetadata dto.TenantMetadata
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: obj.ID.String(), RowStatus: 1}).First(&accountMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
	}
	var resourceMetadata models.Account
	if err := json.Unmarshal(accountMetadata.Metadata, &resourceMetadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshall the data: %w", err)
	}
	return &models.BillingInfo{
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
	}, nil
}

// BillingInfo resolves the BillingInfo field on the Account type
func (r *AccountFieldResolver) BillingAddress(ctx context.Context, obj *models.Account) (*models.BillingAddress, error) {
	var accountMetadata dto.TenantMetadata
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: obj.ID.String(), RowStatus: 1}).First(&accountMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch parent organization: %w", err)
	}
	var resourceMetadata models.Account
	if err := json.Unmarshal(accountMetadata.Metadata, &resourceMetadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshall the data: %w", err)
	}
	return &models.BillingAddress{
		Street:  resourceMetadata.BillingInfo.BillingAddress.Street,
		City:    resourceMetadata.BillingInfo.BillingAddress.City,
		State:   resourceMetadata.BillingInfo.BillingAddress.State,
		Zipcode: resourceMetadata.BillingInfo.BillingAddress.Zipcode,
		Country: resourceMetadata.BillingInfo.BillingAddress.Country,
	}, nil
}
