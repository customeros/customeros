package tracing

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
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
	SpanTagEntityId              = "entity-id"
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

func TraceErr(span opentracing.Span, err error, fields ...log.Field) {
	tracing.TraceErr(span, err, fields...)
}

func LogObjectAsJson(span opentracing.Span, name string, object any) {
	tracing.LogObjectAsJson(span, name, object)
}

func SetNeo4jRepositorySpanTags(ctx context.Context, span opentracing.Span, tenant string) {
	setTenantSpanTag(span, tenant)
	span.SetTag(SpanTagComponent, events.ComponentNeo4jRepository)
}

func SetServiceSpanTags(ctx context.Context, span opentracing.Span, tenant, loggedInUserId string) {
	setTenantSpanTag(span, tenant)
	setUseridSpanTag(span, loggedInUserId)
	span.SetTag(SpanTagComponent, events.ComponentService)
}

func SetCommandHandlerSpanTags(ctx context.Context, span opentracing.Span, tenant, userId string) {
	setTenantSpanTag(span, tenant)
	setUseridSpanTag(span, userId)
	span.SetTag(SpanTagComponent, events.ComponentService)
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
