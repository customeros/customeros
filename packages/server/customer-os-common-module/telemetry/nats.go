package telemetry

import (
	"context"
	"strings"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// InjectTraceContextIntoNatsMsg injects the current trace context into NATS message headers
func InjectTraceContextIntoNatsMsg(ctx context.Context, msg *nats.Msg) {
	if msg == nil {
		return
	}

	// Initialize headers if nil
	if msg.Header == nil {
		msg.Header = nats.Header{}
	}

	// Get the current span context
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		// Try to get span from context value
		if span, ok := ctx.Value(otelSpanKey).(trace.Span); ok && span != nil {
			spanCtx = span.SpanContext()
		} else {
			return
		}
	}

	// Create a context with just the span context
	ctx = trace.ContextWithSpanContext(context.Background(), spanCtx)

	// Get the propagator and inject the context into headers
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.HeaderCarrier(msg.Header)
	propagator.Inject(ctx, carrier)
}

// ExtractTraceContextFromNatsMsg extracts the trace context from NATS message headers
func ExtractTraceContextFromNatsMsg(ctx context.Context, msg *nats.Msg) context.Context {
	if msg == nil || msg.Header == nil {
		return ctx
	}

	// Get the propagator and extract the context from headers
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.HeaderCarrier(msg.Header)

	// Extract into a new context to avoid mixing with existing context
	extractedCtx := propagator.Extract(context.Background(), carrier)

	// Get the span context from the extracted context
	spanCtx := trace.SpanContextFromContext(extractedCtx)
	if !spanCtx.IsValid() {
		// Try to get the traceparent header directly
		traceparent := msg.Header.Get("traceparent")
		if traceparent != "" {
			// The traceparent format is: version-traceid-spanid-flags
			parts := strings.Split(traceparent, "-")
			if len(parts) >= 4 {
				if traceID, err := trace.TraceIDFromHex(parts[1]); err == nil {
					if spanID, err := trace.SpanIDFromHex(parts[2]); err == nil {
						// Create a new span context with the parsed values
						spanCtx = trace.NewSpanContext(trace.SpanContextConfig{
							TraceID:    traceID,
							SpanID:     spanID,
							TraceFlags: trace.FlagsSampled,
						})
					}
				}
			}
		}
	}

	if !spanCtx.IsValid() {
		return ctx
	}

	// Create a new context with the extracted span context
	return trace.ContextWithSpanContext(ctx, spanCtx)
}
