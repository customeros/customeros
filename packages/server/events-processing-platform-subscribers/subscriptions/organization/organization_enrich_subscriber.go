package organization

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	aiConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	ai "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	orgevts "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"strings"
)

type OrganizationEnrichSubscriber struct {
	log                      logger.Logger
	db                       *esdb.Client
	cfg                      *config.Config
	organizationEventHandler *organizationEventHandler
}

func NewOrganizationEnrichSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, caches caches.Cache, grpcClients *grpc_client.Clients, services *service.Services) *OrganizationEnrichSubscriber {
	aiCfg := aiConfig.Config{
		OpenAi: aiConfig.AiModelConfigOpenAi{
			ApiKey:       cfg.Services.Ai.ApiKey,
			Organization: cfg.Services.OpenAi.Organization,
			Model:        "gpt-3.5-turbo-1106", // 1106 has an extra parameter available that locks response as JSON)
		},
		Anthropic: aiConfig.AiModelConfigAnthropic{
			ApiPath: cfg.Services.Ai.ApiPath,
			ApiKey:  cfg.Services.Ai.ApiKey,
		},
	}
	aiModel := ai.NewAiModel(ai.AnthropicModelType, aiCfg)
	return &OrganizationEnrichSubscriber{
		log:                      log,
		db:                       db,
		cfg:                      cfg,
		organizationEventHandler: NewOrganizationEventHandler(services, log, cfg, caches, aiModel, grpcClients),
	}
}

func (s *OrganizationEnrichSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
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

func (s *OrganizationEnrichSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *OrganizationEnrichSubscriber) ProcessEvents(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error {

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
				s.log.Errorf("(OrganizationEnrichSubscriber.when) err: {%v}", err)

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

func (s *OrganizationEnrichSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
		return nil
	}

	acceptedEventTypes := []string{
		orgevts.OrganizationLinkDomainV1,
		orgevts.OrganizationRequestEnrichV1,
	}

	if !utils.Contains(acceptedEventTypes, evt.GetEventType()) {
		return nil
	}

	ctx, span := tracing.StartProjectionTracerSpan(ctx, "OrganizationEnrichSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	switch evt.GetEventType() {
	case orgevts.OrganizationRequestEnrichV1:
		return s.organizationEventHandler.EnrichOrganizationByRequest(ctx, evt)
	default:
		return nil
	}
}
