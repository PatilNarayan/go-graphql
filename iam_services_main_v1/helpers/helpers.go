package helpers

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetGinContext(ctx context.Context) (*gin.Context, error) {
	ginCtx := ctx.Value("GinContextKey")
	if ginCtx == nil {
		return nil, errors.New("unable to retrieve gin.Context")
	}
	return ginCtx.(*gin.Context), nil
}

func CheckValueExists(field string, fallback string) string {
	if field == "" {
		return fallback
	}
	return field
}

func GetTenantID(ctx context.Context) (*uuid.UUID, error) {
	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("gin context not found in the request")
	}
	tenantID, exists := ginCtx.Get("tenantID")

	if !exists {
		return nil, fmt.Errorf("tenant id not found in context")
	}

	switch tenantID := tenantID.(type) {
	case string:
		parsedTenantID, err := uuid.Parse(tenantID)
		if err != nil {
			return nil, fmt.Errorf("error parsing tenant id: %w", err)
		}
		return &parsedTenantID, nil
	case uuid.UUID:
		return &tenantID, nil
	default:
		return nil, fmt.Errorf("invalid tenant id type")
	}
}

func GetUserID(ctx context.Context) (*uuid.UUID, error) {
	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("gin context not found in the request")
	}
	userID, exists := ginCtx.Get("userID")

	if !exists {
		return nil, fmt.Errorf("user id not found in context")
	}

	switch userID := userID.(type) {
	case string:
		parseduserID, err := uuid.Parse(userID)
		if err != nil {
			return nil, fmt.Errorf("error parsing user id: %w", err)
		}
		return &parseduserID, nil
	case uuid.UUID:
		return &userID, nil
	default:
		return nil, fmt.Errorf("invalid user id type")
	}
}

func StructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(input)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		//Skip if field is nil
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		//convert nested struct
		if fieldValue.Kind() == reflect.Struct {
			result[field.Name] = StructToMap(fieldValue.Interface())
		} else {
			result[field.Name] = fieldValue.Interface()
		}

	}

	return result
}

func MergeMaps(existing, updates map[string]interface{}) map[string]interface{} {
	for key, newValue := range updates {
		if newValue == nil {
			// skip nil values in updates
			continue
		}
		// check if the value is a map and recursively merge
		existingValue, exists := existing[key]

		if exists {
			if existingMap, ok := existingValue.(map[string]interface{}); ok {
				if newMap, ok := newValue.(map[string]interface{}); ok {
					existing[key] = MergeMaps(existingMap, newMap)
					continue
				}
			}
		}
		//overwrite or add the value from updates
		existing[key] = newValue
	}
	return existing
}

// Automate struct mapping using reflection
func MapStruct(src interface{}, dst interface{}) error {
	// Get the value of the source struct and destination struct
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	// Ensure both are pointers (to modify the destination)
	if srcValue.Kind() != reflect.Ptr || dstValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both source and destination must be pointers")
	}

	// Dereference the pointers to work with the actual values
	srcValue = srcValue.Elem()
	dstValue = dstValue.Elem()

	// Ensure both are structs
	if srcValue.Kind() != reflect.Struct || dstValue.Kind() != reflect.Struct {
		return fmt.Errorf("both source and destination must be structs")
	}

	// Iterate over the fields of the source struct
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)

		// If the destination field is valid and can be set, copy the value
		if dstField.IsValid() && dstField.CanSet() {
			// Only copy the value if the field types match
			if srcField.Type() == dstField.Type() {
				dstField.Set(srcField)
			}
		}
	}
	return nil
}
