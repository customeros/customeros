package email_validation

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/consumers"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	email_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"golang.org/x/sync/errgroup"

	esdb "github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type EmailValidationConsumer struct {
	log               logger.Logger
	db                *esdb.Client
	cfg               *config.Config
	emailEventHandler *EmailEventHandler
}

func NewEmailValidationConsumer(log logger.Logger, db *esdb.Client, cfg *config.Config, emailCommands *commands.EmailCommands) *EmailValidationConsumer {
	return &EmailValidationConsumer{
		log: log,
		db:  db,
		cfg: cfg,
		emailEventHandler: &EmailEventHandler{
			log:           log,
			cfg:           cfg,
			emailCommands: emailCommands,
		},
	}
}

func (consumer *EmailValidationConsumer) Connect(ctx context.Context, worker consumers.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= consumer.cfg.Subscriptions.EmailValidationSubscription.PoolSize; i++ {
		sub, err := consumer.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			consumer.cfg.Subscriptions.EmailValidationSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{},
		)
		if err != nil {
			return err
		}
		defer sub.Close()

		group.Go(consumer.runWorker(ctx, worker, sub, i))
	}
	return group.Wait()
}

func (consumer *EmailValidationConsumer) runWorker(ctx context.Context, worker consumers.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (consumer *EmailValidationConsumer) ProcessEvents(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error {

	for {
		event := sub.Recv()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if event.SubscriptionDropped != nil {
			consumer.log.Errorf("(SubscriptionDropped) err: {%v}", event.SubscriptionDropped.Error)
			return errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
		}

		if event.EventAppeared != nil {
			consumer.log.ConsumedEvent(consumer.cfg.Subscriptions.EmailValidationSubscription.GroupName, event.EventAppeared.Event, workerID)

			err := consumer.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
			if err != nil {
				consumer.log.Errorf("(EmailValidationConsumer.when) err: {%v}", err)

				if err := sub.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
					consumer.log.Errorf("(stream.Nack) err: {%v}", err)
					return errors.Wrap(err, "stream.Nack")
				}
			}

			err = sub.Ack(event.EventAppeared.Event)
			if err != nil {
				consumer.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}

			consumer.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (consumer *EmailValidationConsumer) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "EmailValidationConsumer.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	switch evt.GetEventType() {

	case email_events.EmailCreatedV1:
		return consumer.emailEventHandler.OnEmailCreate(ctx, evt)
	case email_events.EmailUpdatedV1:
		return nil
	case email_events.EmailValidationFailedV1:
		return nil
	case email_events.EmailValidatedV1:
		return nil

	default:
		consumer.log.Warnf("(EmailValidationConsumer) [When unknown EventType] eventType: {%s}", evt.EventType)
		return eventstore.ErrInvalidEventType
	}
}
