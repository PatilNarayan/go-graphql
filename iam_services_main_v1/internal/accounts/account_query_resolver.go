package accounts

import (
	"context"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountQueryResolver struct {
	DB *gorm.DB
}

// Tenants resolver for fetching all Tenants
func (r *AccountQueryResolver) Accounts(ctx context.Context) (models.OperationResult, error) {
	// resourceType, err := dao.GetResourceTypeByName("Account")
	// if err != nil {
	// 	return nil, fmt.Errorf("resource type not found: %w", err)
	// }

	// // tenantID, _ := helpers.GetTenantID(ctx)

	// conditions := map[string]interface{}{
	// 	"row_status":       1,
	// 	"resource_type_id": resourceType.ResourceTypeID,
	// }

	// if tenantID != nil {
	// 	conditions["tenant_id"] = *tenantID
	// }
	// accountResources, err := dao.GetAllResources(r.DB, conditions)
	// if err != nil {
	// 	return nil, fmt.Errorf("error occured when fetching the resources: %w", err)
	// }
	// accounts := make([]*models.Account, 0, len(accountResources))
	// for _, accountResource := range accountResources {
	// 	account := dto.MapToAccount(accountResource)
	// 	accounts = append(accounts, account)
	// }

	// return accounts, nil

	return nil, nil
}

// GetTenant resolver for fetching a single Tenant by ID
func (r *AccountQueryResolver) Account(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	// resourceType, err := dao.GetResourceTypeByName("Account")
	// if err != nil {
	// 	return nil, fmt.Errorf("resource type not found: %w", err)
	// }

	// tenantID, _ := helpers.GetTenantID(ctx)

	// conditions := map[string]interface{}{
	// 	"resource_id":      id,
	// 	"row_status":       1,
	// 	"resource_type_id": resourceType.ResourceTypeID,
	// }

	// if tenantID != nil {
	// 	conditions["tenant_id"] = *tenantID
	// }

	// accountResource, err := dao.GetResourceItem(r.DB, conditions)

	// if err != nil {
	// 	return nil, fmt.Errorf("error occured when fetching the resources: %w", err)
	// }

	// return dto.MapToAccount(*accountResource), nil

	return nil, nil
}
