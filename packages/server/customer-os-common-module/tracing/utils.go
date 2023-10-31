package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func TraceErr(span opentracing.Span, err error, fields ...log.Field) {
	// Log the error with the fields
	ext.LogError(span, err, fields...)
}
