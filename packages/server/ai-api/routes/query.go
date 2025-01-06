package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/service"
)

func query(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "query")
		defer span.Finish()

		return
	}
}
