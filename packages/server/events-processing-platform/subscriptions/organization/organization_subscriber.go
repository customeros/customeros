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

type OrganizationSubscriber struct {
	log                      logger.Logger
	db                       *esdb.Client
	cfg                      *config.Config
	organizationEventHandler *organizationEventHandler
}

func NewOrganizationSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, orgCommands *command_handler.OrganizationCommands, repositories *repository.Repositories, caches caches.Cache) *OrganizationSubscriber {
	return &OrganizationSubscriber{
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

func (s *OrganizationSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.OrganizationSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.OrganizationSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.OrganizationSubscription.BufferSizeClient,
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

func (s *OrganizationSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *OrganizationSubscriber) ProcessEvents(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error {

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
			s.log.EventAppeared(s.cfg.Subscriptions.OrganizationSubscription.GroupName, event.EventAppeared.Event, workerID)

			err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
			if err != nil {
				s.log.Errorf("(OrganizationSubscriber.when) err: {%v}", err)

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

func (s *OrganizationSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "OrganizationSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	if strings.HasPrefix(evt.GetAggregateID(), "$") {
		return nil
	}

	switch evt.GetEventType() {
	case orgevts.OrganizationCreateV1:
		return s.organizationEventHandler.AdjustNewOrganizationFields(ctx, evt)
	case orgevts.OrganizationUpdateV1:
		return s.organizationEventHandler.AdjustUpdatedOrganizationFields(ctx, evt)
	case orgevts.OrganizationLinkDomainV1:
		return s.organizationEventHandler.WebscrapeOrganizationByDomain(ctx, evt)
	case orgevts.OrganizationRequestRenewalForecastV1:
		return s.organizationEventHandler.OnRenewalForecastRequested(ctx, evt)
	case orgevts.OrganizationRequestNextCycleDateV1:
		return s.organizationEventHandler.OnNextCycleDateRequested(ctx, evt)
	case orgevts.OrganizationRequestScrapeByWebsiteV1:
		return s.organizationEventHandler.WebscrapeOrganizationByWebsite(ctx, evt)
	case
		orgevts.OrganizationPhoneNumberLinkV1,
		orgevts.OrganizationEmailLinkV1,
		orgevts.OrganizationAddSocialV1,
		orgevts.OrganizationUpdateRenewalLikelihoodV1,
		orgevts.OrganizationUpdateRenewalForecastV1,
		orgevts.OrganizationUpdateBillingDetailsV1:
		return nil

	default:
		s.log.Warnf("(OrganizationSubscriber) Unknown EventType: {%s}", evt.EventType)
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}
