package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
)

// TracingMiddleware creates a new span for each request and adds common tags
func TracingMiddleware(parentCtx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get existing custom context if any
		existingCtx := c.Request.Context()

		// Start OpenTelemetry span
		tracer := otel.Tracer("github.com/customeros/mailstack")
		otelCtx, otelSpan := tracer.Start(existingCtx, c.Request.Method+" "+c.FullPath(),
			trace.WithAttributes(
				telemetry.GetDefaultServiceSpanAttributes(existingCtx)...,
			),
		)

		// Extract OpenTelemetry trace context from headers
		otelCtx = otel.GetTextMapPropagator().Extract(otelCtx, propagation.HeaderCarrier(c.Request.Header))
		defer otelSpan.End()

		// Create Spans struct for telemetry operations
		spans := &telemetry.Spans{
			OTel: otelSpan,
		}
		// Tag as REST component for both spans
		telemetry.TagComponentRest(spans)

		// Set default span tags (tenant, user-id, user-email)
		telemetry.SetDefaultServiceSpanAttributes(otelCtx, otelSpan)

		// Add entity ID if present in URL params
		if id := c.Param("id"); id != "" {
			spans.TagEntity(id)
		}

		// Store both spans in context while preserving existing context values
		c.Request = c.Request.WithContext(otelCtx)

		// Process request
		c.Next()

		// Add response status
		if c.Writer.Status() >= 400 {
			spans.TraceError(nil)
		}
	}
}
