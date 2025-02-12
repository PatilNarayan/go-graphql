package roles

import (
	"context"
	"errors"
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	middleware "iam_services_main_v1/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

// RoleQueryResolver handles role-related queries.
type RoleQueryResolver struct {
	DB *gorm.DB
}

// Roles resolves the list of all roles.
func (r *RoleQueryResolver) AllRoles(ctx context.Context) (models.OperationResult, error) {
	//logger.Log.Info("Fetching all roles")
	// Retrieve x-tenant-id from headers
	ginCtx, ok := ctx.Value(middleware.GinContextKey).(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unable to get gin context")
	}
	tenantID := ginCtx.GetHeader("tenantID")
	if tenantID == "" {
		return nil, errors.New("tenantID not found in headers")
	}

	//validate uuid format
	if _, err := uuid.Parse(tenantID); err != nil {
		return nil, fmt.Errorf("invalid tenantID: %w", err)
	}

	var tntResources []dto.TenantResources
	if err := r.DB.Where("tenant_id = ? AND row_status = 1", tenantID).Find(&tntResources).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch roles from the database")
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}
	var tntResourceIDs []uuid.UUID
	for _, tnt := range tntResources {
		tntResourceIDs = append(tntResourceIDs, *&tnt.ResourceID)
	}

	// if id != nil {
	// 	res, err := r.GetAllRolesForAssignableScopeRef(ctx, *id, tntResourceIDs)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return res, nil
	// }

	var roles []dto.TNTRole
	if err := r.DB.Where("row_status = ? AND resource_id IN (?)", 1, tntResourceIDs).Find(&roles).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch roles from the database")
		return nil, err
	}

	var result []*models.Role
	var data []models.Data
	for _, role := range roles {
		convertedRole := convertRoleToGraphQL(&role)
		result = append(result, convertedRole)
		data = append(data, *convertedRole)
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

	//logger.Log.Infof("Fetched %d roles", len(result))

	return &models.SuccessResponse{
		Success: true,
		Message: "Successfully retrieved tenants",
		Data:    data,
	}, nil
}

// GetRole resolves a single role by ID.
func (r *RoleQueryResolver) GetRole(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	//logger.Log.Infof("Fetching role with ID: %s", id)

	var role dto.TNTRole
	if err := r.DB.First(&role, "resource_id = ? AND row_status = 1", id).Error; err != nil {
		//logger.AddContext(err).Warnf("Role with ID %s not found", id)
		return nil, errors.New("role not found")
	}

	//logger.Log.Infof("Role with ID %s fetched successfully", id)
	data := convertRoleToGraphQL(&role)
	return &models.SuccessResponse{
		Success: true,
		Message: "Successfully retrieved tenants",
		Data:    []models.Data{*data},
	}, nil
}

func (r *RoleQueryResolver) GetAllRolesForAssignableScopeRef(ctx context.Context, assignableScopeRef uuid.UUID, tntResourceIDs []uuid.UUID) (models.OperationResult, error) {
	//logger.Log.Infof("Fetching all roles for tenant with ID: %s", assignableScopeRef)

	if assignableScopeRef == uuid.Nil {
		return nil, fmt.Errorf("assignableScopeRef cannot be nil")
	}

	if err := ValidateMstResType(assignableScopeRef); err != nil {
		return nil, fmt.Errorf("invalid assignableScopeRef: %w", err)
	}

	// resourceIds, err := utils.GetResourceTypeIDs([]string{"Role"})
	// if err != nil {
	// 	return nil, err
	// }

	var roles []dto.TNTRole
	if err := r.DB.Where("scope_resource_type_id = ? AND row_status = 1 AND resource_id IN (?)", assignableScopeRef, tntResourceIDs).Find(&roles).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch roles from the database")
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	var result []*models.Role
	var data []models.Data
	for _, role := range roles {
		convertedRole := convertRoleToGraphQL(&role)
		result = append(result, convertedRole)
		data = append(data, *convertedRole)
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

	//logger.Log.Infof("Fetched %d roles for tenant with ID: %s", len(result), assignableScopeRef)

	return &models.SuccessResponse{
		Success: true,
		Message: "Successfully retrieved tenants",
		Data:    data,
	}, nil
}

// Helper function to convert database Role to GraphQL Role models.
func convertRoleToGraphQL(role *dto.TNTRole) *models.Role {
	//logger.Log.Infof("Converting role to GraphQL model for Role ID: %s", role.ResourceID)

	res := &models.Role{
		ID:          role.ResourceID,
		Name:        role.Name,
		Description: ptr.String(role.Description),
		RoleType:    models.RoleTypeEnum(role.RoleType),
		Version:     role.Version,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
		UpdatedBy:   role.UpdatedBy,
		CreatedBy:   role.CreatedBy,
	}

	if res.RoleType == "D" {
		res.RoleType = models.RoleTypeEnumDefault
	} else if res.RoleType == "C" {
		res.RoleType = models.RoleTypeEnumCustom
	}

	permissions, err := GetRolePermissions(role.ResourceID)
	if err != nil {
		//logger.AddContext(err).Error("Failed to fetch role permissions")
		return nil
	}

	// mstPermissions, err := GetMSTRolePermission()
	// if err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch role permissions")
	// 	return nil
	// }
	// permissions = append(permissions, mstPermissions...)
	res.Permissions = permissions

	var ParentResource dto.Mst_ResourceTypes
	if err := config.DB.Where(&dto.Mst_ResourceTypes{ResourceTypeID: role.ScopeResourceTypeID, RowStatus: 1}).First(&ParentResource).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch parent resource")
		return nil
	}

	res.AssignableScope = &models.Root{
		ID:        ParentResource.ResourceTypeID,
		Name:      ParentResource.Name,
		CreatedAt: ParentResource.CreatedAt.String(),
		UpdatedAt: ParentResource.UpdatedAt.String(),
		CreatedBy: ParentResource.CreatedBy,
		UpdatedBy: ParentResource.UpdatedBy,
	}

	//logger.Log.Infof("Successfully converted Role ID: %s to GraphQL model", role.ResourceID)
	return res
}

func GetRolePermissions(id uuid.UUID) ([]*models.Permission, error) {
	//logger.Log.Infof("Fetching role permissions for role ID: %s", id)

	var rolePermissions []dto.TNTRolePermission
	if err := config.DB.Where("role_id = ? AND row_status = 1", id).Find(&rolePermissions).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch role permissions from the database")
		return nil, err
	}

	if len(rolePermissions) == 0 {
		//logger.Log.Infof("No role permissions found for role ID: %s", id)
		return nil, nil
	}
	permissionIDs := make([]uuid.UUID, len(rolePermissions))
	for _, rolePermission := range rolePermissions {
		permissionIDs = append(permissionIDs, rolePermission.PermissionID)
	}

	var permissions []dto.TNTPermission
	if err := config.DB.Where("permission_id in (?) AND row_status = 1", permissionIDs).Find(&permissions).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch permissions from the database")
		return nil, err
	}

	var result []*models.Permission
	for _, rolePermission := range permissions {
		result = append(result, &models.Permission{
			ID:        rolePermission.PermissionID,
			Name:      rolePermission.Name,
			ServiceID: &rolePermission.ServiceID,
			Action:    &rolePermission.Action,
			CreatedAt: rolePermission.CreatedAt.String(),
			UpdatedAt: rolePermission.UpdatedAt.String(),
			UpdatedBy: rolePermission.UpdatedBy,
			CreatedBy: rolePermission.CreatedBy,
		})
	}

	//logger.Log.Infof("Fetched %d role permissions for role ID: %s", len(result), id)
	return result, nil
}

func GetMSTRolePermission(id uuid.UUID) ([]*models.Permission, error) {
	//logger.Log.Infof("Fetching All role permissions")

	var rolePermissions []dto.MstRolePermission
	if err := config.DB.Where("role_id = ? AND row_status = 1", id).Find(&rolePermissions).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch role permissions from the database")
		return nil, err
	}

	if len(rolePermissions) == 0 {
		//logger.Log.Infof("No role permissions found")
		return nil, nil
	}
	permissionIDs := make([]uuid.UUID, len(rolePermissions))
	for _, rolePermission := range rolePermissions {
		permissionIDs = append(permissionIDs, rolePermission.PermissionID)
	}

	var permissions []dto.MstPermission
	if err := config.DB.Where("permission_id in (?) AND row_status = 1", permissionIDs).Find(&permissions).Error; err != nil {
		//logger.AddContext(err).Error("Failed to fetch permissions from the database")
		return nil, err
	}

	var result []*models.Permission
	for _, rolePermission := range permissions {
		result = append(result, &models.Permission{
			ID:        rolePermission.PermissionID,
			Name:      rolePermission.Name,
			ServiceID: &rolePermission.ServiceID,
			Action:    &rolePermission.Action,
			CreatedAt: rolePermission.CreatedAt.String(),
			UpdatedAt: rolePermission.UpdatedAt.String(),
			UpdatedBy: uuid.MustParse(rolePermission.UpdatedBy),
			CreatedBy: uuid.MustParse(rolePermission.CreatedBy),
		})
	}

	//logger.Log.Infof("Fetched %d role permissions", len(result))
	return result, nil
}

func convertMSTRoleToGraphQL(role *dto.MstRole) *models.Role {
	//logger.Log.Infof("Converting role to GraphQL model for Role ID: %s", role.RoleID)

	res := &models.Role{
		ID:          role.RoleID,
		Name:        role.Name,
		Description: ptr.String(role.Description),
		RoleType:    models.RoleTypeEnumDefault,
		Version:     role.Version,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
		UpdatedBy:   role.UpdatedBy,
		CreatedBy:   role.CreatedBy,
	}

	permissions, err := GetMSTRolePermission(role.RoleID)
	if err != nil {
		//logger.AddContext(err).Error("Failed to fetch role permissions")
		return nil
	}

	// mstPermissions, err := GetMSTRolePermission()
	// if err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch role permissions")
	// 	return nil
	// }
	// permissions = append(permissions, mstPermissions...)
	res.Permissions = permissions

	// var childResource, ParentResource dto.TenantResources
	// if err := config.DB.Where(&dto.TenantResources{ResourceID: role.RoleID}).First(&childResource).Error; err != nil {
	// 	logger.AddContext(err).Error("Failed to fetch parent resource")
	// 	return nil
	// }
	// if err := config.DB.Where(&dto.TenantResources{ResourceID: *childResource.TenantID}).First(&ParentResource).Error; err != nil {
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

	//logger.Log.Infof("Successfully converted Role ID: %s to GraphQL model", role.RoleID)
	return res
}
