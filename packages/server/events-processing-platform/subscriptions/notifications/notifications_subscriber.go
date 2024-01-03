package notifications

import (
	"context"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"golang.org/x/sync/errgroup"

	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type NotificationsSubscriber struct {
	log             logger.Logger
	db              *esdb.Client
	cfg             *config.Config
	orgEventHandler *OrganizationEventHandler
}

func NewNotificationsSubscriber(log logger.Logger, db *esdb.Client, repositories *repository.Repositories, grpcClients *grpc_client.Clients, cfg *config.Config) *NotificationsSubscriber {
	return &NotificationsSubscriber{
		log:             log,
		db:              db,
		cfg:             cfg,
		orgEventHandler: NewOrganizationEventHandler(log, repositories, cfg),
	}
}

func (s *NotificationsSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.NotificationsSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.NotificationsSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.NotificationsSubscription.BufferSizeClient,
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

func (consumer *NotificationsSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *NotificationsSubscriber) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {

	for {
		event := stream.Recv()
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
			s.log.EventAppeared(s.cfg.Subscriptions.NotificationsSubscription.GroupName, event.EventAppeared.Event, workerID)

			if event.EventAppeared.Event.Event == nil {
				s.log.Errorf("(NotificationsSubscriber) event.EventAppeared.Event.Event is nil")
			} else {
				err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
				if err != nil {
					s.log.Errorf("(NotificationSubscriber.when) err: {%v}", err)

					if err := stream.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
						s.log.Errorf("(stream.Nack) err: {%v}", err)
						return errors.Wrap(err, "stream.Nack")
					}
				}
			}

			err := stream.Ack(event.EventAppeared.Event)
			if err != nil {
				s.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			s.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (s *NotificationsSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "NotificationSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	if strings.HasPrefix(evt.GetAggregateID(), "$") {
		return nil
	}

	if s.cfg.Subscriptions.NotificationsSubscription.IgnoreEvents {
		return nil
	}

	switch evt.GetEventType() {

	case orgevents.OrganizationUpdateOwnerV1:
		return s.orgEventHandler.OnOrganizationUpdateOwner(ctx, evt)

	default:
		return nil
	}
}
