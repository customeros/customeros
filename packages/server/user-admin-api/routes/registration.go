package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"net/http"
)

func addRegistrationRoutes(rg *gin.RouterGroup, config *config.Config) {
	rg.POST("/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
