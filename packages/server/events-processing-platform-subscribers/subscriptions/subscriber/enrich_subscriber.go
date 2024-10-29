package subscriber

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	emailevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"strings"
)

type EnrichSubscriber struct {
	log               logger.Logger
	db                *esdb.Client
	cfg               *config.Config
	emailEventHandler *email.EmailEventHandler
}

func NewEnrichSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, services *service.Services, caches caches.Cache, grpcClients *grpc_client.Clients) *EnrichSubscriber {
	return &EnrichSubscriber{
		log:               log,
		db:                db,
		cfg:               cfg,
		emailEventHandler: email.NewEmailEventHandler(services, log, cfg, caches, grpcClients),
	}
}

func (s *EnrichSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.EnrichSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.EnrichSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.EnrichSubscription.BufferSizeClient,
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

func (s *EnrichSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *EnrichSubscriber) ProcessEvents(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error {

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
			s.log.EventAppeared(s.cfg.Subscriptions.EnrichSubscription.GroupName, event.EventAppeared.Event, workerID)

			err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
			if err != nil {
				s.log.Errorf("(EnrichSubscriber.when) err: {%v}", err)

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

func (s *EnrichSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
		return nil
	}

	acceptedEventTypes := []string{
		emailevent.EmailValidateV1,
	}
	if !utils.Contains(acceptedEventTypes, evt.GetEventType()) {
		return nil
	}

	ctx, span := tracing.StartProjectionTracerSpan(ctx, "EnrichSubscriber.When", evt)
	defer span.Finish()

	switch evt.GetEventType() {
	case emailevent.EmailValidateV1:
		_ = s.emailEventHandler.OnEmailValidate(ctx, evt)
	}
	return nil
}
