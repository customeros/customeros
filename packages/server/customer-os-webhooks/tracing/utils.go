package tracing

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type spanCtxKey struct{}

var activeSpanCtxKey = spanCtxKey{}

const (
	SpanTagTenant    = "tenant"
	SpanTagComponent = "component"
)

func ExtractSpanCtx(ctx context.Context) opentracing.SpanContext {
	if ctx.Value(activeSpanCtxKey) != nil {
		return ctx.Value(activeSpanCtxKey).(opentracing.SpanContext)
	}
	return nil
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

func InjectSpanContextIntoGrpcMetadata(ctx context.Context, span opentracing.Span) context.Context {
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
