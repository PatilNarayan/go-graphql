package tenants

import (
	"context"
	"fmt"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"go_graphql/permit"

	"gorm.io/gorm"
)

type TenantMutationResolver struct {
	DB *gorm.DB
}

// CreateTenant resolver for adding a new Tenant
func (r *TenantMutationResolver) CreateTenant(ctx context.Context, input models.TenantInput) (*dto.Tenant, error) {

	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	tenantDB := &dto.Tenant{Name: input.Name, ParentOrgID: input.ParentOrgID, RowStatus: 1}
	if input.Description != nil {
		tenantDB.Description = *input.Description
	}
	if input.ContactInfoID != "" {
		tenantDB.ContactInfoID = input.ContactInfoID
	}
	if err := r.DB.Create(tenantDB).Error; err != nil {
		return nil, err
	}

	pc := permit.NewPermitClient()

	tenant, err := pc.APIExecute(ctx, "POST", "tenants", map[string]interface{}{
		"name": input.Name,
		"key":  tenantDB.ID,
	})

	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, fmt.Errorf("Tenant not created")
	}

	// CreateResource := map[string]interface{}{
	// 	"resource": xid.New().String(),
	// 	"tenant":   tenant["key"],
	// }
	// resourceKeyList := []string{}
	// for _, v := range resourceKeyList {
	// 	CreateResource["key"] = v
	// 	_, err = pc.APIExecute(ctx, "POST", "resources", CreateResource)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// }

	return tenantDB, nil
}

// UpdateTenant resolver for updating a Tenant
func (r *TenantMutationResolver) UpdateTenant(ctx context.Context, id string, input models.TenantInput) (*dto.Tenant, error) {
	var Tenant *dto.Tenant
	if err := r.DB.First(&Tenant, id).Error; err != nil {
		return nil, err
	}

	if Tenant != nil {
		pc := permit.NewPermitClient()
		_, err := pc.APIExecute(ctx, "PATCH", "tenants/"+Tenant.ID, map[string]interface{}{
			"name": input.Name,
		})
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Tenant not found")
	}

	Tenant.Name = input.Name
	Tenant.ParentOrgID = input.ParentOrgID

	if err := r.DB.Save(&Tenant).Error; err != nil {
		return nil, err
	}

	return Tenant, nil
}

// DeleteTenant resolver for deleting a Tenant
func (r *TenantMutationResolver) DeleteTenant(ctx context.Context, id string) (bool, error) {
	var tenant *dto.Tenant
	if err := r.DB.First(&tenant, id).Error; err != nil {
		return false, err
	}
	if tenant == nil {
		return false, fmt.Errorf("Tenant not found")
	}
	pc := permit.NewPermitClient()
	_, err := pc.APIExecute(ctx, "DELETE", "tenants/"+tenant.ID, nil)
	if err != nil {
		return false, err
	}

	tenant.RowStatus = 0
	if err := r.DB.Save(&tenant).Error; err != nil {
		return false, err
	}

	if err := r.DB.Delete(&tenant).Error; err != nil {
		return false, err
	}
	return true, nil
}
