package email_validation

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	email_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/projection"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type EmailValidationProjection struct {
	log               logger.Logger
	db                *esdb.Client
	cfg               *config.Config
	emailEventHandler *EmailEventHandler
}

func NewEmailValidationProjection(log logger.Logger, db *esdb.Client, cfg *config.Config, emailCommands *commands.EmailCommands) *EmailValidationProjection {
	return &EmailValidationProjection{
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

func (evp *EmailValidationProjection) Subscribe(ctx context.Context, prefixes []string, poolSize int, worker projection.Worker) error {
	evp.log.Infof("(starting email validation subscription) prefixes: {%+v}", prefixes)

	err := evp.db.CreatePersistentSubscriptionAll(ctx, evp.cfg.Subscriptions.EmailValidationProjectionGroupName, esdb.PersistentAllSubscriptionOptions{
		Filter: &esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: prefixes},
	})
	if err != nil {
		if subscriptionError, ok := err.(*esdb.PersistentSubscriptionError); !ok || ok && (subscriptionError.Code != 6) {
			evp.log.Errorf("(EmailValidationProjection.CreatePersistentSubscriptionAll) err: {%v}", subscriptionError.Error())
		} else if ok && (subscriptionError.Code == 6) {
			// FIXME alexb refactor: call update only if current and new prefixes are different
			settings := esdb.SubscriptionSettingsDefault()
			err = evp.db.UpdatePersistentSubscriptionAll(ctx, evp.cfg.Subscriptions.EmailValidationProjectionGroupName, esdb.PersistentAllSubscriptionOptions{
				Settings: &settings,
				Filter:   &esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: prefixes},
			})
			if err != nil {
				if subscriptionError, ok = err.(*esdb.PersistentSubscriptionError); !ok || ok && (subscriptionError.Code != 6) {
					evp.log.Errorf("(EmailValidationProjection.UpdatePersistentSubscriptionAll) err: {%v}", subscriptionError.Error())
				}
			}
		}
	}

	stream, err := evp.db.ConnectToPersistentSubscription(
		ctx,
		constants.EsAll,
		evp.cfg.Subscriptions.EmailValidationProjectionGroupName,
		esdb.ConnectToPersistentSubscriptionOptions{},
	)
	if err != nil {
		return err
	}
	defer stream.Close()

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i <= poolSize; i++ {
		g.Go(evp.runWorker(ctx, worker, stream, i))
	}
	return g.Wait()
}

func (evp *EmailValidationProjection) runWorker(ctx context.Context, worker projection.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (evp *EmailValidationProjection) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {

	for {
		event := stream.Recv()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if event.SubscriptionDropped != nil {
			evp.log.Errorf("(SubscriptionDropped) err: {%v}", event.SubscriptionDropped.Error)
			return errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
		}

		if event.EventAppeared != nil {
			evp.log.ProjectionEvent(constants.EmailValidationProjection, evp.cfg.Subscriptions.EmailValidationProjectionGroupName, event.EventAppeared, workerID)

			err := evp.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event))
			if err != nil {
				evp.log.Errorf("(EmailValidationProjection.when) err: {%v}", err)

				if err := stream.Nack(err.Error(), esdb.Nack_Retry, event.EventAppeared); err != nil {
					evp.log.Errorf("(stream.Nack) err: {%v}", err)
					return errors.Wrap(err, "stream.Nack")
				}
			}

			err = stream.Ack(event.EventAppeared)
			if err != nil {
				evp.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			evp.log.Infof("(ACK) event commit: {%v}", *event.EventAppeared.Commit)
		}
	}
}

func (evp *EmailValidationProjection) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "EmailValidationProjection.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	switch evt.GetEventType() {

	case email_events.EmailCreatedV1:
		return evp.emailEventHandler.OnEmailCreate(ctx, evt)
	case email_events.EmailValidationFailedV1:
		return nil
	case email_events.EmailValidatedV1:
		return nil
	case "PersistentConfig1":
		return nil

	default:
		// FIXME alexb if event was not recognized, park it
		evp.log.Warnf("(EmailValidationProjection) [When unknown EventType] eventType: {%s}", evt.EventType)
		return eventstore.ErrInvalidEventType
	}
}
