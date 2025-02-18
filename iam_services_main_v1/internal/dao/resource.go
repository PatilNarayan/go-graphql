package dao

import (
	"fmt"
	"iam_services_main_v1/config"
	"iam_services_main_v1/internal/dto"
)

func GetResourceTypeByName(name string) (*dto.Mst_ResourceTypes, error) {
	DB := config.GetDB()
	resourceType := dto.Mst_ResourceTypes{}
	if err := DB.Where("name = ?", name).First(&resourceType).Error; err != nil {
		return nil, fmt.Errorf("resource type not found: %w", err)
	}
	return &resourceType, nil
}

func GetResourceDetails(conditions map[string]interface{}) (*dto.TenantResource, error) {
	DB := config.GetDB()
	resourceDetails := dto.TenantResource{}
	if err := DB.Where(conditions).First(&resourceDetails).Error; err != nil {
		return nil, fmt.Errorf("parent details not found: %w", err)
	}
	return &resourceDetails, nil
}
