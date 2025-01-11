package role

import (
	"context"
	"go_graphql/config"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

// Setup the in-memory test database
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Auto-migrate tables
	db.AutoMigrate(&dto.TNTRole{}, &dto.TNTPermission{}, &dto.TenantResource{}, &dto.Mst_ResourceTypes{})

	return db
}

func TestCreateRole(t *testing.T) {
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
	db.Create(&mstResType)
	existingTenant := dto.TenantResource{
		ResourceID:     uuid.New(),
		Name:           "Existing Tenant",
		ResourceTypeID: mstResType.ResourceTypeID,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&existingTenant)

	t.Run("Valid Input", func(t *testing.T) {
		parentOrgID := existingTenant.ResourceID
		input := models.RoleInput{
			Name:        "Admin Role",
			Description: ptr.String("Admin Description"),
			ParentOrgID: parentOrgID,
			// RoleType:    "CUSTOM",
			Version: "1.0",
			// CreatedBy:   "test_user",
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, input.Name, role.Name)
		assert.Equal(t, input.Description, role.Description)
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		input := models.RoleInput{
			Description: ptr.String("Missing Name"),
			ParentOrgID: uuid.New(),
			// CreatedBy:   "test_user",
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Equal(t, "role name is required", err.Error())
	})

	t.Run("Missing ParentOrgID", func(t *testing.T) {
		input := models.RoleInput{
			Name: "Test Role",
			// CreatedBy: "test_user",
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
	db.Create(&mstResType)

	existingTenant := dto.TenantResource{
		ResourceID:     uuid.New(),
		Name:           "Existing Tenant",
		ResourceTypeID: mstResType.ResourceTypeID,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&existingTenant)

	// Create initial role for testing updates
	initialRole := dto.TNTRole{
		ResourceID:       uuid.New(),
		ParentResourceID: &existingTenant.ResourceID,
		Name:             "Initial Role",
		RoleType:         "CUSTOM",
		Version:          "1.0",
		CreatedBy:        "admin",
		UpdatedBy:        "admin",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	db.Create(&initialRole)

	t.Run("Valid Update", func(t *testing.T) {
		newParentOrgID := existingTenant.ResourceID
		input := models.RoleInput{
			Name:        "Updated Role",
			Description: ptr.String("Updated Description"),
			ParentOrgID: newParentOrgID,
			// RoleType:    models.RoleTypeEnumDefault,
			Version: "2.0",
		}

		updatedRole, err := resolver.UpdateRole(ctx, initialRole.ResourceID, input)
		assert.NoError(t, err)
		assert.NotNil(t, updatedRole)
		assert.Equal(t, input.Name, updatedRole.Name)
		// assert.Equal(t, string(input.RoleType), string(updatedRole.RoleType))
		assert.Equal(t, input.Version, *updatedRole.Version)
	})

	t.Run("Role Not Found", func(t *testing.T) {
		input := models.RoleInput{
			Name:        "Non-existent Role",
			ParentOrgID: existingTenant.ResourceID,
		}

		updatedRole, err := resolver.UpdateRole(ctx, uuid.New(), input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
		assert.Equal(t, "role not found", err.Error())
	})

	t.Run("Missing ParentOrgID", func(t *testing.T) {
		input := models.RoleInput{
			Name: "Test Update",
			// RoleType: models.RoleTypeEnumDefault,
		}

		updatedRole, err := resolver.UpdateRole(ctx, initialRole.ResourceID, input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
		assert.Equal(t, "resource ID is required", err.Error())
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
	db.Create(&mstResType)

	existingTenant := dto.TenantResource{
		ResourceID:     uuid.New(),
		Name:           "Existing Tenant",
		ResourceTypeID: mstResType.ResourceTypeID,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&existingTenant)

	// Create a role to delete
	role := dto.TNTRole{
		ResourceID:       uuid.New(),
		ParentResourceID: &existingTenant.ResourceID,
		Name:             "Role to Delete",
		RoleType:         "CUSTOM",
		Version:          "1.0",
		CreatedBy:        "admin",
		UpdatedBy:        "admin",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	db.Create(&role)

	t.Run("Valid Delete", func(t *testing.T) {
		success, err := resolver.DeleteRole(ctx, role.ResourceID)
		assert.NoError(t, err)
		assert.True(t, success)

		// Verify deletion
		var deletedRole dto.TNTRole
		result := db.First(&deletedRole, "resource_id = ?", role.ResourceID)
		assert.Error(t, result.Error)
		assert.True(t, result.Error == gorm.ErrRecordNotFound)
	})

	t.Run("Delete Non-existent Role", func(t *testing.T) {
		success, err := resolver.DeleteRole(ctx, uuid.New())
		assert.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, "role not found", err.Error())
	})
}
