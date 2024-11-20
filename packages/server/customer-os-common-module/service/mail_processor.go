package service

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

func (p *mailService) initializeTracing(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	tracing.SetDefaultServiceSpanTags(ctx, span)
	return span, ctx
}
