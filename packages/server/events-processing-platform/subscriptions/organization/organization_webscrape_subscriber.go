package organization

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	orgevts "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
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

type OrganizationWebscrapeSubscriber struct {
	log                      logger.Logger
	db                       *esdb.Client
	cfg                      *config.Config
	organizationEventHandler *organizationEventHandler
}

func NewOrganizationWebscrapeSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, orgCommands *command_handler.CommandHandlers, repositories *repository.Repositories, caches caches.Cache) *OrganizationWebscrapeSubscriber {
	return &OrganizationWebscrapeSubscriber{
		log: log,
		db:  db,
		cfg: cfg,
		organizationEventHandler: &organizationEventHandler{
			log:                  log,
			cfg:                  cfg,
			organizationCommands: orgCommands,
			repositories:         repositories,
			caches:               caches,
			domainScraper:        NewDomainScraper(log, cfg, repositories),
		},
	}
}

func (s *OrganizationWebscrapeSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.OrganizationWebscrapeSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.OrganizationWebscrapeSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.OrganizationWebscrapeSubscription.BufferSizeClient,
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

func (s *OrganizationWebscrapeSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *OrganizationWebscrapeSubscriber) ProcessEvents(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error {

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
			s.log.EventAppeared(s.cfg.Subscriptions.OrganizationWebscrapeSubscription.GroupName, event.EventAppeared.Event, workerID)

			err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
			if err != nil {
				s.log.Errorf("(OrganizationWebscrapeSubscriber.when) err: {%v}", err)

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

func (s *OrganizationWebscrapeSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "OrganizationWebscrapeSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	if strings.HasPrefix(evt.GetAggregateID(), "$") {
		return nil
	}

	switch evt.GetEventType() {
	case orgevts.OrganizationLinkDomainV1:
		return s.organizationEventHandler.WebscrapeOrganizationByDomain(ctx, evt)
	case orgevts.OrganizationRequestScrapeByWebsiteV1:
		return s.organizationEventHandler.WebscrapeOrganizationByWebsite(ctx, evt)

	default:
		return nil
	}
}
