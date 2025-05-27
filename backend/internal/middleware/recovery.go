package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery middleware recovers from panics and logs them
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.WithFields(logrus.Fields{
			"error":      recovered,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Error("Panic recovered")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
} 