package role

import (
	"context"
	"go_graphql/config"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"go_graphql/logger"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

// Setup the in-memory test database
func setupTestDB() *gorm.DB {
	logger.InitLogger()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Auto-migrate tables
	db.AutoMigrate(&dto.TNTRole{}, &dto.TNTPermission{}, &dto.TenantResource{}, &dto.Mst_ResourceTypes{}, &dto.TNTRolePermission{}, &dto.MstRole{}, &dto.MstPermission{}, &dto.MstRolePermission{})

	return db
}

func TestMain(m *testing.M) {
	logger.InitLogger()
	//set environment variables
	os.Setenv("PERMIT_PROJECT", "test")
	os.Setenv("PERMIT_ENV", "test")
	os.Setenv("PERMIT_TOKEN", "test")
	os.Setenv("PERMIT_PDP_ENDPOINT", "https://localhost:8080")

	m.Run()
}

func TestCreateRole(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	ctx := context.Background()
	config.DB = db
	resolver := &RoleMutationResolver{DB: db}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Default responder for unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, `{"error": "unmocked request"}`))

	// Register the mock responder for the API endpoint
	httpmock.RegisterResponder("POST", "https://localhost:8080/v2/schema/test/test/resources/Existing%20Tenant/roles",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `
				{
					"key": "Role_123",
					"name": "Role",
					"status": "success"
				}
				`)
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Seed initial data
	mstResType := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Role",
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&mstResType)

	tenantID := uuid.New()
	existingTenant := dto.TenantResource{
		ResourceID:     tenantID,
		Name:           "Existing Tenant",
		ResourceTypeID: mstResType.ResourceTypeID,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&existingTenant)

	t.Run("Valid Input", func(t *testing.T) {
		input := models.CreateRoleInput{
			Name:               "AdminRole",
			Description:        ptr.String("Admin Description"),
			AssignableScopeRef: tenantID,
			RoleType:           "CUSTOM",
			Version:            "1.0",
			Permissions:        []string{},
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, input.Name, role.Name)
		assert.Equal(t, input.Description, role.Description)

		// Verify role permissions were created
		var rolePermissions []dto.TNTRolePermission
		err = db.Where("role_id = ?", role.ID).Find(&rolePermissions).Error
		assert.NoError(t, err)
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		input := models.CreateRoleInput{
			Description:        ptr.String("Missing Name"),
			AssignableScopeRef: tenantID,
			Permissions:        []string{},
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Equal(t, "role name is required", err.Error())
	})

	t.Run("Missing Tenant ID", func(t *testing.T) {
		input := models.CreateRoleInput{
			Name:        "TestRole",
			Permissions: []string{},
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Equal(t, "resource ID is required", err.Error())
	})

}

func TestUpdateRole(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	config.DB = db
	resolver := &RoleMutationResolver{DB: db}

	// Seed initial data
	mstResType := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Role",
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	mstResTypeTenant := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Tenant",
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&mstResType)
	db.Create(&mstResTypeTenant)

	tenantID := uuid.New()
	existingTenant := dto.TenantResource{
		ResourceID:     tenantID,
		Name:           "Existing Tenant",
		ResourceTypeID: mstResTypeTenant.ResourceTypeID,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&existingTenant)

	roleID := uuid.New()
	existingRole := dto.TenantResource{
		ResourceID:     roleID,
		Name:           "Existing Tenant",
		ResourceTypeID: mstResType.ResourceTypeID,
		TenantID:       &tenantID,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&existingRole)

	role := dto.TNTRole{
		ResourceID: roleID,
		Name:       "Existing Role",
		RoleType:   "CUSTOM",
		Version:    "1.0",
		CreatedBy:  "admin",
		UpdatedBy:  "admin",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	db.Create(&role)
	t.Run("Valid Update", func(t *testing.T) {
		input := models.UpdateRoleInput{
			ID:                 roleID,
			Name:               "Updated Role",
			Description:        ptr.String("Updated Description"),
			AssignableScopeRef: tenantID,
			RoleType:           "CUSTOM",
			Version:            "2.0",
			Permissions:        []string{},
		}

		updatedRole, err := resolver.UpdateRole(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, input.Name, updatedRole.Name)
		assert.Equal(t, input.Description, updatedRole.Description)

		// Verify permissions were updated
		var rolePermissions []dto.TNTRolePermission
		err = db.Where("role_id = ? AND deleted_at IS NULL", roleID).Find(&rolePermissions).Error
		assert.NoError(t, err)
	})

	t.Run("Role Not Found", func(t *testing.T) {
		input := models.UpdateRoleInput{
			ID:                 uuid.New(),
			Name:               "Non-existent Role",
			AssignableScopeRef: tenantID,
			Permissions:        []string{},
		}

		updatedRole, err := resolver.UpdateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
		assert.Equal(t, "role not found", err.Error())
	})

	t.Run("Missing Tenant ID", func(t *testing.T) {
		input := models.UpdateRoleInput{
			ID:          roleID,
			Name:        "Updated Role",
			Permissions: []string{},
		}

		updatedRole, err := resolver.UpdateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
		assert.Equal(t, "Tenant ID is required", err.Error())
	})
}

func TestDeleteRole(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	config.DB = db
	resolver := &RoleMutationResolver{DB: db}

	// Seed initial data
	mstResType := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Role",
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	mstResTypeTenant := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Tenant",
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&mstResType)
	db.Create(&mstResTypeTenant)

	tenantID := uuid.New()
	existingTenant := dto.TenantResource{
		ResourceID:     tenantID,
		Name:           "Existing Tenant",
		ResourceTypeID: mstResTypeTenant.ResourceTypeID,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&existingTenant)

	roleID := uuid.New()
	existingRole := dto.TenantResource{
		ResourceID:     roleID,
		Name:           "Existing Tenant",
		ResourceTypeID: mstResType.ResourceTypeID,
		TenantID:       &tenantID,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&existingRole)

	role := dto.TNTRole{
		ResourceID: roleID,
		Name:       "Existing Role",
		RoleType:   "CUSTOM",
		Version:    "1.0",
		CreatedBy:  "admin",
		UpdatedBy:  "admin",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	db.Create(&role)

	t.Run("Valid Delete", func(t *testing.T) {
		success, err := resolver.DeleteRole(ctx, roleID)
		assert.NoError(t, err)
		assert.True(t, success)

		// Verify role is marked as deleted
		var deletedRole dto.TNTRole
		err = db.Unscoped().First(&deletedRole, "resource_id = ?", roleID).Error
		assert.NoError(t, err)
		assert.NotNil(t, deletedRole.DeletedAt)
		assert.Equal(t, 0, deletedRole.RowStatus)
	})

	t.Run("Delete Non-existent Role", func(t *testing.T) {
		success, err := resolver.DeleteRole(ctx, uuid.New())
		assert.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, "role not found", err.Error())
	})
}
