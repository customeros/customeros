package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

func TracingEnhancer(ctx context.Context, endpoint string) func(c *gin.Context) {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(ctx, endpoint, c.Request.Header)
		for k, v := range c.Request.Header {
			span.LogFields(log.String("request.header.key", k), log.Object("request.header.value", v))
		}
		defer span.Finish()
		c.Next()
	}
}
