package tracing

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
)

const (
	SpanTagTenant                = "tenant"
	SpanTagUserId                = "user-id"
	SpanTagComponent             = "component"
	SpanTagAggregateId           = "aggregateID"
	SpanTagRedundantEventSkipped = "redundantEventSkipped"
)

func StartGrpcServerTracerSpan(ctx context.Context, operationName string) (context.Context, opentracing.Span) {
	textMapCarrierFromMetaData := GetTextMapCarrierFromMetaData(ctx)

	span, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, textMapCarrierFromMetaData)
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(span))
	ctx = opentracing.ContextWithSpan(ctx, serverSpan)
	return ctx, serverSpan
}

func StartProjectionTracerSpan(ctx context.Context, operationName string, event eventstore.Event) (context.Context, opentracing.Span) {
	textMapCarrierFromMetaData := GetTextMapCarrierFromEvent(event)

	span, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, textMapCarrierFromMetaData)
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(span))
	ctx = opentracing.ContextWithSpan(ctx, serverSpan)
	return ctx, serverSpan
}

func GetTextMapCarrierFromEvent(event eventstore.Event) opentracing.TextMapCarrier {
	metadataMap := make(opentracing.TextMapCarrier)
	err := json.Unmarshal(event.GetMetadata(), &metadataMap)
	if err != nil {
		return metadataMap
	}
	return metadataMap
}

func GetTextMapCarrierFromMetaData(ctx context.Context) opentracing.TextMapCarrier {
	metadataMap := make(opentracing.TextMapCarrier)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key := range md.Copy() {
			metadataMap.Set(key, md.Get(key)[0])
		}
	}
	return metadataMap
}

func InjectTextMapCarrier(spanCtx opentracing.SpanContext) (opentracing.TextMapCarrier, error) {
	m := make(opentracing.TextMapCarrier)
	if err := opentracing.GlobalTracer().Inject(spanCtx, opentracing.TextMap, m); err != nil {
		return nil, err
	}
	return m, nil
}

func ExtractTextMapCarrier(spanCtx opentracing.SpanContext) opentracing.TextMapCarrier {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return make(opentracing.TextMapCarrier)
	}
	return textMapCarrier
}

func TraceErr(span opentracing.Span, err error, fields ...log.Field) {
	tracing.TraceErr(span, err, fields...)
}

func LogObjectAsJson(span opentracing.Span, name string, object any) {
	tracing.LogObjectAsJson(span, name, object)
}

func SetNeo4jRepositorySpanTags(ctx context.Context, span opentracing.Span, tenant string) {
	setTenantSpanTag(span, tenant)
	span.SetTag(SpanTagComponent, constants.ComponentNeo4jRepository)
}

func SetServiceSpanTags(ctx context.Context, span opentracing.Span, tenant, userId string) {
	setTenantSpanTag(span, tenant)
	setUseridSpanTag(span, userId)
	span.SetTag(SpanTagComponent, constants.ComponentService)
}

func SetCommandHandlerSpanTags(ctx context.Context, span opentracing.Span, tenant, userId string) {
	setTenantSpanTag(span, tenant)
	setUseridSpanTag(span, userId)
	span.SetTag(SpanTagComponent, constants.ComponentService)
}

func setTenantSpanTag(span opentracing.Span, tenant string) {
	if tenant != "" {
		span.SetTag(SpanTagTenant, tenant)
	}
}

func setUseridSpanTag(span opentracing.Span, userId string) {
	if userId != "" {
		span.SetTag(SpanTagUserId, userId)
	}
}
