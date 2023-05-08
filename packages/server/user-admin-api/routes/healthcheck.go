package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func addHealthRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
