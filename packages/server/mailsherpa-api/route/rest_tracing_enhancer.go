package route

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func RestTracingEnhancer(ctx context.Context, endpoint string) func(c *gin.Context) {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartHttpServerTracerSpanWithHeader(ctx, endpoint, c.Request.Header)
		defer spans.Finish()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
