package tracing

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

type spanCtxKey struct{}

var activeSpanCtxKey = spanCtxKey{}

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
