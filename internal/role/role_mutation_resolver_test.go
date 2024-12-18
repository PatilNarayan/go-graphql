package role

import (
	"context"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"testing"

	"github.com/glebarez/sqlite"
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
	db.AutoMigrate(&dto.Role{}, &dto.RoleAssignment{}, &dto.Permission{})
	return db
}

// TestCreateRole validates the CreateRole function
func TestRole(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	// Seed permissions to validate them later
	permission1 := dto.Permission{PermissionID: "perm_1", Name: "Permission 1"}
	permission2 := dto.Permission{PermissionID: "perm_2", Name: "Permission 2"}
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

	// Check if role assignments are created
	var roleAssignments []dto.RoleAssignment
	db.Find(&roleAssignments, "role_id = ?", role.ID)
	assert.Equal(t, 2, len(roleAssignments))

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
	role, err = resolver.UpdateRole(ctx, role.ID, Update)
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "Updated Role", role.Name)
	assert.Equal(t, "Updated Description", *role.Description)
	assert.Equal(t, models.RoleTypeEnumDefault, role.RoleType)
	assert.Equal(t, "1.0", *role.Version)

	// Check if role assignments are updated
	db.Find(&roleAssignments, "role_id = ?", role.ID)
	assert.Equal(t, 1, len(roleAssignments))

	// Call the resolver function
	// Delete the role
	flag, err := resolver.DeleteRole(ctx, role.ID)
	assert.NoError(t, err)
	assert.NotNil(t, flag)
}

func ptrString(s string) *string {
	return &s
}
