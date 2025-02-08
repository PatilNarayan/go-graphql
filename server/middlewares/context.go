package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
)

const GinContextKey = "GinContextKey"

// GinContextToContextMiddleware attaches the Gin context to the request context.
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
