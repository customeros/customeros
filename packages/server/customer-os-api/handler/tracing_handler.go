package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/opentracing/opentracing-go/log"
)

func TracingEnhancer(ctx context.Context, endpoint string) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxWithSpan, span := tracing.StartHttpServerTracerSpanWithHeader(ctx, endpoint, c.Request.Header)
		defer span.Finish()
		for k, v := range c.Request.Header {
			span.LogFields(log.String("request.header.key", k), log.Object("request.header.value", v))
		}
		c.Request = c.Request.WithContext(ctxWithSpan)
		c.Next()
	}
}
