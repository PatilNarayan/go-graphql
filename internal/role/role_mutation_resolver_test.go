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
	db.AutoMigrate(&dto.TNTRole{}, &dto.TNTPermission{})
	return db
}

// TestCreateRole validates the CreateRole function
func TestRole(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	// Seed permissions to validate them later
	permission1 := dto.TNTPermission{PermissionID: "perm_1", Name: "Permission 1"}
	permission2 := dto.TNTPermission{PermissionID: "perm_2", Name: "Permission 2"}
	db.Create(&permission1)
	db.Create(&permission2)

	// Input for the role creation
	input := models.RoleInput{
		Name:           "Test Role",
		Description:    ptrString("Test Role Description"),
		PermissionsIds: []string{"perm_1", "perm_2"},
		ResourceID:     ptrString("resource_123"),
		RoleType:       "CUSTOM",
		Version:        ptr.String("1.0"),
		CreatedBy:      "user_123",
		UpdatedBy:      ptr.String("user_123"),
	}

	// Call the resolver function
	role, err := resolver.CreateRole(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "Test Role", role.Name)

	// Input for the role update
	Update := models.RoleInput{
		Name:           "Updated Role",
		Description:    ptrString("Updated Description"),
		PermissionsIds: []string{"perm_2"},
		ResourceID:     ptrString("resource_123"),
		RoleType:       models.RoleTypeEnumDefault,
		Version:        ptr.String("1.0"),
		CreatedBy:      "user_123",
		UpdatedBy:      ptr.String("user_123"),
	}

	// Call the resolver function
	role, err = resolver.UpdateRole(ctx, uuid.MustParse(*role.ResourceID), Update)
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "Updated Role", role.Name)
	assert.Equal(t, "Updated Description", *role.Description)
	assert.Equal(t, models.RoleTypeEnumDefault, role.RoleType)
	assert.Equal(t, "1.0", *role.Version)

	// Call the resolver function
	// Delete the role
	flag, err := resolver.DeleteRole(ctx, uuid.MustParse(*role.ResourceID))
	assert.NoError(t, err)
	assert.NotNil(t, flag)
}

func ptrString(s string) *string {
	return &s
}

func TestCreateRole(t *testing.T) {
	// Setup test database and resolver
	db := setupTestDB()
	ctx := context.Background()
	resolver := &RoleMutationResolver{DB: db}
	// Seed permissions to validate them later
	permission1 := dto.TNTPermission{PermissionID: "perm_1", Name: "Permission 1"}
	permission2 := dto.TNTPermission{PermissionID: "perm_2", Name: "Permission 2"}
	db.Create(&permission1)
	db.Create(&permission2)

	t.Run("Valid Input", func(t *testing.T) {
		input := models.RoleInput{
			Name:           "Admin",
			Description:    ptr.String("Administrator role"),
			ResourceID:     ptr.String("resource_123"),
			RoleType:       models.RoleTypeEnum("DEFAULT"),
			PermissionsIds: []string{"perm_1", "perm_2"},
			Version:        ptr.String("1.0"),
			CreatedBy:      "user_123",
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, "Admin", role.Name)
	})

	t.Run("Missing Role Name", func(t *testing.T) {
		input := models.RoleInput{
			Description:    ptr.String("No name role"),
			ResourceID:     ptr.String("resource_123"),
			RoleType:       models.RoleTypeEnum("DEFAULT"),
			PermissionsIds: []string{"perm_1", "perm_2"},
			Version:        ptr.String("1.0"),
			CreatedBy:      "user_123",
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Equal(t, "role name is required", err.Error())
	})

	t.Run("Invalid Permissions", func(t *testing.T) {
		input := models.RoleInput{
			Name:           "Admin",
			Description:    ptr.String("Invalid permissions role"),
			ResourceID:     ptr.String("resource_123"),
			RoleType:       models.RoleTypeEnum("DEFAULT"),
			PermissionsIds: []string{"invalid_perm"},
			Version:        ptr.String("1.0"),
			CreatedBy:      "user_123",
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Contains(t, err.Error(), "invalid permissions")
	})

	t.Run("JSON Marshal Error", func(t *testing.T) {
		input := models.RoleInput{
			Name:           "Admin",
			Description:    ptr.String("Invalid JSON"),
			ResourceID:     ptr.String("resource_123"),
			RoleType:       models.RoleTypeEnum("DEFAULT"),
			PermissionsIds: []string{"perm_1", string([]byte{0xff})},
			Version:        ptr.String("1.0"),
			CreatedBy:      "user_123",
		}

		role, err := resolver.CreateRole(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, role)
	})
}

func TestUpdateRole(t *testing.T) {
	// Setup test database and resolver
	db := setupTestDB()
	ctx := context.Background()
	resolver := &RoleMutationResolver{DB: db}

	// Seed the database with a role
	role := &dto.TNTRole{
		ResourceID: "role_123",
		Name:       "Old Role",
		// Description:    "Old Description",
		// PermissionsIDs: `["perm_1", "perm_2"]`,
		RoleType:  "DEFAULT",
		Version:   *ptrString("1.0"),
		UpdatedBy: "user_123",
	}
	db.Create(role)

	t.Run("Valid Input", func(t *testing.T) {
		input := models.RoleInput{
			Name:           "Updated Role",
			Description:    ptr.String("Updated Description"),
			RoleType:       models.RoleTypeEnum("DEFAULT"),
			PermissionsIds: []string{"perm_1", "perm_2"},
			Version:        ptr.String("2.0"),
			CreatedBy:      "user_123",
			UpdatedBy:      ptr.String("user_456"),
		}

		updatedRole, err := resolver.UpdateRole(ctx, uuid.MustParse(role.ResourceID), input)
		assert.NoError(t, err)
		assert.NotNil(t, updatedRole)
		assert.Equal(t, "Updated Role", updatedRole.Name)
	})

	t.Run("Role Not Found", func(t *testing.T) {
		input := models.RoleInput{
			Name:      "Non-existent Role",
			UpdatedBy: ptr.String("user_456"),
		}

		updatedRole, err := resolver.UpdateRole(ctx, uuid.MustParse("non-existent-id"), input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
		assert.Equal(t, "role not found", err.Error())
	})

	t.Run("Missing UpdatedBy", func(t *testing.T) {
		input := models.RoleInput{
			Name: "Missing UpdatedBy",
		}

		updatedRole, err := resolver.UpdateRole(ctx, uuid.MustParse(role.ResourceID), input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
		assert.Equal(t, "updatedBy is required", err.Error())
	})

	t.Run("Invalid RoleType", func(t *testing.T) {
		input := models.RoleInput{
			RoleType:  "invalid_type",
			UpdatedBy: ptr.String("user_456"),
		}

		updatedRole, err := resolver.UpdateRole(ctx, uuid.MustParse(role.ResourceID), input)
		assert.Error(t, err)
		assert.Nil(t, updatedRole)
	})
}

func TestDeleteRole(t *testing.T) {
	// Setup test database and resolver
	db := setupTestDB()
	ctx := context.Background()
	resolver := &RoleMutationResolver{DB: db}

	// Seed roles and permissions to validate them later
	role := &dto.TNTRole{
		ResourceID: "role_123",
		Name:       "Old Role",
		// Description:    "Old Description",
		// PermissionsIDs: `["perm_1", "perm_2"]`,
		RoleType: "DEFAULT",
		// ResourceID:     "resource_123",
		Version:   *ptrString("1.0"),
		UpdatedBy: "user_123",
	}
	db.Create(role)

	t.Run("Valid Role Deletion", func(t *testing.T) {
		// Test deleting an existing role
		var roleDB []dto.TNTRole
		if err := db.Find(&roleDB).Error; err != nil {
			t.Fatal(err) // errors.New("role not found")
		}

		roleID := roleDB[0].ResourceID

		success, err := resolver.DeleteRole(ctx, uuid.MustParse(roleID))
		assert.NoError(t, err)
		assert.True(t, success)

		// Ensure the role is deleted from the database
		var deletedRole dto.TNTRole
		db.First(&deletedRole, "role_id = ?", roleID)
		assert.Equal(t, deletedRole.ResourceID, "") // Role should not exist
	})

	t.Run("Role Does Not Exist", func(t *testing.T) {
		// Test deleting a non-existing role
		roleID := "nonexistent_role"

		success, err := resolver.DeleteRole(ctx, uuid.MustParse(roleID))
		assert.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, "role not found", err.Error())
	})

}
