package utils

import (
	"errors"
	"go_graphql/config"
	"go_graphql/internal/dto"

	"github.com/google/uuid"
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
