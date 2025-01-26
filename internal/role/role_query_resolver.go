package role

import (
	"context"
	"errors"
	"fmt"

	"go_graphql/config"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"go_graphql/internal/utils"
	"go_graphql/logger"

	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

// RoleQueryResolver handles role-related queries.
type RoleQueryResolver struct {
	DB *gorm.DB
}

// Roles resolves the list of all roles.
func (r *RoleQueryResolver) AllRoles(ctx context.Context, id uuid.UUID) ([]*models.Role, error) {
	logger.Log.Info("Fetching all roles")

	if id != uuid.Nil {
		res, err := r.GetAllRolesForAssignableScopeRef(ctx, id)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	var roles []dto.TNTRole
	if err := r.DB.Where("row_status = ?", 1).Find(&roles).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch roles from the database")
		return nil, err
	}

	var result []*models.Role
	for _, role := range roles {
		convertedRole := convertRoleToGraphQL(&role)
		result = append(result, convertedRole)
	}

	// var mstroles []dto.MstRole
	// if err := r.DB.Find(&mstroles).Error; err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch roles from the database")
	// 	return nil, err
	// }

	// for _, role := range mstroles {
	// 	convertedRole := convertMSTRoleToGraphQL(&role)
	// 	result = append(result, convertedRole)
	// }

	logger.Log.Infof("Fetched %d roles", len(result))
	return result, nil
}

// GetRole resolves a single role by ID.
func (r *RoleQueryResolver) GetRole(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	logger.Log.Infof("Fetching role with ID: %s", id)

	var role dto.TNTRole
	if err := r.DB.First(&role, "resource_id = ? AND row_status = 1", id).Error; err != nil {
		logger.AddContext(err).Warnf("Role with ID %s not found", id)
		return nil, errors.New("role not found")
	}

	logger.Log.Infof("Role with ID %s fetched successfully", id)
	return convertRoleToGraphQL(&role), nil
}

func (r *RoleQueryResolver) GetAllRolesForAssignableScopeRef(ctx context.Context, assignableScopeRef uuid.UUID) ([]*models.Role, error) {
	logger.Log.Infof("Fetching all roles for tenant with ID: %s", assignableScopeRef)

	if assignableScopeRef == uuid.Nil {
		return nil, fmt.Errorf("assignableScopeRef cannot be nil")
	}

	if err := utils.ValidateResourceID(assignableScopeRef); err != nil {
		return nil, fmt.Errorf("invalid assignableScopeRef: %w", err)
	}

	resourceIds, err := utils.GetResourceTypeIDs([]string{"Role"})
	if err != nil {
		return nil, err
	}

	var assignableResource []dto.TenantResource
	if err := r.DB.Where("parent_resource_id = ? AND row_status = 1 AND resource_type_id IN (?)", assignableScopeRef, resourceIds).Find(&assignableResource).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tenant resource: %w", err)
	}

	if len(assignableResource) == 0 {
		return nil, fmt.Errorf("assignableScopeRef %s does not exist", assignableScopeRef)
	}

	roleIds := make([]uuid.UUID, len(assignableResource))
	for i, resource := range assignableResource {
		roleIds[i] = resource.ResourceID
	}

	var roles []dto.TNTRole
	if err := r.DB.Where("resource_id IN (?) AND row_status = 1", roleIds).Find(&roles).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch roles from the database")
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	var result []*models.Role
	for _, role := range roles {
		convertedRole := convertRoleToGraphQL(&role)
		result = append(result, convertedRole)
	}

	// var mstroles []dto.MstRole
	// if err := r.DB.Find(&mstroles).Error; err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch roles from the database")
	// 	return nil, err
	// }

	// for _, role := range mstroles {
	// 	convertedRole := convertMSTRoleToGraphQL(&role)
	// 	result = append(result, convertedRole)
	// }

	logger.Log.Infof("Fetched %d roles for tenant with ID: %s", len(result), assignableScopeRef)
	return result, nil
}

// Helper function to convert database Role to GraphQL Role models.
func convertRoleToGraphQL(role *dto.TNTRole) *models.Role {
	logger.Log.Infof("Converting role to GraphQL model for Role ID: %s", role.ResourceID)

	res := &models.Role{
		ID:          role.ResourceID,
		Name:        role.Name,
		Description: ptr.String(role.Description),
		RoleType:    models.RoleTypeEnum(role.RoleType),
		Version:     role.Version,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   ptr.String(role.UpdatedAt.String()),
		UpdatedBy:   &role.UpdatedBy,
		CreatedBy:   &role.CreatedBy,
	}

	permissions, err := GetRolePermissions(role.ResourceID)
	if err != nil {
		logger.AddContext(err).Error("Failed to fetch role permissions")
		return nil
	}

	// mstPermissions, err := GetMSTRolePermission()
	// if err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch role permissions")
	// 	return nil
	// }
	// permissions = append(permissions, mstPermissions...)
	res.Permissions = permissions

	var childResource dto.TenantResource
	if err := config.DB.Where(&dto.TenantResource{ResourceID: role.ResourceID, RowStatus: 1}).First(&childResource).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch parent resource")
		return nil
	}
	var ParentResource dto.Mst_ResourceTypes
	if err := config.DB.Where(&dto.Mst_ResourceTypes{ResourceTypeID: *childResource.ParentResourceID, RowStatus: 1}).First(&ParentResource).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch parent resource")
		return nil
	}

	res.AssignableScope = &models.Root{
		ID:        ParentResource.ResourceTypeID,
		Name:      ParentResource.Name,
		CreatedAt: ParentResource.CreatedAt.String(),
		UpdatedAt: ptr.String(ParentResource.UpdatedAt.String()),
		CreatedBy: &ParentResource.CreatedBy,
		UpdatedBy: &ParentResource.UpdatedBy,
	}

	logger.Log.Infof("Successfully converted Role ID: %s to GraphQL model", role.ResourceID)
	return res
}

func GetRolePermissions(id uuid.UUID) ([]*models.Permission, error) {
	logger.Log.Infof("Fetching role permissions for role ID: %s", id)

	var rolePermissions []dto.TNTRolePermission
	if err := config.DB.Where("role_id = ? AND row_status = 1", id).Find(&rolePermissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch role permissions from the database")
		return nil, err
	}

	if len(rolePermissions) == 0 {
		logger.Log.Infof("No role permissions found for role ID: %s", id)
		return nil, nil
	}
	permissionIDs := make([]uuid.UUID, len(rolePermissions))
	for _, rolePermission := range rolePermissions {
		permissionIDs = append(permissionIDs, rolePermission.PermissionID)
	}

	var permissions []dto.TNTPermission
	if err := config.DB.Where("permission_id in (?) AND row_status = 1", permissionIDs).Find(&permissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch permissions from the database")
		return nil, err
	}

	var result []*models.Permission
	for _, rolePermission := range permissions {
		result = append(result, &models.Permission{
			ID:        rolePermission.PermissionID,
			Name:      rolePermission.Name,
			ServiceID: &rolePermission.ServiceID,
			Action:    &rolePermission.Action,
			CreatedAt: ptr.String(rolePermission.CreatedAt.String()),
			UpdatedAt: ptr.String(rolePermission.UpdatedAt.String()),
			UpdatedBy: &rolePermission.UpdatedBy,
			CreatedBy: rolePermission.CreatedBy,
		})
	}

	logger.Log.Infof("Fetched %d role permissions for role ID: %s", len(result), id)
	return result, nil
}

func GetMSTRolePermission(id uuid.UUID) ([]*models.Permission, error) {
	logger.Log.Infof("Fetching All role permissions")

	var rolePermissions []dto.MstRolePermission
	if err := config.DB.Where("role_id = ? AND row_status = 1", id).Find(&rolePermissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch role permissions from the database")
		return nil, err
	}

	if len(rolePermissions) == 0 {
		logger.Log.Infof("No role permissions found")
		return nil, nil
	}
	permissionIDs := make([]uuid.UUID, len(rolePermissions))
	for _, rolePermission := range rolePermissions {
		permissionIDs = append(permissionIDs, rolePermission.PermissionID)
	}

	var permissions []dto.MstPermission
	if err := config.DB.Where("permission_id in (?) AND row_status = 1", permissionIDs).Find(&permissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch permissions from the database")
		return nil, err
	}

	var result []*models.Permission
	for _, rolePermission := range permissions {
		result = append(result, &models.Permission{
			ID:        rolePermission.PermissionID,
			Name:      rolePermission.Name,
			ServiceID: &rolePermission.ServiceID,
			Action:    &rolePermission.Action,
			CreatedAt: ptr.String(rolePermission.CreatedAt.String()),
			UpdatedAt: ptr.String(rolePermission.UpdatedAt.String()),
			UpdatedBy: &rolePermission.UpdatedBy,
			CreatedBy: rolePermission.CreatedBy,
		})
	}

	logger.Log.Infof("Fetched %d role permissions", len(result))
	return result, nil
}

func convertMSTRoleToGraphQL(role *dto.MstRole) *models.Role {
	logger.Log.Infof("Converting role to GraphQL model for Role ID: %s", role.RoleID)

	res := &models.Role{
		ID:          role.RoleID,
		Name:        role.Name,
		Description: ptr.String(role.Description),
		RoleType:    models.RoleTypeEnumDefault,
		Version:     role.Version,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   ptr.String(role.UpdatedAt.String()),
		UpdatedBy:   &role.UpdatedBy,
		CreatedBy:   &role.CreatedBy,
	}

	permissions, err := GetMSTRolePermission(role.RoleID)
	if err != nil {
		logger.AddContext(err).Error("Failed to fetch role permissions")
		return nil
	}

	// mstPermissions, err := GetMSTRolePermission()
	// if err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch role permissions")
	// 	return nil
	// }
	// permissions = append(permissions, mstPermissions...)
	res.Permissions = permissions

	// var childResource, ParentResource dto.TenantResource
	// if err := config.DB.Where(&dto.TenantResource{ResourceID: role.RoleID}).First(&childResource).Error; err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch parent resource")
	// 	return nil
	// }
	// if err := config.DB.Where(&dto.TenantResource{ResourceID: *childResource.TenantID}).First(&ParentResource).Error; err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch parent resource")
	// 	return nil
	// }

	// res.AssignableScope = &models.Root{
	// 	ID:        ParentResource.ResourceID,
	// 	Name:      ParentResource.Name,
	// 	CreatedAt: ParentResource.CreatedAt.String(),
	// 	UpdatedAt: ptr.String(ParentResource.UpdatedAt.String()),
	// 	CreatedBy: &ParentResource.CreatedBy,
	// 	UpdatedBy: &ParentResource.UpdatedBy,
	// }

	logger.Log.Infof("Successfully converted Role ID: %s to GraphQL model", role.RoleID)
	return res
}
