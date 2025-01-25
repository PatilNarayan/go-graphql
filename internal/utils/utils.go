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

func ValidateMstResType(resourceID uuid.UUID) error {
	var count int64
	if err := config.DB.Model(&dto.Mst_ResourceTypes{}).
		// Where("resource_id = ? AND row_status = 1 AND resource_type_id IN (?)", resourceID, resourceIds).
		Where("resource_type_id = ? AND row_status = 1", resourceID).
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
	if err := config.DB.Model(&dto.Mst_ResourceTypes{}).Where("resource_type_id = ?", resourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data.Name, nil
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

func GetPermissionAction(permissionsIds []string) ([]string, error) {
	var actions []string
	//validate permissionIds
	for _, permissionID := range permissionsIds {
		var data dto.TNTPermission
		if err := config.DB.Model(&dto.TNTPermission{}).Where("permission_id = ?", permissionID).First(&data).Error; err != nil {
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

//	{
//	    "data": [
//	        {
//	            "key": "doc",
//	            "id": "d68ce6b3c9904dcfb0f0dec9502658f1",
//	            "organization_id": "b76ca5174a8c46488b80e21d4174eaaf",
//	            "project_id": "7c7d5c3b8eb644ea9f11f279708ce58d",
//	            "environment_id": "e92b04b3c0824b55b28073c986f8ffc4",
//	            "created_at": "2025-01-21T20:54:33+00:00",
//	            "updated_at": "2025-01-21T20:54:33+00:00",
//	            "name": "doc",
//	            "urn": "prn:devPatil:dev:doc",
//	            "description": null,
//	            "actions": {
//	                "getbuckets": {
//	                    "name": "getbuckets",
//	                    "description": null,
//	                    "attributes": {},
//	                    "id": "58a0361462b6452799ed02adf98dbe20",
//	                    "key": "getbuckets"
//	                },
//	                "fff": {
//	                    "name": "ff",
//	                    "description": null,
//	                    "attributes": {},
//	                    "id": "d0e82e74d3d743988427c0cb2f8e1088",
//	                    "key": "fff"
//	                },
//	                "read": {
//	                    "name": "read",
//	                    "description": null,
//	                    "attributes": null,
//	                    "id": "baf899c077df47b7b87b30828cd415e1",
//	                    "key": "read"
//	                },
//	                "create": {
//	                    "name": "create",
//	                    "description": null,
//	                    "attributes": null,
//	                    "id": "3922b425f9a4416c97ad10d587e477b1",
//	                    "key": "create"
//	                },
//	                "update": {
//	                    "name": "update",
//	                    "description": null,
//	                    "attributes": null,
//	                    "id": "6544e41053a344f9bc1412a1f17b5acb",
//	                    "key": "update"
//	                },
//	                "delete": {
//	                    "name": "delete",
//	                    "description": null,
//	                    "attributes": null,
//	                    "id": "51d7c8387ea241bbb5e675bfa2cf01b5",
//	                    "key": "delete"
//	                }
//	            },
//	            "type_attributes": null,
//	            "attributes": {},
//	            "roles": {
//	                "ad": {
//	                    "name": "ad",
//	                    "description": null,
//	                    "permissions": [],
//	                    "attributes": {},
//	                    "extends": [],
//	                    "granted_to": null,
//	                    "key": "ad",
//	                    "id": "93c8506dbb4f4912abee09c7a887367f",
//	                    "organization_id": "b76ca5174a8c46488b80e21d4174eaaf",
//	                    "project_id": "7c7d5c3b8eb644ea9f11f279708ce58d",
//	                    "environment_id": "e92b04b3c0824b55b28073c986f8ffc4",
//	                    "resource_id": "d68ce6b3c9904dcfb0f0dec9502658f1",
//	                    "resource": "doc",
//	                    "created_at": "2025-01-25T04:27:43+00:00",
//	                    "updated_at": "2025-01-25T04:27:43+00:00"
//	                }
//	            },
//	            "relations": {},
//	            "action_groups": {}
//	        }
//	    ],
//	    "total_count": 1,
//	    "page_count": 1
//	}
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
