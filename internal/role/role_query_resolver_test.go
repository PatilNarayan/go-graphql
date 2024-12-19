package role

import (
	"context"
	"testing"
	"time"

	"go_graphql/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestRoles(t *testing.T) {
	db := setupTestDB()
	resolver := RoleQueryResolver{DB: db}

	db.Create(&dto.Role{RoleID: "1", Name: "Admin", Description: "Administrator role", RoleType: "DEFAULT", Version: "0.0.1", CreatedBy: "1", UpdatedBy: "1", PermissionsIDs: "[]", ResourceID: "1", UpdatedAt: time.Now(), CreatedAt: time.Now()})
	db.Create(&dto.Role{RoleID: "2", Name: "User", Description: "User role", RoleType: "DEFAULT", Version: "0.0.1", CreatedBy: "1", UpdatedBy: "1", PermissionsIDs: "[]", ResourceID: "1", UpdatedAt: time.Now(), CreatedAt: time.Now()})

	ctx := context.Background()

	// Test
	roles, err := resolver.Roles(ctx)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)
	assert.Equal(t, "Admin", roles[0].Name)
	assert.Equal(t, "User", roles[1].Name)
}

func TestGetRole_Success(t *testing.T) {
	db := setupTestDB()
	resolver := RoleQueryResolver{DB: db}

	// Seed data
	roleDB := &dto.Role{RoleID: "1", Name: "Admin", Description: "Administrator role", RoleType: "DEFAULT", Version: "0.0.1", CreatedBy: "1", UpdatedBy: "1", PermissionsIDs: "[]", ResourceID: "1", UpdatedAt: time.Now(), CreatedAt: time.Now()}
	db.Create(roleDB)

	ctx := context.Background()

	// Test
	role, err := resolver.GetRole(ctx, roleDB.RoleID)
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "Admin", role.Name)
	assert.Equal(t, "Administrator role", *role.Description)
}

func TestGetRole_NotFound(t *testing.T) {
	db := setupTestDB()
	resolver := RoleQueryResolver{DB: db}

	ctx := context.Background()

	// Test
	role, err := resolver.GetRole(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Nil(t, role)
	assert.Equal(t, "role not found", err.Error())
}
