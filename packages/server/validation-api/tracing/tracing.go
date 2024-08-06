package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

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
