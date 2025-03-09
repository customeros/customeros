package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

// HealthCheck provides a simple health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Status returns the current status of all mailboxes
func Status(imapService interfaces.IMAPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := imapService.Status()
		c.JSON(http.StatusOK, status)
	}
}
