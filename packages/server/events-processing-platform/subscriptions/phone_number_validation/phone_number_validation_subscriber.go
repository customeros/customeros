package phone_number_validation

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"golang.org/x/sync/errgroup"
	"strings"

	esdb "github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type PhoneNumberValidationSubscriber struct {
	log                     logger.Logger
	db                      *esdb.Client
	cfg                     *config.Config
	phoneNumberEventHandler *PhoneNumberEventHandler
	repositories            *repository.Repositories
}

func NewPhoneNumberValidationSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, phoneNumberCommands *command_handler.PhoneNumberCommands, repositories *repository.Repositories) *PhoneNumberValidationSubscriber {
	return &PhoneNumberValidationSubscriber{
		log: log,
		db:  db,
		cfg: cfg,
		phoneNumberEventHandler: &PhoneNumberEventHandler{
			log:                 log,
			cfg:                 cfg,
			phoneNumberCommands: phoneNumberCommands,
			repositories:        repositories,
		},
	}
}

func (s *PhoneNumberValidationSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.PhoneNumberValidationSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.PhoneNumberValidationSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.PhoneNumberValidationSubscription.BufferSizeClient,
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

func (s *PhoneNumberValidationSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *PhoneNumberValidationSubscriber) ProcessEvents(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error {

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
			s.log.EventAppeared(s.cfg.Subscriptions.PhoneNumberValidationSubscription.GroupName, event.EventAppeared.Event, workerID)

			err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
			if err != nil {
				s.log.Errorf("(PhoneNumberValidationSubscriber.when) err: {%v}", err)

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

func (s *PhoneNumberValidationSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "PhoneNumberValidationSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	if strings.HasPrefix(evt.GetAggregateID(), "$") {
		return nil
	}

	switch evt.GetEventType() {
	case
		events.PhoneNumberCreateV1:
		return s.phoneNumberEventHandler.OnPhoneNumberCreate(ctx, evt)
	case
		events.PhoneNumberUpdateV1,
		events.PhoneNumberValidationFailedV1,
		events.PhoneNumberValidationSkippedV1,
		events.PhoneNumberValidatedV1:
		return nil

	default:
		s.log.Warnf("(PhoneNumberValidationSubscriber) Unknown EventType: {%s}", evt.EventType)
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}
