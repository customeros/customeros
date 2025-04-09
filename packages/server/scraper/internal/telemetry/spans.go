package telemetry

import (
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/customeros/customeros/packages/server/scraper/internal/utils"
)

// nopCloser implements io.Closer for when tracing is disabled
type nopCloser struct{}

func (nopCloser) Close() error { return nil }

// otelCloser implements io.Closer for cleanup
type otelCloser struct {
	tp *sdktrace.TracerProvider
}

func (c *otelCloser) Close() error {
	return c.tp.Shutdown(context.Background())
}

// StartHttpServerTracerSpanWithHeader creates a span from HTTP headers
func StartHttpServerTracerSpanWithHeader(ctx context.Context, operationName string, headers http.Header) (context.Context, trace.Span) {
	// Extract context from headers
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.HeaderCarrier(headers)
	ctx = propagator.Extract(ctx, carrier)

	// Start a new span
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(ctx, operationName)

	// Inject the span context back into headers
	propagator.Inject(ctx, carrier)

	return ctx, span
}

// SetDefaultSpanTags adds common tags to a span
func SetDefaultSpanTags(ctx context.Context, span trace.Span) {
	tenant := utils.GetTenantFromContext(ctx)
	loggedInUserId := utils.GetUserIdFromContext(ctx)
	loggedInUserEmail := utils.GetUserEmailFromContext(ctx)

	if tenant != "" {
		span.SetAttributes(attribute.String(SpanTagTenant, tenant))
	}
	if loggedInUserId != "" {
		span.SetAttributes(attribute.String(SpanTagUserId, loggedInUserId))
	}
	if loggedInUserEmail != "" {
		span.SetAttributes(attribute.String(SpanTagUserEmail, loggedInUserEmail))
	}
}

// TraceErr records an error in a span
func TraceErr(span trace.Span, err error, fields ...attribute.KeyValue) {
	if span == nil || err == nil {
		return
	}

	// Mark the span as error and log details
	span.RecordError(err, trace.WithAttributes(fields...))
	span.SetStatus(codes.Error, err.Error())
}

// LogObjectAsJson logs an object as JSON to a span
func LogObjectAsJson(span trace.Span, name string, object any) {
	if span == nil {
		return
	}

	var value string
	if object == nil {
		value = "nil"
	} else {
		jsonObject, err := json.Marshal(object)
		if err == nil {
			value = string(jsonObject)
		} else {
			value = "error marshaling object to JSON"
		}
	}

	span.SetAttributes(attribute.String(name, value))
}
