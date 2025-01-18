package service

import (
	"context"

	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	registry "github.com/openline-ai/openline-customer-os/packages/server/events/event/_registry"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	generic "github.com/openline-ai/openline-customer-os/packages/server/events/services"

	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
)

type eventStoreService struct {
	log                      logger.Logger
	aggregateStore           eventstore.AggregateStore
	eventStoreGenericService generic.EventStoreGenericService
}

func NewEventStoreService(log logger.Logger, aggregateStore eventstore.AggregateStore, genericSrv generic.EventStoreGenericService) *eventStoreService {
	return &eventStoreService{
		log:                      log,
		aggregateStore:           aggregateStore,
		eventStoreGenericService: genericSrv,
	}
}

func (s *eventStoreService) StoreEvent(ctx context.Context, request *eventstorepb.StoreEventGrpcRequest) (*eventstorepb.StoreEventGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EventStoreService.StoreEvent")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "request", request)

	eventPayload, err := registry.UnmarshalBaseEventPayload(request.GetEventDataBytes())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, grpcerr.ErrResponse(err)
	}

	entityId, err := s.eventStoreGenericService.Store(ctx, eventPayload, eventstore.LoadAggregateOptions{})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, grpcerr.ErrResponse(err)
	}

	return &eventstorepb.StoreEventGrpcResponse{Id: *entityId}, nil
}

func (s *eventStoreService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
