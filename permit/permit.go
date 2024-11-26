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

type PermitClient struct {
	client *permit.Client
	ctx    context.Context
}

// NewPermitClient initializes a new Permit client
func NewPermitClient() (*PermitClient, error) {
	ctx := context.Background()

	// Environment variables
	pdpEndpoint := os.Getenv("PDP_ENDPOINT")
	permitToken := os.Getenv("PERMIT_TOKEN")
	project := os.Getenv("PROJECT")
	env := os.Getenv("ENV")
	DefaultFactsSyncTimeout := 10 * time.Second

	// Config setup
	permitContext := config.NewPermitContext(config.EnvironmentAPIKeyLevel, project, env)
	client := permit.New(config.NewConfigBuilder(permitToken).
		WithPdpUrl(pdpEndpoint).
		WithApiUrl(pdpEndpoint).
		WithContext(permitContext).
		WithLogger(zap.NewExample()).
		WithProxyFactsViaPDP(true).
		WithFactsSyncTimeout(DefaultFactsSyncTimeout).
		Build())

	return &PermitClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// CreateTenant creates a new tenant
func (p *PermitClient) CreateTenant(tenantName, description string) (*models.TenantRead, error) {
	tenantCreate := models.NewTenantCreate(xid.New().String(), tenantName)
	tenantCreate.SetName(tenantName)
	tenantCreate.SetDescription(description)

	tenant, err := p.client.Api.Tenants.Create(p.ctx, *tenantCreate)
	if err != nil {
		return nil, err.(PermitErrors.PermitError)
	}

	return tenant, nil
}

// DeleteTenant deletes an existing tenant by ID
func (p *PermitClient) DeleteTenant(tenantID string) error {
	err := p.client.Api.Tenants.Delete(p.ctx, tenantID)
	if err != nil {
		return err.(PermitErrors.PermitError)
	}
	return nil
}

// UpdateTenant updates an existing tenant's information
func (p *PermitClient) UpdateTenant(tenantID string, tenantName string) (*models.TenantRead, error) {
	tenantUpdate := models.NewTenantUpdate()
	tenantUpdate.SetName(tenantName)

	tenant, err := p.client.Api.Tenants.Update(p.ctx, tenantID, *tenantUpdate)
	if err != nil {
		return nil, err.(PermitErrors.PermitError)
	}

	return tenant, nil
}
