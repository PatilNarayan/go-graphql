package permit

import (
	"context"
	"os"
	"time"

	"github.com/permitio/permit-golang/pkg/config"
	PermitErrors "github.com/permitio/permit-golang/pkg/errors"
	"github.com/permitio/permit-golang/pkg/models"
	"github.com/permitio/permit-golang/pkg/permit"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

func NewPermitClientPDP() permit.Client {
	project := os.Getenv("PERMIT_PROJECT")
	env := os.Getenv("PERMIT_ENV")
	token := os.Getenv("PERMIT_TOKEN")
	DefaultPDPUrl := os.Getenv("PERMIT_PDP_ENDPOINT")
	permitContext := config.NewPermitContext(config.EnvironmentAPIKeyLevel, project, env)
	permitClient := permit.New(config.NewConfigBuilder(token).
		WithPdpUrl(DefaultPDPUrl).
		WithApiUrl(DefaultPDPUrl).
		WithContext(permitContext).
		WithLogger(zap.NewExample()).
		WithProxyFactsViaPDP(true).
		WithFactsSyncTimeout(10 * time.Second).
		Build())
	return *permitClient
}

func CreateTenant(tenantname string) (*models.TenantRead, error) {
	ctx := context.Background()
	permitClient := NewPermitClientPDP()
	tenantCreate := models.NewTenantCreate(xid.New().String(), tenantname)
	tenantCreate.SetName(tenantname)
	tenant, err := permitClient.Api.Tenants.Create(ctx, *tenantCreate)
	if err != nil {
		return nil, err.(PermitErrors.PermitError)
	}
	return tenant, nil
}

func UpdateTenant(tenantid string, tenantname string) (*models.TenantRead, error) {
	ctx := context.Background()
	permitClient := NewPermitClientPDP()
	tenantUpdate := models.NewTenantUpdate()
	tenantUpdate.SetName(tenantname)
	tenant, err := permitClient.Api.Tenants.Update(ctx, tenantid, *tenantUpdate)
	if err != nil {
		return nil, err.(PermitErrors.PermitError)
	}
	return tenant, nil
}

func DeleteTenant(tenantid string) (bool, error) {
	ctx := context.Background()
	permitClient := NewPermitClientPDP()
	err := permitClient.Api.Tenants.Delete(ctx, "tenantKey")
	if err != nil {
		return false, err.(PermitErrors.PermitError)
	}
	return true, nil
}

func ResourceInstance(tenantKey string, resourceName string) (*models.ResourceInstanceRead, error) {
	ctx := context.Background()
	permitClient := NewPermitClientPDP()
	resourceInstanceCreate := models.ResourceInstanceCreate{
		Key:      xid.New().String(),
		Tenant:   &tenantKey,
		Resource: resourceName,
	}
	resourceInstance, err := permitClient.Api.ResourceInstances.Create(ctx, resourceInstanceCreate)
	if err != nil {
		return nil, err.(PermitErrors.PermitError)
	}
	return resourceInstance, nil
}
