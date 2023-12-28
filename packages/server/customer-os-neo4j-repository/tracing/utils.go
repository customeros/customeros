package tracing

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const neo4jRepository = "neo4jRepository"

const SpanTagEntityId = tracing.SpanTagEntityId

func setTenantSpanTag(span opentracing.Span, tenant string) {
	if tenant != "" {
		span.SetTag(tracing.SpanTagTenant, tenant)
	}
}

func SetNeo4jRepositorySpanTags(span opentracing.Span, tenant string) {
	setTenantSpanTag(span, tenant)
	span.SetTag(tracing.SpanTagComponent, neo4jRepository)
}

func LogObjectAsJson(span opentracing.Span, name string, object any) {
	tracing.LogObjectAsJson(span, name, object)
}

func TraceErr(span opentracing.Span, err error, fields ...log.Field) {
	tracing.TraceErr(span, err, fields...)
}
