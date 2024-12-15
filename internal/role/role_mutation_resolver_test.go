package role

import (
	"context"
	"testing"
	"time"

	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateRole(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	input := models.RoleInput{
		Name:        "Test Role",
		Description: ptrString("Test Description"),
		RoleType:    models.RoleTypeEnumDefault,
		Version:     ptrString("v1.0"),
	}

	role, err := resolver.CreateRole(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "Test Role", role.Name)
	assert.Equal(t, "Test Description", *role.Description)
	assert.Equal(t, models.RoleTypeEnumDefault, role.RoleType)
	assert.Equal(t, "v1.0", *role.Version)
}

func TestCreateRole_MissingName(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	input := models.RoleInput{
		Name:        "",
		Description: ptrString("Test Description"),
		RoleType:    models.RoleTypeEnumDefault,
		Version:     ptrString("v1.0"),
	}

	role, err := resolver.CreateRole(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, role)
	assert.Equal(t, "role name is required", err.Error())
}

func TestUpdateRole(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	// Seed a role in the database
	seedRole := &dto.Role{
		RoleID:      "role_123",
		Name:        "Old Role",
		Description: "Old Description",
		RoleType:    "DEFAULT",
		Version:     "v1.0",
		CreatedAt:   time.Now(),
	}
	db.Create(seedRole)

	input := models.RoleInput{
		Name:        "Updated Role",
		Description: ptrString("Updated Description"),
		RoleType:    models.RoleTypeEnumCustom,
		Version:     ptrString("v2.0"),
	}

	updatedRole, err := resolver.UpdateRole(ctx, seedRole.RoleID, input)
	assert.NoError(t, err)
	assert.NotNil(t, updatedRole)
	assert.Equal(t, "Updated Role", updatedRole.Name)
	assert.Equal(t, "Updated Description", *updatedRole.Description)
	assert.Equal(t, models.RoleTypeEnumCustom, updatedRole.RoleType)
	assert.Equal(t, "v2.0", *updatedRole.Version)
}

func TestUpdateRole_NotFound(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	input := models.RoleInput{
		Name:        "Updated Role",
		Description: ptrString("Updated Description"),
		RoleType:    models.RoleTypeEnumCustom,
		Version:     ptrString("v2.0"),
	}

	updatedRole, err := resolver.UpdateRole(ctx, "nonexistent_role", input)
	assert.Error(t, err)
	assert.Nil(t, updatedRole)
	assert.Equal(t, "role not found", err.Error())
}

func TestDeleteRole(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	// Seed a role in the database
	seedRole := &dto.Role{
		RoleID:    "role_123",
		Name:      "Role to Delete",
		CreatedAt: time.Now(),
	}
	db.Create(seedRole)

	deleted, err := resolver.DeleteRole(ctx, seedRole.RoleID)
	assert.NoError(t, err)
	assert.True(t, deleted)

	// Verify the role is deleted
	var deletedRole dto.Role
	result := db.First(&deletedRole, "role_id = ?", seedRole.RoleID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

// Helper function to create a pointer to a string
func ptrString(s string) *string {
	return &s
}
