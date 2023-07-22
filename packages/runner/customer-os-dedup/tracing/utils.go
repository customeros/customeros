package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

const (
	SpanTagTenant    = "tenant"
	SpanTagComponent = "component"
	SpanTagSource    = "source"
)

const ComponentNeo4jRepository = "neo4jRepository"

func StartTracerSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
	return serverSpan, opentracing.ContextWithSpan(ctx, serverSpan)
}

func TraceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_msg", err.Error())
}

func SetDefaultNeo4jRepositorySpanTags(span opentracing.Span) {
	span.SetTag(SpanTagComponent, ComponentNeo4jRepository)
}
