package role

import (
	"context"
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
	db.AutoMigrate(&dto.Role{}, &dto.RoleAssignment{}, &dto.Permission{})
	return db
}

// TestCreateRole validates the CreateRole function
func TestCreateRole(t *testing.T) {
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
}

// TestUpdateRole validates the UpdateRole function
func TestUpdateRole(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	// Seed a role and permissions
	role := dto.Role{
		RoleID:         uuid.New().String(),
		Name:           "Old Role",
		RoleType:       "user",
		PermissionsIDs: `["perm_1"]`,
		CreatedAt:      time.Now(),
		CreatedBy:      "user_123",
	}
	db.Create(&role)

	permission := dto.Permission{PermissionID: "perm_2", Name: "Permission 2"}
	db.Create(&permission)

	// Update input
	input := models.RoleInput{
		Name:           "Updated Role",
		Description:    ptrString("Updated Description"),
		PermissionsIds: []string{"perm_2"},
		UpdatedBy:      ptrString("user_456"),
	}

	// Call the resolver function
	updatedRole, err := resolver.UpdateRole(ctx, role.RoleID, input)
	assert.NoError(t, err)
	assert.NotNil(t, updatedRole)
	assert.Equal(t, "Updated Role", updatedRole.Name)
	assert.Equal(t, "Updated Description", updatedRole.Description)

	// Check role assignments are updated
	var roleAssignments []dto.RoleAssignment
	db.Find(&roleAssignments, "role_id = ?", role.RoleID)
	assert.Equal(t, 1, len(roleAssignments))
	assert.Equal(t, "perm_2", roleAssignments[0].PermissionID)
}

// TestDeleteRole validates the DeleteRole function
func TestDeleteRole(t *testing.T) {
	db := setupTestDB()
	resolver := &RoleMutationResolver{DB: db}
	ctx := context.Background()

	// Seed a role and assignments
	role := dto.Role{
		RoleID: uuid.New().String(),
		Name:   "Role to Delete",
	}
	db.Create(&role)

	roleAssignment := dto.RoleAssignment{
		RoleAssignmentID: uuid.New().String(),
		RoleID:           role.RoleID,
		PermissionID:     "perm_1",
		CreatedAt:        time.Now(),
		CreatedBy:        "user_123",
	}
	db.Create(&roleAssignment)

	// Call the resolver function
	deleted, err := resolver.DeleteRole(ctx, role.RoleID)
	assert.NoError(t, err)
	assert.True(t, deleted)

	// Ensure role and assignments are deleted
	var roleCount int64
	db.Model(&dto.Role{}).Where("role_id = ?", role.RoleID).Count(&roleCount)
	assert.Equal(t, int64(0), roleCount)

	var assignmentCount int64
	db.Model(&dto.RoleAssignment{}).Where("role_id = ?", role.RoleID).Count(&assignmentCount)
	assert.Equal(t, int64(0), assignmentCount)
}

// Utility functions for pointers
func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}
