package interactionEvent

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"golang.org/x/sync/errgroup"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type InteractionEventSubscriber struct {
	log                     logger.Logger
	db                      *esdb.Client
	cfg                     *config.Config
	interactionEventHandler *interactionEventHandler
}

func NewInteractionEventSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, commands *command_handler.CommandHandlers, repositories *repository.Repositories) *InteractionEventSubscriber {
	return &InteractionEventSubscriber{
		log: log,
		db:  db,
		cfg: cfg,
		interactionEventHandler: &interactionEventHandler{
			log:                      log,
			cfg:                      cfg,
			interactionEventCommands: commands,
			repositories:             repositories,
		},
	}
}

func (s *InteractionEventSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.InteractionEventSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.InteractionEventSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.InteractionEventSubscription.BufferSizeClient,
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

func (s *InteractionEventSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *InteractionEventSubscriber) ProcessEvents(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error {

	for {
		event := sub.Recv()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if event.SubscriptionDropped != nil {
			s.log.Errorf("(SubscriptionDropped) err: {%v}", event.SubscriptionDropped.Error)
			return errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
		}

		if event.EventAppeared != nil {
			s.log.EventAppeared(s.cfg.Subscriptions.InteractionEventSubscription.GroupName, event.EventAppeared.Event, workerID)

			err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
			if err != nil {
				s.log.Errorf("(InteractionEventSubscriber.when) err: {%v}", err)

				if err := sub.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
					s.log.Errorf("(stream.Nack) err: {%v}", err)
					return errors.Wrap(err, "stream.Nack")
				}
			}

			err = sub.Ack(event.EventAppeared.Event)
			if err != nil {
				s.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}

			s.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (s *InteractionEventSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "InteractionEventSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	if strings.HasPrefix(evt.GetAggregateID(), "$") {
		return nil
	}

	switch evt.GetEventType() {
	case event.InteractionEventRequestSummaryV1:
		return s.interactionEventHandler.GenerateSummaryForEmail(ctx, evt)
	case event.InteractionEventRequestActionItemsV1:
		return s.interactionEventHandler.GenerateActionItemsForEmail(ctx, evt)
	case event.InteractionEventReplaceSummaryV1,
		event.InteractionEventReplaceActionItemsV1,
		event.InteractionEventCreateV1,
		event.InteractionEventUpdateV1:
		return nil

	default:
		tracing.TraceErr(span, eventstore.ErrInvalidEventType)
		s.log.Warnf("Unknown EventType: {%s}", evt.EventType)
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}
