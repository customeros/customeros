package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/core-crm/internal/utils"
)

type ginContextKey struct{}

func TenantValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant := ""
		for _, header := range utils.TenantHeaders {
			if value := c.GetHeader(header); value != "" {
				tenant = value
				break
			}
		}

		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant header is required"})
			c.Abort()
			return
		}

		// Store in gin context for later use
		c.Set("Tenant", tenant)
		c.Next()
	}
}
