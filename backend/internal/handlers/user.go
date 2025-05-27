package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/readitlater/backend/internal/models"
	"github.com/readitlater/backend/internal/services"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	userService *services.UserService
	logger      *logrus.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetProfile retrieves the current user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Placeholder implementation
	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"email":    "user@example.com",
		"username": "placeholder-user",
	})
}

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Placeholder implementation
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// ChangePassword changes the current user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.PasswordChange
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Placeholder implementation
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// GetStats retrieves user statistics
func (h *UserHandler) GetStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Placeholder implementation
	c.JSON(http.StatusOK, gin.H{
		"total_articles":    0,
		"read_articles":     0,
		"unread_articles":   0,
		"favorite_articles": 0,
		"storage_used":      0,
	})
}

// GetExtensionUser retrieves user info for browser extension
func (h *UserHandler) GetExtensionUser(c *gin.Context) {
	// Placeholder implementation for extension API
	c.JSON(http.StatusOK, gin.H{
		"id":       "extension-user-id",
		"username": "extension-user",
	})
} 