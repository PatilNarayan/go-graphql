package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go_graphql/config"
	"go_graphql/gql/models"
	"go_graphql/internal/constants"
	"go_graphql/internal/dto"
	"go_graphql/internal/utils"
	"go_graphql/logger"
	"go_graphql/permit"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleMutationResolver handles role-related mutations.
type RoleMutationResolver struct {
	DB *gorm.DB
}

// CreateRole handles creating a new role.
func (r *RoleMutationResolver) CreateRole(ctx context.Context, input models.CreateRoleInput) (*models.Role, error) {
	logger.Log.Info("Starting CreateRole")

	// Validate input
	if input.Name == "" {
		logger.Log.Warn("Role name is required")
		return nil, errors.New("role name is required")
	} else {
		if err := utils.ValidateName(input.Name); err != nil {
			logger.AddContext(err).Error("Invalid role name")
			return nil, fmt.Errorf("invalid role name: %w", err)
		}
	}

	if input.AssignableScopeRef == uuid.Nil {
		logger.Log.Warn("Account ID / Client ID is required")
		return nil, errors.New("account ID / client ID is required")
	} else {
		if err := utils.ValidateResourceID(input.AssignableScopeRef); err != nil {
			logger.AddContext(err).Error("invalid account ID / client ID")
			return nil, fmt.Errorf("invalid account ID / client ID: %w", err)
		}
	}

	// check if role already exists
	var roleExists dto.TNTRole
	if err := r.DB.Where("name = ? AND row_status = 1 AND role_type = ? AND resource_id = ? AND deleted_at IS NULL", input.Name, dto.RoleTypeEnum(input.RoleType), input.AssignableScopeRef).First(&roleExists).Error; err == nil {
		logger.AddContext(err).Error("Role already exists")
		return nil, fmt.Errorf("role already exists")
	}

	if err := CheckPermissions(input.Permissions); err != nil {
		return nil, err
	}

	var assignableScopeRef dto.TenantResource
	if err := r.DB.Where("resource_id = ? AND row_status = 1", input.AssignableScopeRef).First(&assignableScopeRef).Error; err != nil {
		logger.AddContext(err).Error("Invalid TenantID")
		return nil, fmt.Errorf("invalid TenantID")
	}

	//create role in permit
	pc := permit.NewPermitClient()
	_, err := pc.APIExecute(ctx, "POST", fmt.Sprintf("resources/%s/roles", assignableScopeRef.Name), map[string]interface{}{
		"name": input.Name,
		"key":  input.Name,
	})

	if err != nil {
		return nil, err
	}

	role := dto.TNTRole{}
	role.Name = input.Name
	if input.Description != nil {
		role.Description = *input.Description
	}

	role.RoleType = dto.RoleTypeEnum(input.RoleType)
	role.Version = input.Version
	role.CreatedAt = time.Now()
	role.CreatedBy = constants.DefaltCreatedBy
	role.UpdatedBy = constants.DefaltCreatedBy
	role.RowStatus = 1

	logger.Log.Info("Fetching resource type")
	resourceType := dto.Mst_ResourceTypes{}
	if err := r.DB.Where("name = ?", "Role").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return nil, fmt.Errorf("resource type not found: %w", err)
	}

	logger.Log.Info("Creating tenant resource")
	tenantResource := &dto.TenantResource{
		ResourceID:       uuid.New(),
		ParentResourceID: &input.AssignableScopeRef,
		ResourceTypeID:   resourceType.ResourceTypeID,
		Name:             input.Name,
		RowStatus:        1,
		CreatedBy:        constants.DefaltCreatedBy,
		UpdatedBy:        constants.DefaltCreatedBy,
		CreatedAt:        time.Now(),
		TenantID:         &input.AssignableScopeRef,
	}

	if err := r.DB.Create(&tenantResource).Error; err != nil {
		logger.AddContext(err).Error("Failed to create tenant resource")
		return nil, err
	}

	role.ResourceID = tenantResource.ResourceID
	logger.Log.Info("Saving role to database")
	if err := r.DB.Create(&role).Error; err != nil {
		logger.AddContext(err).Error("Failed to save role")
		return nil, err
	}

	for _, permissionID := range input.Permissions {
		tx := r.DB.Create(&dto.TNTRolePermission{
			ID:           uuid.New(),
			RoleID:       role.ResourceID,
			PermissionID: uuid.MustParse(permissionID),
			RowStatus:    1,
			CreatedBy:    constants.DefaltCreatedBy,
			UpdatedBy:    constants.DefaltCreatedBy,
		})

		if tx.Error != nil {
			logger.AddContext(tx.Error).Error("Failed to save role permission")
			return nil, tx.Error
		}
	}
	logger.Log.Info("Role created successfully")
	return convertRoleToGraphQL(&role), nil
}

// UpdateRole handles updating an existing role.
func (r *RoleMutationResolver) UpdateRole(ctx context.Context, input models.UpdateRoleInput) (*models.Role, error) {
	logger.Log.Infof("Starting UpdateRole for ID: %s", input.ID)

	var role dto.TNTRole
	if err := r.DB.First(&role, "resource_id = ?", input.ID).Error; err != nil {
		logger.AddContext(err).Warn("Role not found")
		return nil, errors.New("role not found")
	}

	// Update fields
	if input.Name != "" {
		role.Name = input.Name
	}

	if input.AssignableScopeRef == uuid.Nil {
		logger.Log.Warn("Tenant ID is required")
		return nil, errors.New("Tenant ID is required")
	} else {
		if err := utils.ValidateResourceID(input.AssignableScopeRef); err != nil {
			logger.AddContext(err).Error("Invalid TenantID")
			return nil, fmt.Errorf("invalid TenantID")
		}
	}

	if err := CheckPermissions(input.Permissions); err != nil {
		return nil, err
	}

	// check if role already exists
	var roleExists dto.TNTRole
	if input.Name != role.Name && input.Name != "" {
		if err := r.DB.Where("name = ? AND row_status = 1 AND role_type = ? AND resource_id = ? AND deleted_at IS NULL", input.Name, dto.RoleTypeEnum(input.RoleType), input.AssignableScopeRef).First(&roleExists).Error; err == nil {
			logger.AddContext(err).Error("Role already exists")
			return nil, fmt.Errorf("role already exists")
		}
	}

	logger.Log.Infof("Updating role in database for ID: %s", input.ID)
	updateData := map[string]interface{}{
		"version":    input.Version,
		"role_type":  dto.RoleTypeEnum(input.RoleType),
		"updated_by": constants.DefaltCreatedBy,
		"updated_at": time.Now(),
	}
	if input.Name != "" {
		updateData["name"] = input.Name
	}
	if input.Description != nil {
		updateData["description"] = *input.Description
	}

	if err := r.DB.Model(&role).Updates(updateData).Error; err != nil {
		logger.AddContext(err).Error("Failed to update role")
		return nil, err
	}

	var pdata []dto.TNTRolePermission
	if err := r.DB.Where("role_id = ?", input.ID).Find(&pdata).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch role permissions")
		return nil, err
	}
	for _, pid := range input.Permissions {
		exists := false
		for _, p := range pdata {
			if p.PermissionID.String() == pid {
				exists = true
				break
			}
		}
		if !exists {
			tx := r.DB.Create(&dto.TNTRolePermission{
				ID:           uuid.New(),
				RoleID:       input.ID,
				PermissionID: uuid.MustParse(pid),
				RowStatus:    1,
				CreatedBy:    constants.DefaltCreatedBy,
				UpdatedBy:    constants.DefaltCreatedBy,
			})

			if tx.Error != nil {
				logger.AddContext(tx.Error).Error("Failed to save role permission")
				return nil, tx.Error
			}

		}
	}
	removeIDs := make([]uuid.UUID, 0)
	for _, p := range pdata {
		exists := false
		for _, pid := range input.Permissions {
			if p.PermissionID.String() == pid {
				exists = true
				break
			}
		}
		if !exists {
			removeIDs = append(removeIDs, p.ID)
		}
	}

	for _, id := range removeIDs {
		if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_permission_id = ?", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
			logger.AddContext(err).Error("Failed to delete role")
			return nil, fmt.Errorf("failed to delete role: %w", err)
		}
	}

	logger.Log.Infof("Role updated successfully for ID: %s", input.ID)
	return convertRoleToGraphQL(&role), nil
}

// DeleteRole handles deleting a role by ID.
func (r *RoleMutationResolver) DeleteRole(ctx context.Context, id uuid.UUID) (bool, error) {
	logger.Log.Infof("Starting DeleteRole for ID: %s", id)

	var roleDB dto.TNTRole
	if err := r.DB.First(&roleDB, "resource_id = ?", id).Error; err != nil {
		logger.AddContext(err).Warn("Role not found")
		return false, errors.New("role not found")
	}

	updates := map[string]interface{}{
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}

	logger.Log.Infof("Marking role as deleted for ID: %s", id)
	if err := r.DB.Model(&dto.TNTRole{}).Where("resource_id = ?", id).Updates(updates).Error; err != nil {
		logger.AddContext(err).Error("Failed to delete role")
		return false, err
	}

	if err := r.DB.Model(&dto.TNTRolePermission{}).Where("role_id = ?", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		logger.AddContext(err).Error("Failed to delete role")
		return false, fmt.Errorf("failed to delete role: %w", err)
	}

	if err := r.DB.Model(&dto.TenantResource{}).Where("resource_id = ?", id).Updates(utils.UpdateDeletedMap()).Error; err != nil {
		logger.AddContext(err).Error("Failed to delete role")
		return false, fmt.Errorf("failed to delete role: %w", err)
	}

	logger.Log.Infof("Role deleted successfully for ID: %s", id)
	return true, nil
}

func CheckPermissions(permissionsIds []string) error {
	//validate permissionIds
	for _, permissionID := range permissionsIds {
		if err := utils.ValidatePermissionID(permissionID); err != nil {
			logger.AddContext(err).Error("Invalid permission ID")
			return fmt.Errorf("invalid permission ID: %w", err)
		}
	}
	return nil
}

func CreateMstRole(tenantID uuid.UUID) error {
	var mstRole []*dto.MstRole

	if err := config.DB.Find(&mstRole).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch roles from the database")
		return err
	}

	var mstRolePermissions []*dto.MstRolePermission
	if err := config.DB.Find(&mstRolePermissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch role permissions from the database")
		return err
	}
	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ?", "Role").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return fmt.Errorf("resource type not found: %w", err)
	}

	ResDefaultPermissions, err := AddDefaultPermissions()
	if err != nil {
		return err
	}

	for _, mrole := range mstRole {
		tenantResource := &dto.TenantResource{
			ResourceID:       uuid.New(),
			ParentResourceID: &tenantID,
			ResourceTypeID:   resourceType.ResourceTypeID,
			Name:             mrole.Name,
			RowStatus:        1,
			CreatedBy:        constants.DefaltCreatedBy,
			UpdatedBy:        constants.DefaltCreatedBy,
			CreatedAt:        time.Now(),
			TenantID:         &tenantID,
		}

		if err := config.DB.Save(tenantResource).Error; err != nil {
			logger.AddContext(err).Error("Failed to save role")
			return err
		}

		role := dto.TNTRole{
			ResourceID:     tenantResource.ResourceID,
			Name:           mrole.Name,
			Description:    mrole.Description,
			RoleType:       dto.RoleTypeEnumDefault,
			ResourceTypeID: resourceType.ResourceTypeID,
			Version:        mrole.Version,
			CreatedAt:      mrole.CreatedAt,
			CreatedBy:      mrole.CreatedBy,
			UpdatedBy:      mrole.UpdatedBy,
			UpdatedAt:      mrole.UpdatedAt,
			RowStatus:      mrole.RowStatus,
		}

		if err := config.DB.Save(&role).Error; err != nil {
			logger.AddContext(err).Error("Failed to save role")
			return err
		}

		for _, permissionID := range mstRolePermissions {
			if permissionID.RoleID != mrole.RoleID {
				continue
			}
			tx := config.DB.Create(&dto.TNTRolePermission{
				ID:           uuid.New(),
				RoleID:       role.ResourceID,
				PermissionID: ResDefaultPermissions[permissionID.PermissionID],
				RowStatus:    1,
				CreatedBy:    constants.DefaltCreatedBy,
				UpdatedBy:    constants.DefaltCreatedBy,
			})

			if tx.Error != nil {
				logger.AddContext(tx.Error).Error("Failed to save role permission")
				return tx.Error
			}
		}
	}

	return nil

}

func AddDefaultPermissions() (map[uuid.UUID]uuid.UUID, error) {
	var mstPermissions []*dto.MstPermission

	if err := config.DB.Find(&mstPermissions).Error; err != nil {
		logger.AddContext(err).Error("Failed to fetch permissions from the database")
		return nil, err
	}

	res := make(map[uuid.UUID]uuid.UUID)
	for _, mpermission := range mstPermissions {
		permission := &dto.TNTPermission{
			PermissionID: uuid.New(),
			Name:         mpermission.Name,
			ServiceID:    mpermission.ServiceID,
			Action:       mpermission.Action,
			RowStatus:    1,
			// RoleID:       *input.RoleID,
			CreatedBy: constants.DefaltCreatedBy,
			UpdatedBy: constants.DefaltUpdatedBy,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		res[mpermission.PermissionID] = permission.PermissionID
		if err := config.DB.Create(permission).Error; err != nil {
			logger.AddContext(err).Error("Failed to create permission")
			return nil, err
		}
	}

	return res, nil
}
