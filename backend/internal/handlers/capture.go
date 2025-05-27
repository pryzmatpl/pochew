package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/readitlater/backend/internal/services"
)

// CaptureHandler handles content capture endpoints
type CaptureHandler struct {
	captureService *services.CaptureService
	logger         *logrus.Logger
}

// NewCaptureHandler creates a new capture handler
func NewCaptureHandler(captureService *services.CaptureService, logger *logrus.Logger) *CaptureHandler {
	return &CaptureHandler{
		captureService: captureService,
		logger:         logger,
	}
}

// CaptureURL captures content from a URL
func (h *CaptureHandler) CaptureURL(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		URL string `json:"url" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Placeholder implementation
	c.JSON(http.StatusAccepted, gin.H{
		"capture_id": "placeholder-capture-id",
		"status":     "processing",
		"url":        req.URL,
	})
}

// GetCaptureStatus retrieves the status of a capture operation
func (h *CaptureHandler) GetCaptureStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	captureID := c.Param("id")
	
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Placeholder implementation
	c.JSON(http.StatusOK, gin.H{
		"id":       captureID,
		"status":   "completed",
		"progress": 100,
	})
}

// ExtensionCapture handles capture requests from browser extension
func (h *CaptureHandler) ExtensionCapture(c *gin.Context) {
	var req struct {
		URL     string `json:"url" binding:"required"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Placeholder implementation for extension capture
	c.JSON(http.StatusCreated, gin.H{
		"id":     "extension-article-id",
		"url":    req.URL,
		"title":  req.Title,
		"status": "saved",
	})
} 