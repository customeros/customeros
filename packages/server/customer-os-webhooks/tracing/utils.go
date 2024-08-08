package tracing

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/opentracing/opentracing-go"
)

const (
	SpanTagTenant         = tracing.SpanTagTenant
	SpanTagComponent      = tracing.SpanTagComponent
	SpanTagExternalSystem = tracing.SpanTagExternalSystem
)

func setDefaultSpanTags(ctx context.Context, span opentracing.Span) {
	tenant := common.GetTenantFromContext(ctx)
	if tenant != "" {
		span.SetTag(SpanTagTenant, tenant)
	}
}

func SetDefaultServiceSpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(SpanTagComponent, constants.ComponentService)
}

func SetDefaultNeo4jRepositorySpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(SpanTagComponent, constants.ComponentNeo4jRepository)
}
