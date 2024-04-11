package tracing

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
)

const (
	SpanTagTenant         = "tenant"
	SpanTagUserId         = "user-id"
	SpanTagUserEmail      = "user-email"
	SpanTagEntityId       = "entity-id"
	SpanTagComponent      = "component"
	SpanTagExternalSystem = "external-system"
)

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

func setDefaultSpanTags(ctx context.Context, span opentracing.Span) {
	tenant := common.GetTenantFromContext(ctx)
	loggedInUserId := common.GetUserIdFromContext(ctx)
	loggedInUserEmail := common.GetUserEmailFromContext(ctx)
	if tenant != "" {
		span.SetTag(SpanTagTenant, tenant)
	}
	if loggedInUserId != "" {
		span.SetTag(SpanTagUserId, loggedInUserId)
	}
	if loggedInUserEmail != "" {
		span.SetTag(SpanTagUserEmail, loggedInUserEmail)
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

func TraceErr(span opentracing.Span, err error, fields ...log.Field) {
	// Log the error with the fields
	if span == nil {
		return
	}
	ext.LogError(span, err, fields...)
}

func LogObjectAsJson(span opentracing.Span, name string, object any) {
	if object == nil {
		span.LogFields(log.String(name, "nil"))
	}
	jsonObject, err := json.Marshal(object)
	if err == nil {
		span.LogFields(log.String(name, string(jsonObject)))
	} else {
		span.LogFields(log.Object(name, object))
	}
}
