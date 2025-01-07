package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/service"
)

// AskAI handles AI requests
func AskAI(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "query")
		defer span.Finish()

		// Add your AI processing logic here
		c.JSON(200, gin.H{"message": "AI endpoint hit"})
	}
}
