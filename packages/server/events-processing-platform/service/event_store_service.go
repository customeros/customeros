package service

import (
	"context"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	registry "github.com/openline-ai/openline-customer-os/packages/server/events/event/_registry"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type eventStoreService struct {
	services       *Services
	log            logger.Logger
	aggregateStore eventstore.AggregateStore
}

func NewEventStoreService(services *Services, log logger.Logger, aggregateStore eventstore.AggregateStore) *eventStoreService {
	return &eventStoreService{
		services:       services,
		log:            log,
		aggregateStore: aggregateStore,
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

	entityId, err := s.services.EventStoreGenericService.Store(ctx, eventPayload, eventstore.LoadAggregateOptions{})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, grpcerr.ErrResponse(err)
	}

	return &eventstorepb.StoreEventGrpcResponse{Id: *entityId}, nil
}

func (s *eventStoreService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
