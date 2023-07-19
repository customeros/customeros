package tracing

import (
	"context"
	local_utils "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"github.com/opentracing/opentracing-go"
)

const ComponentNeo4jRepository = "neo4jRepository"
const ComponentPostgresRepository = "postgresRepository"
const ComponentSyncService = "syncService"

func StartTracerSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
	return serverSpan, opentracing.ContextWithSpan(ctx, serverSpan)
}

func TraceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_msg", err.Error())
}

func setDefaultSpanTags(ctx context.Context, span opentracing.Span) {
	if local_utils.GetTenantFromContext(ctx) != "" {
		span.SetTag(SpanTagTenant, local_utils.GetTenantFromContext(ctx))
	}
	if local_utils.GetSourceFromContext(ctx) != "" {
		span.SetTag(SpanTagSource, local_utils.GetSourceFromContext(ctx))
	}
}

func SetDefaultSyncServiceSpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(SpanTagComponent, ComponentSyncService)
}

func SetDefaultPostgresRepositorySpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(SpanTagComponent, ComponentSyncService)
}

func SetDefaultNeo4jRepositorySpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(SpanTagComponent, ComponentNeo4jRepository)
}
