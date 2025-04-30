package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

// CustomContextMiddleware adds custom context to all requests
func CustomContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := utils.WithCustomContextFromGinRequest(c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinContextToContextMiddleware adds the Gin context to the request context
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey{}, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinContextFromContext extracts the Gin context from a context
func GinContextFromContext(ctx context.Context) (*gin.Context, bool) {
	ginContext, ok := ctx.Value(ginContextKey{}).(*gin.Context)
	return ginContext, ok
}
