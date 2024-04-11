package tracing

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

const (
	SpanTagComponent = tracing.SpanTagComponent
	SpanTagEntityId  = tracing.SpanTagEntityId
)

type spanCtxKey struct{}

var activeSpanCtxKey = spanCtxKey{}

func ExtractSpanCtx(ctx context.Context) opentracing.SpanContext {
	if ctx.Value(activeSpanCtxKey) != nil {
		return ctx.Value(activeSpanCtxKey).(opentracing.SpanContext)
	}
	return nil
}

func EnrichCtxWithSpanCtxForGraphQL(ctx context.Context, operationContext *graphql.OperationContext) context.Context {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(operationContext.Headers))
	if err != nil {
		return ctx
	}
	if ExtractSpanCtx(ctx) != nil {
		return ctx
	}
	return context.WithValue(ctx, activeSpanCtxKey, spanCtx)
}

func StartGraphQLTracerSpan(ctx context.Context, operationName string, operationContext *graphql.OperationContext) (context.Context, opentracing.Span) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(operationContext.Headers))

	if err != nil {
		rootSpan := opentracing.GlobalTracer().StartSpan(operationName)
		opentracing.GlobalTracer().Inject(rootSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(operationContext.Headers))
		return opentracing.ContextWithSpan(ctx, rootSpan), rootSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(spanCtx))
	return opentracing.ContextWithSpan(ctx, serverSpan), serverSpan
}

func StartHttpServerTracerSpanWithHeader(ctx context.Context, operationName string, headers http.Header) (context.Context, opentracing.Span) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))

	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		opentracing.GlobalTracer().Inject(serverSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
		return opentracing.ContextWithSpan(ctx, serverSpan), serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(spanCtx))
	return opentracing.ContextWithSpan(ctx, serverSpan), serverSpan
}

func TraceErr(span opentracing.Span, err error, fields ...log.Field) {
	tracing.TraceErr(span, err, fields...)
}

func setDefaultSpanTags(ctx context.Context, span opentracing.Span) {
	tenant := common.GetTenantFromContext(ctx)
	loggedInUserId := common.GetUserIdFromContext(ctx)
	loggedInUserEmail := common.GetUserEmailFromContext(ctx)
	if tenant != "" {
		span.SetTag(tracing.SpanTagTenant, tenant)
	}
	if loggedInUserId != "" {
		span.SetTag(tracing.SpanTagUserId, loggedInUserId)
	}
	if loggedInUserEmail != "" {
		span.SetTag(tracing.SpanTagUserEmail, loggedInUserEmail)
	}
}

func SetDefaultResolverSpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentResolver)
}

func LogObjectAsJson(span opentracing.Span, name string, object any) {
	tracing.LogObjectAsJson(span, name, object)
}
