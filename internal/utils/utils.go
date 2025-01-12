package utils

import (
	"errors"
	"fmt"
	"go_graphql/config"
	"go_graphql/internal/dto"
	"go_graphql/logger"
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
	var count int64
	if err := config.DB.Model(&dto.TenantResource{}).Where("resource_id = ?", resourceID).Count(&count).Error; err != nil {
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
