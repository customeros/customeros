package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

func UserIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := ""
		for _, header := range utils.UserIdHeaders {
			if value := c.GetHeader(header); value != "" {
				userId = value
				break
			}
		}

		// Store in gin context for later use
		c.Set("UserId", userId)
		c.Next()
	}
}
