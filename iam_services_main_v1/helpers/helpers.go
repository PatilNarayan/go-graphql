package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

func GetTenant(ctx *gin.Context) (string, error) {
	tenantID := ctx.GetHeader("tenantID")
	if tenantID == "" {
		return "", errors.New("tenantID not found in headers")
	}

	//validate uuid format
	if _, err := uuid.Parse(tenantID); err != nil {
		return "", fmt.Errorf("invalid tenantID: %w", err)
	}
	return tenantID, nil
}

func GetUserID(ctx *gin.Context) (string, error) {
	userID := ctx.GetHeader("userID")
	if userID == "" {
		return "", errors.New("userID not found in headers")
	}
	return userID, nil
}

// StructToMap converts a struct to map[string]interface{}
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	// Marshal the struct into JSON
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %w", err)
	}

	// Unmarshal JSON into map
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}

	return result, nil
}
