package telemetry

import (
	"go.opentelemetry.io/otel/trace"
)

// ExtractOTelTraceID extracts OpenTelemetry trace ID from span
func ExtractOTelTraceID(span trace.Span) string {
	if span == nil {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

// ExtractOTelSpanID extracts OpenTelemetry span ID from span
func ExtractOTelSpanID(span trace.Span) string {
	if span == nil {
		return ""
	}
	return span.SpanContext().SpanID().String()
}
