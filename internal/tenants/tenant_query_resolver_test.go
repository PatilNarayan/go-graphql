package tenants

import (
	"context"
	"go_graphql/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to set up a test database.
// func setupTestDB() *gorm.DB {
// 	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	db.AutoMigrate(&dto.Tenant{})
// 	return db
// }

func TestTenants(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	// Seed the database with tenants.
	tenantsList := []dto.Tenant{
		{ID: "1", Name: "Tenant 1", RowStatus: 1},
		{ID: "2", Name: "Tenant 2", RowStatus: 1},
	}
	for _, tenant := range tenantsList {
		db.Create(&tenant)
	}

	resolver := TenantQueryResolver{DB: db}

	result, err := resolver.Tenants(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Tenant 1", result[0].Name)
	assert.Equal(t, "Tenant 2", result[1].Name)
}

func TestTenants_NoRecords(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	resolver := TenantQueryResolver{DB: db}

	result, err := resolver.Tenants(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestGetTenant(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	// Seed the database with a tenant.
	tenant := dto.Tenant{Name: "Tenant 1", RowStatus: 1}
	db.Create(&tenant)

	resolver := TenantQueryResolver{DB: db}

	tenantID := tenant.ID
	result, err := resolver.GetTenant(ctx, tenantID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Tenant 1", result.Name)
}

func TestGetTenant_NotFound(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	resolver := TenantQueryResolver{DB: db}

	tenantID := "999" // ID not present in the database.
	result, err := resolver.GetTenant(ctx, tenantID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTenant_NilID(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	resolver := TenantQueryResolver{DB: db}

	var tenantID string = ""
	result, err := resolver.GetTenant(ctx, tenantID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "id cannot be nil", err.Error())
}
