package tracing

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

const (
	SpanTagComponent = "component"
)

const ComponentNeo4jRepository = "neo4jRepository"

func StartTracerSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
	return serverSpan, opentracing.ContextWithSpan(ctx, serverSpan)
}

func TraceErr(span opentracing.Span, err error) {
	tracing.TraceErr(span, err)
}

func SetDefaultNeo4jRepositorySpanTags(span opentracing.Span) {
	span.SetTag(SpanTagComponent, ComponentNeo4jRepository)
}

func SetDefaultServiceSpanTags(span opentracing.Span) {
	span.SetTag(SpanTagComponent, ComponentNeo4jRepository)
}
