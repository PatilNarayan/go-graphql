package resource

import (
	"gorm.io/gorm"
)

type ResourceQueryResolver struct {
	DB *gorm.DB
}

// func (r *ResourceQueryResolver) GetResource(ctx context.Context, id uuid.UUID) (models.Resource, error) {
// 	return nil, nil
// }
// func (r *ResourceQueryResolver) AllResources(ctx context.Context) ([]models.Resource, error) {
// 	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
// 	if !ok {
// 		return nil, fmt.Errorf("gin context not found in the request")
// 	}

// 	var resources []dto.TenantResources
// 	tenantID, exists := ginCtx.Get("tenantID")
// 	if !exists {
// 		if err := r.DB.Where(&dto.TenantResources{RowStatus: 1}).Find(&resources).Error; err != nil {
// 			return nil, fmt.Errorf("failed to fetch resources: %w", err)
// 		}
// 	} else {
// 		tenantIDStr, ok := tenantID.(string)
// 		if !ok {
// 			return nil, fmt.Errorf("tenantID is not a string")
// 		}
// 		parsedTenantID, err := uuid.Parse(tenantIDStr)
// 		if err != nil {
// 			return nil, fmt.Errorf("error parsing resource type  UUID: %w", err)
// 		}
// 		if err := r.DB.Where(&dto.TenantResources{RowStatus: 1, TenantID: &parsedTenantID}).Find(&resources).Error; err != nil {
// 			return nil, fmt.Errorf("failed to fetch resources: %w", err)
// 		}
// 	}

// 	allResources := make([]models.Resource, 0, len(resources))
// 	for _, resource := range resources {
// 		var metadata dto.TenantMetadata
// 		if err := r.DB.Where(&dto.TenantMetadata{ResourceID: resource.ResourceID, RowStatus: 1}).First(&metadata).Error; err != nil {
// 			return nil, fmt.Errorf("failed to fetch resource metadata: %w", err)
// 		}
// 		var resourceMetadata models.Account
// 		if err := json.Unmarshal(metadata.Metadata, &resourceMetadata); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshall the data: %w", err)
// 		}
// 		switch resource.ResourceTypeID.String() {
// 		case "ed113f30-bbda-11ef-87ea-c03c5946f955":
// 			allResources = append(allResources, &models.Account{
// 				ID:          resource.ResourceID,
// 				Name:        resource.Name,
// 				Description: resourceMetadata.Description,
// 				CreatedAt:   resource.CreatedAt.String(),
// 				BillingInfo: &models.BillingInfo{
// 					CreditCardNumber: resourceMetadata.BillingInfo.CreditCardNumber,
// 					CreditCardType:   resourceMetadata.BillingInfo.CreditCardType,
// 					Cvv:              resourceMetadata.BillingInfo.Cvv,
// 					ExpirationDate:   resourceMetadata.BillingInfo.ExpirationDate,
// 					BillingAddress: &models.BillingAddress{
// 						Street:  resourceMetadata.BillingInfo.BillingAddress.Street,
// 						City:    resourceMetadata.BillingInfo.BillingAddress.City,
// 						State:   resourceMetadata.BillingInfo.BillingAddress.State,
// 						Zipcode: resourceMetadata.BillingInfo.BillingAddress.Zipcode,
// 						Country: resourceMetadata.BillingInfo.BillingAddress.Country,
// 					},
// 				},
// 				CreatedBy: &resource.CreatedBy,
// 				UpdatedAt: ptr.String(resource.UpdatedAt.String()),
// 				UpdatedBy: &resource.UpdatedBy,
// 			})
// 		}

// 	}

// 	return allResources, nil
// }
