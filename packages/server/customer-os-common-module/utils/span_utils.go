package utils

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func GetTraceIDFromSpan(span opentracing.Span) string {
	if span == nil {
		return ""
	}
	// For Jaeger specifically
	if jaegerSpan, ok := span.(*jaeger.Span); ok {
		spanContext := jaegerSpan.Context().(jaeger.SpanContext)
		return spanContext.TraceID().String()
	} else {
		// Fallback for other tracers - get carrier with all span context info
		carrier := opentracing.TextMapCarrier{}
		err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
		if err == nil {
			// Many tracers use these standard field names
			return carrier["uber-trace-id"]
		}
	}
	return ""
}
