package tracing

import (
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	SpanTagTenant         = "tenant"
	SpanTagUserId         = "user-id"
	SpanTagUserEmail      = "user-email"
	SpanTagEntityId       = "entity-id"
	SpanTagComponent      = "component"
	SpanTagExternalSystem = "external-system"
)

func TraceErr(span opentracing.Span, err error, fields ...log.Field) {
	// Log the error with the fields
	if span == nil {
		return
	}
	ext.LogError(span, err, fields...)
}

func LogObjectAsJson(span opentracing.Span, name string, object any) {
	if object == nil {
		span.LogFields(log.String(name, "nil"))
	}
	jsonObject, err := json.Marshal(object)
	if err == nil {
		span.LogFields(log.String(name, string(jsonObject)))
	} else {
		span.LogFields(log.Object(name, object))
	}
}
