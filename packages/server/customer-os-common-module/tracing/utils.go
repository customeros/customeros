package tracing

import "github.com/opentracing/opentracing-go"

func TraceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_msg", err.Error())
}
