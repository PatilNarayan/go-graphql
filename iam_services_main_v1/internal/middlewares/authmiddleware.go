package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Define the secret key to validate JWTs (replace with your actual secret)
var jwtSecretKey []byte

// JWTClaims structure
type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// // Retrieve the secret key from the environment variable
		// jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
		// if jwtSecretKey == nil || len(jwtSecretKey) == 0 {
		// 	log.Fatal("JWT_SECRET_KEY not set in environment variables")
		// }
		// // Get the token from the Authorization header
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		// 	c.Abort()
		// 	return
		// }

		// // Check if the token is prefixed with "Bearer"
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// if tokenString == authHeader {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token missing"})
		// 	c.Abort()
		// 	return
		// }

		// // Parse and validate the token
		// token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 	// Ensure the signing method matches the expected method
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// 	}
		// 	return jwtSecretKey, nil
		// })

		// // Handle errors during parsing
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// 	c.Abort()
		// 	return
		// }

		// Validate the token claims (e.g., expiration)
		// claims, ok := token.Claims.(*JWTClaims)
		// if ok && token.Valid {
		// 	// Set the user ID in the request context
		// 	c.Set("userID", claims.UserID)

		// 	tenantID := c.GetHeader("X-Tenant-ID")

		// 	if tenantID != "" {
		// 		c.Set("tenantID", tenantID)
		// 	}

		// 	// Proceed to the next handler
		// 	c.Next()
		// } else {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		// 	c.Abort()
		// 	return
		// }
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID != "" {
			c.Set("tenantID", tenantID)
		}
		userID := c.GetHeader("userID")
		if userID != "" {
			c.Set("userID", userID)
		}

		fmt.Println("userID", userID)
		// Proceed to the next handler
	}
}
