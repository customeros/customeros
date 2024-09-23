package service

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	eventcompletionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const aggregateType = "event_completion"
const eventType = "V1_EVENT_COMPLETED"

type Event struct {
	Tenant    string `json:"tenant"`
	EventType string `json:"eventType"`
	Entity    string `json:"entity"`
	EntityID  string `json:"entityId"`
}

type eventCompletionService struct {
	eventcompletionpb.UnimplementedEventCompletionGrpcServiceServer
	services       *Services
	log            logger.Logger
	aggregateStore eventstore.AggregateStore
	cfg            *config.Config
}

func NewEventCompletionService(services *Services, log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config) *eventCompletionService {
	return &eventCompletionService{
		services:       services,
		log:            log,
		aggregateStore: aggregateStore,
		cfg:            cfg,
	}
}

func (s *eventCompletionService) NotifyEventProcessed(ctx context.Context, request *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EventStoreService.NotifyEventProcessed")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, "")
	tracing.LogObjectAsJson(span, "request", request)

	// ignore request if tenant is missing
	if request.Tenant == "" {
		span.LogFields(log.String("result", "missing tenant"))
		return &emptypb.Empty{}, nil
	}

	// prepare stream id
	streamID := aggregateType + "-" + request.Tenant

	// Check if stream exists
	updateStreamMetadata := false
	err := s.aggregateStore.Exists(ctx, streamID)
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return &emptypb.Empty{}, nil
		} else {
			// This is first event for this tenant
			updateStreamMetadata = true
		}
	}

	eventData := Event{
		Tenant:    request.Tenant,
		EventType: request.EventType,
		Entity:    request.Entity,
		EntityID:  request.EntityId,
	}

	aggr := eventstore.NewCommonAggregateWithId(aggregateType, request.Tenant)
	if aggr == nil {
		tracing.TraceErr(span, errors.New("invalid aggregate"))
		return &emptypb.Empty{}, nil
	}
	aggr.SetTemporal(true)

	err = eventstore.LoadAggregate(ctx, s.aggregateStore, aggr, *eventstore.NewLoadAggregateOptions().WithSkipLoadEvents())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error loading aggregate"))
		return &emptypb.Empty{}, nil
	}

	// if aggregate version is divided by 1000, then update stream metadata
	if aggr.GetVersion()%1000 == 0 {
		updateStreamMetadata = true
	}

	event := eventstore.NewBaseEvent(aggr, eventType)
	if err = event.SetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error setting json data"))
		return &emptypb.Empty{}, nil
	}
	err = aggr.Apply(event)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error applying event"))
		return &emptypb.Empty{}, nil
	}

	err = s.aggregateStore.Save(ctx, aggr)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error saving aggregate"))
		return &emptypb.Empty{}, nil
	}

	if updateStreamMetadata {
		// 7 days in seconds
		maxAgeSeconds := 7 * 24 * 60 * 60
		streamMetadata := esdb.StreamMetadata{}
		streamMetadata.SetMaxAge(time.Duration(maxAgeSeconds) * time.Second)

		err = s.aggregateStore.UpdateStreamMetadata(ctx, streamID, streamMetadata)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error updating stream metadata"))
			s.log.Errorf("Error while updating stream metadata: %s", err)
			return &emptypb.Empty{}, nil
		}
	}

	return &emptypb.Empty{}, nil
}
