package utils

import (
	"errors"
	"fmt"
	"go_graphql/config"
	"go_graphql/internal/dto"
	"go_graphql/logger"
	"regexp"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func StringValue(s *string) string {
	if s != nil {
		return *s
	}
	return "" // or return your preferred default value
}

func ValidateResourceID(resourceID uuid.UUID) error {
	// Check if the resource ID exists in the database

	// resourceIds, err := GetResourceTypeIDs([]string{"Account", "Client"})
	// if err != nil {
	// 	return err
	// }

	var count int64
	if err := config.DB.Model(&dto.TenantResource{}).
		// Where("resource_id = ? AND row_status = 1 AND resource_type_id IN (?)", resourceID, resourceIds).
		Where("resource_id = ? AND row_status = 1", resourceID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func ValidateRoleID(roleId uuid.UUID) error {
	// Check if the resource ID exists in the database
	var count int64
	if err := config.DB.Model(&dto.TNTRole{}).Where("resource_id = ? AND row_status = 1", roleId).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func ValidatePermissionID(permissionId string) error {
	// Check if the resource ID exists in the database
	var count int64
	if err := config.DB.Model(&dto.TNTPermission{}).Where("permission_id = ? AND row_status = 1", permissionId).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func UpdateDeletedMap() map[string]interface{} {
	return map[string]interface{}{
		"row_status": 0,
	}
}

func ValidateTenantID(tenantID uuid.UUID) error {
	// Check if the resource ID exists in the database

	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ? AND row_status = 1", "Tenant").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return fmt.Errorf("resource type not found: %w", err)
	}
	var count int64
	if err := config.DB.Model(&dto.TenantResource{}).Where("resource_id = ? AND row_status = 1 AND resource_type_id = ?", tenantID, resourceType.ResourceTypeID).Count(&count).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func ValidateMstResType(resourceID uuid.UUID) error {
	var count int64
	if err := config.DB.Model(&dto.Mst_ResourceTypes{}).
		Where("resource_type_id = ? AND row_status = 1", resourceID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func ValidateRole(resourceID uuid.UUID) error {
	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ? AND row_status = 1", "Role").First(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return fmt.Errorf("resource type not found: %w", err)
	}
	var count int64
	if err := config.DB.Model(&dto.TenantResource{}).
		// Where("resource_id = ? AND row_status = 1 AND resource_type_id IN (?)", resourceID, resourceIds).
		Where("resource_id = ? AND row_status = 1 AND resource_type_id = ?", resourceID, resourceType.ResourceTypeID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func GetResourceTypeIDName(resourceID uuid.UUID) (*string, error) {
	var data dto.Mst_ResourceTypes
	if err := config.DB.Model(&dto.Mst_ResourceTypes{}).Where("resource_type_id = ? AND row_status = 1", resourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data.Name, nil
}

func GetResourceTypeIDs(resourceName []string) ([]string, error) {
	resourceType := []dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name in (?) AND row_status = 1", resourceName).Find(&resourceType).Error; err != nil {
		logger.AddContext(err).Error("Resource type not found")
		return nil, err
	}
	var resourceIds []string
	for _, resource := range resourceType {
		resourceIds = append(resourceIds, resource.ResourceTypeID.String())
	}

	return resourceIds, nil
}

// ValidateName validates that the input string matches the regex "^[A-Za-z0-9\\-_]+$".
func ValidateName(name string) error {
	// Define the regex pattern
	pattern := `^[A-Za-z0-9\-_]+$`
	// Compile the regex
	re := regexp.MustCompile(pattern)
	// Check if the name matches the regex
	if !re.MatchString(name) {
		return errors.New("invalid name: must contain only alphanumeric characters, hyphens, or underscores")
	}
	return nil
}

func GetPermissionAction(permissionsIds []string) ([]string, error) {
	var actions []string
	//validate permissionIds
	for _, permissionID := range permissionsIds {
		var data dto.TNTPermission
		if err := config.DB.Model(&dto.TNTPermission{}).Where("permission_id = ? AND row_status = 1", permissionID).First(&data).Error; err != nil {
			return nil, err
		}
		actions = append(actions, data.Action)
	}
	return actions, nil
}

func CreateActionMap(store map[string]interface{}, actions []string) map[string]interface{} {
	for _, action := range actions {
		store[action] = map[string]interface{}{
			"name": action,
		}
	}
	return store
}

func GetActionMap(data []interface{}, key string) map[string]interface{} {
	actionMap := make(map[string]interface{})
	for _, d := range data {
		d := d.(map[string]interface{})
		if d["key"].(string) == key {
			actionMap = d["actions"].(map[string]interface{})
		}
	}
	for _, value := range actionMap {
		value := value.(map[string]interface{})
		for key1 := range value {
			if key1 != "name" {
				delete(value, key1)
			}
		}
	}
	return actionMap
}

func GetPermissionResourceAction(resourceID uuid.UUID) ([]string, error) {
	var resources []dto.TenantResource
	if err := config.DB.Model(&dto.TenantResource{}).Where("parent_resource_id = ? AND row_status = 1", resourceID).Find(&resources).Error; err != nil {
		return nil, err
	}

	resourceIds := []string{}
	for _, resource := range resources {
		resourceIds = append(resourceIds, resource.ResourceID.String())
	}

	var rolePermissions []dto.TNTRolePermission
	if err := config.DB.Model(&dto.TNTRolePermission{}).Where("role_id in (?) AND row_status = 1", resourceIds).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	permissionIds := []string{}
	for _, rolePermission := range rolePermissions {
		permissionIds = append(permissionIds, rolePermission.PermissionID.String())
	}

	actions, err := GetPermissionAction(permissionIds)
	if err != nil {
		return nil, err
	}
	return actions, nil
}
