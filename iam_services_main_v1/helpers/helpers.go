package helpers

import (
	"context"
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
