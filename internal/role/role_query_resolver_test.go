package role

import (
	"context"
	"testing"
	"time"

	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)
func TestRoles(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	resolver := RoleQueryResolver{DB: db}

	// Seed roles
	role1 := dto.TNTRole{
		ResourceID: uuid.New(),
		Name:       "Admin",
		RoleType:   "DEFAULT",
		Version:    "1.0",
		CreatedBy:  "admin",
		UpdatedBy:  "admin",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	role2 := dto.TNTRole{
		ResourceID: uuid.New(),
		Name:       "User",
		RoleType:   "DEFAULT",
		Version:    "1.0",
		CreatedBy:  "admin",
		UpdatedBy:  "admin",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	db.Create(&role1)
	db.Create(&role2)

	// Test AllRoles
	roles, err := resolver.AllRoles(ctx)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	// Validate returned roles
	expectedNames := []string{"Admin", "User"}
	for i, role := range roles {
		assert.Equal(t, expectedNames[i], role.Name)
		assert.Equal(t, "DEFAULT", string(role.RoleType))
	}
}

func TestGetRole(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	resolver := RoleQueryResolver{DB: db}

	// Seed role
	role := dto.TNTRole{
		ResourceID: uuid.New(),
		Name:       "Admin",
		RoleType:   "CUSTOM",
		Version:    "1.0",
		CreatedBy:  "admin",
		UpdatedBy:  "admin",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	db.Create(&role)

	t.Run("Role Exists", func(t *testing.T) {
		result, err := resolver.GetRole(ctx, role.ResourceID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Admin", result.Name)
		assert.Equal(t, models.RoleTypeEnum("CUSTOM"), result.RoleType)
	})

	t.Run("Role Not Found", func(t *testing.T) {
		randomID := uuid.New()
		result, err := resolver.GetRole(ctx, randomID)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "role not found", err.Error())
	})
}

func TestGetRole_NotFound(t *testing.T) {
	db := setupTestDB()
	resolver := RoleQueryResolver{DB: db}

	ctx := context.Background()

	// Test
	role, err := resolver.GetRole(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, role)
	assert.Equal(t, "role not found", err.Error())
}
