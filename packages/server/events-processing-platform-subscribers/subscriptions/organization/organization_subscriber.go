package organization

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	orgevts "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	aiConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	ai "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type OrganizationSubscriber struct {
	log                      logger.Logger
	db                       *esdb.Client
	cfg                      *config.Config
	organizationEventHandler *organizationEventHandler
}

func NewOrganizationSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, repositories *repository.Repositories, caches caches.Cache, grpcClients *grpc_client.Clients) *OrganizationSubscriber {
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
	return &OrganizationSubscriber{
		log:                      log,
		db:                       db,
		cfg:                      cfg,
		organizationEventHandler: NewOrganizationEventHandler(repositories, log, cfg, caches, aiModel, grpcClients, nil),
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

	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
		return nil
	}

	switch evt.GetEventType() {
	case orgevts.OrganizationCreateV1:
		_ = s.organizationEventHandler.AdjustNewOrganizationFields(ctx, evt)
		return nil
	case orgevts.OrganizationUpdateV1:
		_ = s.organizationEventHandler.AdjustUpdatedOrganizationFields(ctx, evt)
		return nil
	case orgevts.OrganizationAdjustIndustryV1:
		_ = s.organizationEventHandler.OnAdjustIndustry(ctx, evt)
		return nil
	default:
		return nil
	}
}
