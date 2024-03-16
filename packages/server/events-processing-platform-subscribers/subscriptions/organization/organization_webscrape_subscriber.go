package organization

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/grpc_client"
	"strings"

	aiConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	ai "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	orgevts "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"golang.org/x/sync/errgroup"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OrganizationWebscrapeSubscriber struct {
	log                      logger.Logger
	db                       *esdb.Client
	cfg                      *config.Config
	organizationEventHandler *organizationEventHandler
}

func NewOrganizationWebscrapeSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, repositories *repository.Repositories, caches caches.Cache, grpcClients *grpc_client.Clients) *OrganizationWebscrapeSubscriber {
	aiCfg := aiConfig.Config{
		OpenAi: aiConfig.AiModelConfigOpenAi{
			ApiKey:       cfg.Services.OpenAi.ApiKey,
			Organization: cfg.Services.OpenAi.Organization,
			Model:        "gpt-3.5-turbo-1106", // 1106 has an extra parameter available that locks response as JSON)
		},
		Anthropic: aiConfig.AiModelConfigAnthropic{
			ApiPath: cfg.Services.Anthropic.ApiPath,
			ApiKey:  cfg.Services.Anthropic.ApiKey,
		},
	}
	domainScraper := NewDomainScraper(log, cfg, repositories, ai.NewAiModel(ai.OpenAiModelType, aiCfg))
	aiModel := ai.NewAiModel(ai.AnthropicModelType, aiCfg)
	return &OrganizationWebscrapeSubscriber{
		log:                      log,
		db:                       db,
		cfg:                      cfg,
		organizationEventHandler: NewOrganizationEventHandler(repositories, log, cfg, caches, domainScraper, aiModel, grpcClients),
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

	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
		return nil
	}

	switch evt.GetEventType() {
	case orgevts.OrganizationLinkDomainV1:
		return s.organizationEventHandler.WebScrapeOrganizationByDomain(ctx, evt)
	case orgevts.OrganizationRequestScrapeByWebsiteV1:
		return s.organizationEventHandler.WebScrapeOrganizationByWebsite(ctx, evt)
	default:
		return nil
	}
}
