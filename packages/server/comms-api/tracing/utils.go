package tracing

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

func StartHttpServerTracerSpanWithHeader(ctx context.Context, operationName string, headers http.Header) (opentracing.Span, context.Context) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))

	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		opentracing.GlobalTracer().Inject(serverSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
		return serverSpan, opentracing.ContextWithSpan(ctx, serverSpan)
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(spanCtx))
	return serverSpan, opentracing.ContextWithSpan(ctx, serverSpan)
}

func TraceErr(span opentracing.Span, err error) {
	tracing.TraceErr(span, err)
}
