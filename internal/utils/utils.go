package utils

import (
	"errors"
	"fmt"
	"go_graphql/config"
	"go_graphql/internal/dto"
	"go_graphql/logger"
	"regexp"
	"time"

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

	resourceIds, err := GetResourceTypeIDs([]string{"Account", "Client"})
	if err != nil {
		return err
	}

	var count int64
	if err := config.DB.Model(&dto.TenantResource{}).
		Where("resource_id = ? AND row_status = 1 AND resource_type_id IN (?)", resourceID, resourceIds).
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
	if err := config.DB.Model(&dto.TNTRole{}).Where("resource_id = ?", roleId).Count(&count).Error; err != nil {
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
	if err := config.DB.Model(&dto.TNTPermission{}).Where("permission_id = ?", permissionId).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("resource ID does not exist")
	}
	return nil
}

func UpdateDeletedMap() map[string]interface{} {
	return map[string]interface{}{
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		"row_status": 0,
	}
}

func ValidateTenantID(tenantID uuid.UUID) error {
	// Check if the resource ID exists in the database

	resourceType := dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name = ?", "Tenant").First(&resourceType).Error; err != nil {
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
func GetResourceTypeIDs(resourceName []string) ([]string, error) {
	resourceType := []dto.Mst_ResourceTypes{}
	if err := config.DB.Where("name in (?)", resourceName).Find(&resourceType).Error; err != nil {
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
