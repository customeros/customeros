package utils

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func EventCompleted(ctx context.Context, tenant, entity, entityId, entityType string, grpcClients *grpc_client.Clients) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventCompleted")
	defer span.Finish()
	span.LogKV("tenant", tenant, "entity", entity, "entityID", entityId, "entityType", entityType)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := CallEventsPlatformGRPCWithRetry[*emptypb.Empty](func() (*emptypb.Empty, error) {
		return grpcClients.EventCompletionClient.NotifyEventProcessed(ctx, &event_completion_grpc_service.NotifyEventProcessedRequest{
			Tenant:    tenant,
			EventType: entityType,
			Entity:    entity,
			EntityId:  entityId,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
}
