package graph_low_prio

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/opentracing/opentracing-go"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
)

type GraphLowPrioSubscriber struct {
	log                      logger.Logger
	db                       *esdb.Client
	cfg                      *config.Config
	organizationEventHandler *OrganizationEventHandler
}

func NewGraphLowPrioSubscriber(log logger.Logger, db *esdb.Client, services *service.Services, grpcClients *grpc_client.Clients, cfg *config.Config) *GraphLowPrioSubscriber {
	return &GraphLowPrioSubscriber{
		log:                      log,
		db:                       db,
		cfg:                      cfg,
		organizationEventHandler: NewOrganizationEventHandler(log, services, grpcClients),
	}
}

func (s *GraphLowPrioSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.GraphLowPrioritySubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.GraphLowPrioritySubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.GraphLowPrioritySubscription.BufferSizeClient,
			},
		)
		if err != nil {
			return err
		}
		defer sub.Close()

		group.Go(s.runWorker(ctx, worker, sub, i))
	}
	return group.Wait()
}

func (s *GraphLowPrioSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *GraphLowPrioSubscriber) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {

	for {
		event := stream.Recv()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if event.SubscriptionDropped != nil {
			span, _ := opentracing.StartSpanFromContext(ctx, "GraphLowPrioSubscriber.ProcessEvents")
			defer span.Finish()
			wrappedErr := errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
			tracing.TraceErr(span, wrappedErr)
			s.log.Errorf(wrappedErr.Error())
			return wrappedErr
		}

		if event.EventAppeared != nil {
			s.log.EventAppeared(s.cfg.Subscriptions.GraphLowPrioritySubscription.GroupName, event.EventAppeared.Event, workerID)

			if event.EventAppeared.Event.Event == nil {
				span, _ := opentracing.StartSpanFromContext(ctx, "GraphLowPrioSubscriber.ProcessEvents")
				defer span.Finish()
				err := errors.Wrap(errors.New("event.EventAppeared.Event.Event is nil"), "GraphLowPrioSubscriber")
				tracing.TraceErr(span, err)
				s.log.Errorf(err.Error())
			} else {
				err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
				if err != nil {
					span, _ := opentracing.StartSpanFromContext(ctx, "GraphLowPrioSubscriber.ProcessEvents")
					defer span.Finish()
					tracing.TraceErr(span, err)
					s.log.Errorf("(GraphLowPrioSubscriber.when) err: {%v}", err)

					if err := stream.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("(stream.Nack) err: {%v}", err)
						return errors.Wrap(err, "stream.Nack")
					}
				}
			}

			err := stream.Ack(event.EventAppeared.Event)
			if err != nil {
				span, _ := opentracing.StartSpanFromContext(ctx, "GraphLowPrioSubscriber.ProcessEvents")
				defer span.Finish()
				tracing.TraceErr(span, err)
				s.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			s.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (s *GraphLowPrioSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
		return nil
	}

	ctx, span := tracing.StartProjectionTracerSpan(ctx, "GrpahLowPrioSubscriber.When", evt)
	defer span.Finish()

	switch evt.GetEventType() {

	case orgevents.OrganizationRefreshLastTouchpointV1:
		return s.organizationEventHandler.OnRefreshLastTouchPointV1(ctx, evt)
	default:
		return nil
	}
}
