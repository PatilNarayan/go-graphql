package role

import (
	"context"
	"log"
	"testing"

	"go_graphql/internal/dto"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&dto.Role{}, &dto.Permission{}, &dto.RoleAssignment{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestRoles(t *testing.T) {
	db := setupTestDB()
	resolver := RoleQueryResolver{DB: db}

	// Seed data
	db.Create(&dto.Role{RoleID: "1", Name: "Admin", Description: "Administrator role", RoleType: "DEFAULT", Version: "0.0.1"})
	db.Create(&dto.Role{RoleID: "2", Name: "User", Description: "User role", RoleType: "DEFAULT", Version: "0.0.1"})

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
	roleDB := &dto.Role{RoleID: "1", Name: "Admin", Description: "Administrator role", RoleType: "DEFAULT", Version: "0.0.1"}
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
