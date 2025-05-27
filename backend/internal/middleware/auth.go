package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/readitlater/backend/internal/auth"
	"github.com/readitlater/backend/internal/services"
)

// Auth middleware validates JWT tokens
func Auth(authService *auth.Service, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Development bypass for placeholder tokens
		if tokenString == "placeholder-jwt-token" || tokenString == "new-placeholder-jwt-token" {
			// Set placeholder user info in context for development
			c.Set("user_id", "placeholder-user-id")
			c.Set("user_email", "user@example.com")
			c.Set("user_username", "placeholder-user")
			c.Next()
			return
		}

		claims, err := authService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)

		c.Next()
	}
}

// ExtensionAuth middleware for browser extension API
func ExtensionAuth(cfg interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Placeholder for extension authentication
		c.Next()
	}
}

// RateLimit middleware for rate limiting
func RateLimit(cfg interface{}, redisClient interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Placeholder for rate limiting
		c.Next()
	}
} 