package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Define a key type to avoid context key collisions
type contextKey string

const GinContextKey contextKey = "GinContextKey"

// Middleware to inject *gin.Context into the request context
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new context containing the gin.Context
		ctx := context.WithValue(c.Request.Context(), GinContextKey, c)

		// Replace the request with the new context
		c.Request = c.Request.WithContext(ctx)

		// Continue processing
		c.Next()
	}
}
