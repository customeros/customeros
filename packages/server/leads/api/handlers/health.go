package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck provides a simple health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
