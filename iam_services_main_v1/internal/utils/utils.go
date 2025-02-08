package utils

import (
	"errors"
	"regexp"
)

func UpdateDeletedMap() map[string]interface{} {
	return map[string]interface{}{
		"row_status": 0,
	}
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
