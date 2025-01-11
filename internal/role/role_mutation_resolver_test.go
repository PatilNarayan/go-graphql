package role

import (
	"context"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"testing"

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

	// Seed resource type
	resourceType := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		Name:           "Role",
	}
	db.Create(&resourceType)

	return db
}

func TestCreateRole(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	resolver := &RoleMutationResolver{DB: db}

	t.Run("Valid Input", func(t *testing.T) {
		parentOrgID := uuid.New()
		input := models.RoleInput{
			Name:        "Admin Role",
			Description: ptr.String("Admin Description"),
			ParentOrgID: parentOrgID,
			RoleType:    "CUSTOM",
			Version:     "1.0",
			CreatedBy:   "test_user",
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, input.Name, role.Name)
		assert.Equal(t, input.Description, role.Description)
		assert.Equal(t, input.ParentOrgID.String(), role.ParentOrgID)
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		input := models.RoleInput{
			Description: ptr.String("Missing Name"),
			ParentOrgID: uuid.New(),
			CreatedBy:   "test_user",
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Equal(t, "role name is required", err.Error())
	})

	t.Run("Missing ParentOrgID", func(t *testing.T) {
		input := models.RoleInput{
			Name:      "Test Role",
			CreatedBy: "test_user",
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
	resolver := &RoleMutationResolver{DB: db}

	// Create initial role for testing updates
	parentOrgID := uuid.New()
	initialRole := dto.TNTRole{
		ResourceID:       uuid.New(),
		ParentResourceID: &parentOrgID,
		Name:             "Initial Role",
		RoleType:         "CUSTOM",
		Version:          "1.0",
		CreatedBy:        "test_user",
		UpdatedBy:        "test_user",
	}
	db.Create(&initialRole)

	t.Run("Valid Update", func(t *testing.T) {
		newParentOrgID := uuid.New()
		input := models.RoleInput{
			Name:        "Updated Role",
			Description: ptr.String("Updated Description"),
			ParentOrgID: newParentOrgID,
			RoleType:    "DEFAULT",
			Version:     "2.0",
			UpdatedBy:   ptr.String("update_user"),
		}

		updatedRole, err := resolver.UpdateRole(ctx, initialRole.ResourceID, input)
		assert.NoError(t, err)
		assert.NotNil(t, updatedRole)
		assert.Equal(t, input.Name, updatedRole.Name)
		assert.Equal(t, input.RoleType, string(updatedRole.RoleType))
		assert.Equal(t, input.Version, *updatedRole.Version)
	})

	t.Run("Role Not Found", func(t *testing.T) {
		input := models.RoleInput{
			Name:        "Non-existent Role",
			ParentOrgID: uuid.New(),
			UpdatedBy:   ptr.String("update_user"),
		}

		updatedRole, err := resolver.UpdateRole(ctx, uuid.New(), input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
		assert.Equal(t, "role not found", err.Error())
	})

	t.Run("Missing UpdatedBy", func(t *testing.T) {
		input := models.RoleInput{
			Name:        "Test Update",
			ParentOrgID: uuid.New(),
		}

		updatedRole, err := resolver.UpdateRole(ctx, initialRole.ResourceID, input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
		assert.Equal(t, "updatedBy is required", err.Error())
	})
}

func TestDeleteRole(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	resolver := &RoleMutationResolver{DB: db}

	// Create a role to delete
	parentOrgID := uuid.New()
	role := dto.TNTRole{
		ResourceID:       uuid.New(),
		ParentResourceID: &parentOrgID,
		Name:             "Role to Delete",
		RoleType:         "CUSTOM",
		Version:          "1.0",
		CreatedBy:        "test_user",
		UpdatedBy:        "test_user",
	}
	db.Create(&role)

	t.Run("Valid Delete", func(t *testing.T) {
		success, err := resolver.DeleteRole(ctx, role.ResourceID)
		assert.NoError(t, err)
		assert.True(t, success)

		// Verify deletion
		var deletedRole dto.TNTRole
		result := db.First(&deletedRole, "role_id = ?", role.ResourceID)
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
