package routes

import (
	"github.com/gin-gonic/gin"
)

func addHealthRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", healthCheckHandler)
	rg.GET("/readiness", healthCheckHandler)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}
