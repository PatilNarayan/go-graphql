package clientorganizationunits

import (
	"context"
	"iam_services_main_v1/gql/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientOrganizationUnitQueryResolver struct {
	DB *gorm.DB
}

func (r *ClientOrganizationUnitQueryResolver) ClientOrganizationUnit(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	return nil, nil
}
func (r *ClientOrganizationUnitQueryResolver) ClientOrganizationUnits(ctx context.Context) (models.OperationResult, error) {
	return nil, nil
}

// func (r *ClientOrganizationUnitQueryResolver) GetClientOrganizationUnit(ctx context.Context, id uuid.UUID) (*models.ClientOrganizationUnit, error) {
// 	logger := log.WithContext(ctx).WithFields(log.Fields{
// 		"className":  "organization_query_resolver",
// 		"methodName": "AllOrganizations",
// 	})
// 	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
// 	if !ok {
// 		logger.Error("gin context not found")
// 		return nil, fmt.Errorf("gin context not found in the request")
// 	}
// 	resourceType, err := dao.GetResourceTypeByName("ClientOrganizationUnit")
// 	if err != nil {
// 		return nil, fmt.Errorf("resource type not found: %w", err)
// 	}

// 	var resource *dto.TenantResources
// 	tenantID, exists := ginCtx.Get("tenantID")
// 	if !exists {
// 		log.Println("without tenantid")
// 		if err := r.DB.Where(&dto.TenantResources{ResourceTypeID: resourceType.ResourceTypeID, RowStatus: 1}).Find(&resource).Error; err != nil {
// 			return nil, fmt.Errorf("failed to fetch tenants: %w", err)
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
// 		log.Println("with tenantid", parsedTenantID)
// 		if err := r.DB.Where(&dto.TenantResources{ResourceTypeID: resourceType.ResourceTypeID, RowStatus: 1, TenantID: &parsedTenantID}).Find(&resource).Error; err != nil {
// 			return nil, fmt.Errorf("failed to fetch tenants: %w", err)
// 		}
// 	}

// 	var resourceMetadata dto.TenantMetadata
// 	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: resource.ResourceID}).First(&resourceMetadata).Error; err != nil {
// 		return nil, errors.New("unable to find resource metadata")
// 	}

// 	var data map[string]interface{}
// 	json.Unmarshal([]byte(resourceMetadata.Metadata), &data)
// 	if err == nil {
// 		updatedAt := data["updated_at"].(string)
// 		// createdBy := data["created_by"].(string)
// 		// updatedBy := data["updated_by"].(string)

// 		_, exists := data["description"]
// 		var description string
// 		if exists {
// 			description = data["description"].(string)
// 		}
// 		unit := &models.ClientOrganizationUnit{
// 			ID:          uuid.MustParse(data["resource_id"].(string)),
// 			Name:        data["name"].(string),
// 			CreatedAt:   data["created_at"].(string),
// 			Description: &description,
// 			UpdatedAt:   updatedAt,
// 			// CreatedBy:   cast.createdBy,
// 			// UpdatedBy:   updatedBy,
// 		}
// 		return unit, nil
// 	}

// 	return nil, errors.New("unable to unmarshal organization details")
// }

// func (r *ClientOrganizationUnitQueryResolver) AllClientOrganizationUnits(ctx context.Context) ([]*models.ClientOrganizationUnit, error) {
// 	logger := log.WithContext(ctx).WithFields(log.Fields{
// 		"className":  "organization_query_resolver",
// 		"methodName": "AllOrganizations",
// 	})

// 	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
// 	if !ok {
// 		return nil, fmt.Errorf("gin context not found in the request")
// 	}
// 	resourceType, err := dao.GetResourceTypeByName("ClientOrganizationUnit")
// 	if err != nil {
// 		return nil, fmt.Errorf("resource type not found: %w", err)
// 	}

// 	var resources []*dto.TenantResources
// 	tenantID, exists := ginCtx.Get("tenantID")
// 	if !exists {
// 		log.Println("without tenantid")
// 		if err := r.DB.Where(&dto.TenantResources{ResourceTypeID: resourceType.ResourceTypeID, RowStatus: 1}).Find(&resources).Error; err != nil {
// 			return nil, fmt.Errorf("failed to fetch tenants: %w", err)
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
// 		log.Println("with tenantid", parsedTenantID)
// 		if err := r.DB.Where(&dto.TenantResources{ResourceTypeID: resourceType.ResourceTypeID, RowStatus: 1, TenantID: &parsedTenantID}).Find(&resources).Error; err != nil {
// 			return nil, fmt.Errorf("failed to fetch tenants: %w", err)
// 		}
// 	}

// 	var orgs []*models.ClientOrganizationUnit
// 	for _, resource := range resources {
// 		if resource.RowStatus == 1 {
// 			var data map[string]interface{}
// 			resourceStr, _ := json.Marshal(resource)

// 			var resourceMetadata *dto.TenantMetadata
// 			if err := r.DB.Find(&resourceMetadata).Where(&dto.TenantMetadata{ResourceID: resource.ResourceID}).Error; err != nil {
// 				logger.Error("unable to find resource metadata")
// 				return nil, errors.New("unable to find metadata for resource")
// 			}
// 			resourceUnmarshalErr := json.Unmarshal(resourceStr, &data)
// 			var metadataMap map[string]interface{}
// 			metdataUnmarshalErr := json.Unmarshal([]byte(resourceMetadata.Metadata), &metadataMap)

// 			if resourceUnmarshalErr == nil && metdataUnmarshalErr == nil {
// 				// updatedAt := data["updated_at"].(string)
// 				// createdBy := data["created_by"].(string)
// 				// updatedBy := data["updated_by"].(string)
// 				_, exists := metadataMap["description"]
// 				var desciption string
// 				if exists {
// 					desciption = metadataMap["description"].(string)
// 				}
// 				unit := &models.ClientOrganizationUnit{
// 					ID:          uuid.MustParse(data["resource_id"].(string)),
// 					Name:        data["name"].(string),
// 					Description: &desciption,
// 					CreatedAt:   data["created_at"].(string),
// 					// UpdatedAt:   &updatedAt,
// 					// CreatedBy:   &createdBy,
// 					// UpdatedBy:   &updatedBy,
// 				}
// 				orgs = append(orgs, unit)
// 			}
// 		}

// 	}
// 	return orgs, nil
// }
