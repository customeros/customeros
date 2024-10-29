package utils

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventCompletedDetails struct {
	Create bool
	Update bool
	Delete bool
}

func NewEventCompletedDetails() *EventCompletedDetails {
	return &EventCompletedDetails{}
}

func (ecd *EventCompletedDetails) WithCreate() *EventCompletedDetails {
	ecd.Create = true
	return ecd
}

func (ecd *EventCompletedDetails) WithUpdate() *EventCompletedDetails {
	ecd.Update = true
	return ecd
}

func (ecd *EventCompletedDetails) WithDelete() *EventCompletedDetails {
	ecd.Delete = true
	return ecd
}

// deprecated
func EventCompleted(ctx context.Context, tenant, entity, entityId string, grpcClients *grpc_client.Clients, details *EventCompletedDetails) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventCompleted")
	defer span.Finish()
	span.LogKV("tenant", tenant, "entity", entity, "entityID", entityId)

	request := event_completion_grpc_service.NotifyEventProcessedRequest{
		Tenant:   tenant,
		Entity:   entity,
		EntityId: entityId,
	}
	if details != nil {
		request.Create = details.Create
		request.Update = details.Update
		request.Delete = details.Delete
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := CallEventsPlatformGRPCWithRetry[*emptypb.Empty](func() (*emptypb.Empty, error) {
		return grpcClients.EventCompletionClient.NotifyEventProcessed(ctx, &request)
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
}
