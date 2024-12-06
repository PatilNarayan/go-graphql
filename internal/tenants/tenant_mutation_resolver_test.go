package tenants

import (
	"context"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockPermitClient struct {
	mock.Mock
}

func (m *MockPermitClient) APIExecute(ctx context.Context, method, endpoint string, body interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, method, endpoint, body)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&dto.Tenant{})
	return db
}

func TestCreateTenant(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	mockPermit := new(MockPermitClient)
	resolver := &TenantMutationResolver{DB: db}

	input := models.TenantInput{
		Name:        "Test Tenant",
		Description: ptr.String("Test Tenant Description"),
		ParentOrgID: "org_123",
	}

	mockPermit.On("APIExecute", ctx, "POST", "tenants", mock.Anything).Return(map[string]interface{}{
		"key": "tenant_123",
	}, nil)

	tenant, err := resolver.CreateTenant(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, "Test Tenant", tenant.Name)
	assert.Equal(t, "org_123", tenant.ParentOrgID)
}

func TestCreateTenant_MissingName(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	// mockPermit := new(MockPermitClient)
	resolver := &TenantMutationResolver{DB: db}

	input := models.TenantInput{
		Name: "",
	}

	tenant, err := resolver.CreateTenant(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, tenant)
	assert.Equal(t, "name is required", err.Error())
}

func TestUpdateTenant(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	mockPermit := new(MockPermitClient)
	resolver := &TenantMutationResolver{DB: db}

	// Seed the database with a tenant
	tenant := &dto.Tenant{ID: "tenant_123", Name: "Old Name", ParentOrgID: "org_123"}
	db.Create(tenant)

	input := models.TenantInput{
		Name:        "Updated Name",
		ParentOrgID: "org_456",
	}

	mockPermit.On("APIExecute", ctx, "PATCH", "tenants/"+tenant.ID, mock.Anything).Return(nil, nil)

	updatedTenant, err := resolver.UpdateTenant(ctx, "tenant_123", input)
	assert.NoError(t, err)
	assert.NotNil(t, updatedTenant)
	assert.Equal(t, "Updated Name", updatedTenant.Name)
	assert.Equal(t, "org_456", updatedTenant.ParentOrgID)
}

func TestDeleteTenant(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	mockPermit := new(MockPermitClient)
	resolver := &TenantMutationResolver{DB: db}

	// Seed the database with a tenant
	tenant := &dto.Tenant{ID: "tenant_123", Name: "Tenant to Delete", RowStatus: 1}
	db.Create(tenant)

	mockPermit.On("APIExecute", ctx, "DELETE", "tenants/"+tenant.ID, nil).Return(nil, nil)

	deleted, err := resolver.DeleteTenant(ctx, "tenant_123")
	assert.NoError(t, err)
	assert.True(t, deleted)

	var deletedTenant dto.Tenant
	result := db.First(&deletedTenant, "tenant_123")
	assert.Error(t, result.Error)
}
