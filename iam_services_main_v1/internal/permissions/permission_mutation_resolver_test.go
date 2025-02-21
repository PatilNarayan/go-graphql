package permissions

import (
	"iam_services_main_v1/internal/dto"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&dto.MstPermission{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

// func TestDeletePermission(t *testing.T) {
// 	logger.InitLogger()
// 	db := setupTestDB(t)
// 	resolver := &PermissionMutationResolver{DB: db}

// 	// Insert a sample permission
// 	permissionID := uuid.New()
// 	db.Create(&dto.TNTPermission{
// 		PermissionID: permissionID,
// 		Name:         "Test Permission",
// 		ServiceID:    uuid.New().String(),
// 		Action:       "test_action",
// 		RowStatus:    1,
// 		CreatedAt:    time.Now(),
// 		UpdatedAt:    time.Now(),
// 	})

// 	// Test case 1: Successful deletion
// 	success, err := resolver.DeletePermission(context.Background(), permissionID)
// 	assert.NoError(t, err)
// 	assert.True(t, success)

// 	// Test case 2: Non-existent permission
// 	success, err = resolver.DeletePermission(context.Background(), uuid.New())
// 	assert.Error(t, err)
// 	assert.False(t, success)

// 	// Test case 3: Database error (simulate by closing DB)
// 	sqlDB, _ := db.DB()
// 	sqlDB.Close()
// 	_, err = resolver.DeletePermission(context.Background(), permissionID)
// 	assert.Error(t, err)
// }

// func TestCreatePermission(t *testing.T) {
// 	logger.InitLogger()
// 	db := setupTestDB(t)
// 	resolver := &PermissionMutationResolver{DB: db}

// 	// Test case 1: Successful creation
// 	input := &models.CreatePermission{
// 		Name:      "Test Permission",
// 		ServiceID: uuid.New(),
// 		Action:    "test_action",
// 	}
// 	createdPermission, err := resolver.CreatePermission(context.Background(), input)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, createdPermission)
// 	assert.Equal(t, input.Name, createdPermission.Name)

// 	// Test case 2: Missing input
// 	_, err = resolver.CreatePermission(context.Background(), nil)
// 	assert.Error(t, err)

// 	// Test case 3: Duplicate name
// 	_, err = resolver.CreatePermission(context.Background(), input)
// 	assert.Error(t, err)
// }

// func TestUpdatePermission(t *testing.T) {
// 	logger.InitLogger()
// 	db := setupTestDB(t)
// 	resolver := &PermissionMutationResolver{DB: db}

// 	// Insert a sample permission
// 	permissionID := uuid.New()
// 	db.Create(&dto.TNTPermission{
// 		PermissionID: permissionID,
// 		Name:         "Test Permission",
// 		ServiceID:    uuid.New().String(),
// 		Action:       "test_action",
// 		RowStatus:    1,
// 		CreatedAt:    time.Now(),
// 		UpdatedAt:    time.Now(),
// 	})

// 	serviceID := uuid.New()

// 	// Test case 1: Successful update
// 	input := &models.UpdatePermission{
// 		ID:        permissionID,
// 		Name:      "Updated Permission",
// 		ServiceID: &serviceID,
// 		Action:    ptr.String("action"),
// 	}
// 	updatedPermission, err := resolver.UpdatePermission(context.Background(), input)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, updatedPermission)
// 	assert.Equal(t, input.Name, updatedPermission.Name)

// 	// Test case 2: Missing input
// 	_, err = resolver.UpdatePermission(context.Background(), nil)
// 	assert.Error(t, err)

// 	serviceID = uuid.New()

// 	// Test case 3: Non-existent permission
// 	nonExistentInput := &models.UpdatePermission{
// 		ID:        uuid.New(),
// 		Name:      "Non-existent Permission",
// 		ServiceID: &serviceID,
// 		Action:    ptr.String("action"),
// 	}
// 	_, err = resolver.UpdatePermission(context.Background(), nonExistentInput)
// 	assert.Error(t, err)

// 	// Test case 4: Duplicate name
// 	duplicatePermission := dto.TNTPermission{
// 		PermissionID: uuid.New(),
// 		Name:         "Duplicate Permission",
// 		ServiceID:    uuid.New().String(),
// 		Action:       "action",
// 		RowStatus:    1,
// 		CreatedAt:    time.Now(),
// 		UpdatedAt:    time.Now(),
// 	}
// 	db.Create(&duplicatePermission)

// 	input.Name = duplicatePermission.Name
// 	_, err = resolver.UpdatePermission(context.Background(), input)
// 	assert.Error(t, err)
// }
