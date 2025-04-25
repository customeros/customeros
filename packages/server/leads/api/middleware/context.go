package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/customeros/mailstack/internal/utils"
)

// CustomContextMiddleware adds custom context to all requests
func CustomContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := utils.WithCustomContextFromGinRequest(c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
