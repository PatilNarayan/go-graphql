package helpers

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

func GetGinContext(ctx context.Context) (*gin.Context, error) {
	ginCtx := ctx.Value("GinContextKey")
	if ginCtx == nil {
		return nil, errors.New("unable to retrieve gin.Context")
	}
	return ginCtx.(*gin.Context), nil
}
