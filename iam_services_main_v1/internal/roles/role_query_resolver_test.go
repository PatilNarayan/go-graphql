package roles

import (
	"context"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAllRoles(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	config.DB = db
	resolver := &RoleQueryResolver{DB: db}

	// Setup test data
	tenantID := uuid.New()
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

	// Create tenant resource
	tenantResource := dto.TenantResource{
		ResourceID:     tenantID,
		Name:           "TestTenant",
		ResourceTypeID: mstResType.ResourceTypeID,
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&tenantResource)

	// Create TNT roles
	role1 := dto.TNTRole{
		ResourceID:  uuid.New(),
		Name:        "Role 1",
		Description: "Test Role 1",
		RoleType:    "CUSTOM",
		Version:     "1.0",
		RowStatus:   1,
		CreatedBy:   "admin",
		UpdatedBy:   "admin",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&role1)

	// Create role resource
	roleResource := dto.TenantResource{
		ResourceID:     role1.ResourceID,
		Name:           role1.Name,
		ResourceTypeID: mstResType.ResourceTypeID,
		TenantID:       &tenantID,
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&roleResource)

	// Create MST roles
	mstRole := dto.MstRole{
		RoleID:      uuid.New(),
		Name:        "Master Role",
		Description: "Master Role Description",
		Version:     "1.0",
		RowStatus:   1,
		CreatedBy:   "admin",
		UpdatedBy:   "admin",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&mstRole)

	t.Run("Successfully fetch all roles", func(t *testing.T) {
		roles, err := resolver.AllRoles(ctx, &uuid.Nil)
		assert.NoError(t, err)
		assert.NotNil(t, roles)
		assert.Equal(t, 1, len(roles)) // 1 TNT role + 1 MST role

		// Verify TNT role
		var tntRoleFound bool
		for _, role := range roles {
			if role.ID == role1.ResourceID {
				tntRoleFound = true
				assert.Equal(t, role1.Name, role.Name)
				assert.Equal(t, role1.Description, *role.Description)
				assert.Equal(t, models.RoleTypeEnum(role1.RoleType), role.RoleType)
			}
		}
		assert.True(t, tntRoleFound)
	})
}

func TestGetRole(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	config.DB = db
	resolver := &RoleQueryResolver{DB: db}

	// Setup test data
	tenantID := uuid.New()
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

	// Create tenant resource
	tenantResource := dto.TenantResource{
		ResourceID:     tenantID,
		Name:           "Test Tenant",
		ResourceTypeID: mstResType.ResourceTypeID,
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&tenantResource)

	roleID := uuid.New()
	role := dto.TNTRole{
		ResourceID:  roleID,
		Name:        "Test Role",
		Description: "Test Description",
		RoleType:    "CUSTOM",
		Version:     "1.0",
		RowStatus:   1,
		CreatedBy:   "admin",
		UpdatedBy:   "admin",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&role)

	// Create role resource
	roleResource := dto.TenantResource{
		ResourceID:     roleID,
		Name:           role.Name,
		ResourceTypeID: mstResType.ResourceTypeID,
		TenantID:       &tenantID,
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&roleResource)

	t.Run("Successfully get role by ID", func(t *testing.T) {
		result, err := resolver.GetRole(ctx, roleID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, role.Name, result.Name)
		assert.Equal(t, role.Description, *result.Description)
		assert.Equal(t, models.RoleTypeEnum(role.RoleType), result.RoleType)
	})

	t.Run("Role not found", func(t *testing.T) {
		result, err := resolver.GetRole(ctx, uuid.New())
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "role not found", err.Error())
	})
}

func TestGetGetAllRolesForAssignableScopeRef(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()
	config.DB = db
	resolver := &RoleQueryResolver{DB: db}

	// Setup test data
	tenantID := uuid.New()
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

	mstResTypeTenant := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Tenant",
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&mstResTypeTenant)

	// Create tenant resource
	tenantResource := dto.TenantResource{
		ResourceID:     tenantID,
		Name:           "Test Tenant",
		ResourceTypeID: mstResTypeTenant.ResourceTypeID,
		RowStatus:      1,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(&tenantResource)

	// Create roles for the tenant
	for i := 1; i <= 3; i++ {
		roleID := uuid.New()
		role := dto.TNTRole{
			ResourceID:  roleID,
			Name:        fmt.Sprintf("Role %d", i),
			Description: fmt.Sprintf("Description %d", i),
			RoleType:    "CUSTOM",
			Version:     "1.0",
			RowStatus:   1,
			CreatedBy:   "admin",
			UpdatedBy:   "admin",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		db.Create(&role)

		// Create role resource
		roleResource := dto.TenantResource{
			ResourceID:     roleID,
			Name:           role.Name,
			ResourceTypeID: mstResType.ResourceTypeID,
			TenantID:       &tenantID,
			RowStatus:      1,
			CreatedBy:      "admin",
			UpdatedBy:      "admin",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		db.Create(&roleResource)
	}

	// Create MST role
	mstRole := dto.MstRole{
		RoleID:      uuid.New(),
		Name:        "Master Role",
		Description: "Master Role Description",
		Version:     "1.0",
		RowStatus:   1,
		CreatedBy:   "admin",
		UpdatedBy:   "admin",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&mstRole)

	t.Run("Successfully get all roles for tenant", func(t *testing.T) {
		roles, err := resolver.GetAllRolesForAssignableScopeRef(ctx, tenantID, []uuid.UUID{})
		assert.NoError(t, err)
		assert.NotNil(t, roles)
		assert.Equal(t, 4, len(roles)) // 3 TNT roles + 1 MST role

		// Verify MST role is included
		var mstRoleFound bool
		for _, role := range roles {
			if role.ID == mstRole.RoleID {
				mstRoleFound = true
				assert.Equal(t, mstRole.Name, role.Name)
				assert.Equal(t, mstRole.Description, *role.Description)
				assert.Equal(t, models.RoleTypeEnumDefault, role.RoleType)
			}
		}
		assert.True(t, mstRoleFound)
	})

	t.Run("Invalid tenant ID", func(t *testing.T) {
		roles, err := resolver.GetAllRolesForAssignableScopeRef(ctx, uuid.Nil, []uuid.UUID{})
		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Equal(t, "assignableScopeRef cannot be nil", err.Error())
	})

	t.Run("Tenant not found", func(t *testing.T) {
		roles, err := resolver.GetAllRolesForAssignableScopeRef(ctx, uuid.New(), []uuid.UUID{})
		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Contains(t, err.Error(), "invalid TenantID")
	})
}

func TestGetRolePermissions(t *testing.T) {
	db := setupTestDB()
	config.DB = db

	roleID := uuid.New()
	permID1 := uuid.New()
	permID2 := uuid.New()

	// Create permissions
	perm1 := dto.TNTPermission{
		PermissionID: permID1,
		Name:         "Permission 1",
		ServiceID:    uuid.New().String(),
		Action:       "action1",
		RowStatus:    1,
		CreatedBy:    "admin",
		UpdatedBy:    "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	db.Create(&perm1)

	perm2 := dto.TNTPermission{
		PermissionID: permID2,
		Name:         "Permission 2",
		ServiceID:    uuid.New().String(),
		Action:       "action2",
		RowStatus:    1,
		CreatedBy:    "admin",
		UpdatedBy:    "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	db.Create(&perm2)

	// Create role permissions
	rolePerms := []dto.TNTRolePermission{
		{
			ID:           uuid.New(),
			RoleID:       roleID,
			PermissionID: permID1,
			RowStatus:    1,
			CreatedBy:    "admin",
			UpdatedBy:    "admin",
		},
		{
			ID:           uuid.New(),
			RoleID:       roleID,
			PermissionID: permID2,
			RowStatus:    1,
			CreatedBy:    "admin",
			UpdatedBy:    "admin",
		},
	}
	for _, rp := range rolePerms {
		db.Create(&rp)
	}

	t.Run("Successfully get role permissions", func(t *testing.T) {
		permissions, err := GetRolePermissions(roleID)
		assert.NoError(t, err)
		assert.NotNil(t, permissions)
		assert.Equal(t, 2, len(permissions))

		// Verify permissions
		for _, p := range permissions {
			if p.ID == permID1 {
				assert.Equal(t, perm1.Name, p.Name)
				assert.Equal(t, perm1.Action, *p.Action)
			} else if p.ID == permID2 {
				assert.Equal(t, perm2.Name, p.Name)
				assert.Equal(t, perm2.Action, *p.Action)
			}
		}
	})

	t.Run("Role with no permissions", func(t *testing.T) {
		permissions, err := GetRolePermissions(uuid.New())
		assert.NoError(t, err)
		assert.Nil(t, permissions)
	})
}
