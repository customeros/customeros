package tracing

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
	"net/http"
)

const (
	SpanTagTenant    = tracing.SpanTagTenant
	SpanTagComponent = tracing.SpanTagComponent
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
	if tenant != "" {
		span.SetTag(tracing.SpanTagTenant, tenant)
	}
	if loggedInUserId != "" {
		span.SetTag(tracing.SpanTagUserId, loggedInUserId)
	}
}

func SetDefaultResolverSpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentResolver)
}

func SetDefaultServiceSpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
}

func SetDefaultNeo4jRepositorySpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentNeo4jRepository)
}

func LogObjectAsJson(span opentracing.Span, name string, object any) {
	tracing.LogObjectAsJson(span, name, object)
}

func InjectSpanIntoGrpcRequestMetadata(ctx context.Context, span opentracing.Span) context.Context {
	if span != nil {
		// Inject the span context into the gRPC request metadata.
		textMapCarrier := make(opentracing.TextMapCarrier)
		err := span.Tracer().Inject(span.Context(), opentracing.TextMap, textMapCarrier)
		if err == nil {
			// Add the injected metadata to the gRPC context.
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.New(nil)
			}
			for key, val := range textMapCarrier {
				md.Set(key, val)
			}
			ctx = metadata.NewOutgoingContext(ctx, md)
			return ctx
		}
	}
	return ctx
}
