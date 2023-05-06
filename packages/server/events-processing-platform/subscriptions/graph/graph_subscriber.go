package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	contact_event_handlers "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event_handlers"
	contact_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	email_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	organization_event_handlers "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/event_handlers"
	organization_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	phone_number_event_handlers "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/event_handlers"
	phone_number_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	user_event_handlers "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/event_handlers"
	user_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"golang.org/x/sync/errgroup"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphSubscriber struct {
	log                      logger.Logger
	db                       *esdb.Client
	cfg                      *config.Config
	repositories             *repository.Repositories
	phoneNumberEventHandler  *phone_number_event_handlers.GraphPhoneNumberEventHandler
	contactEventHandler      *contact_event_handlers.GraphContactEventHandler
	organizationEventHandler *organization_event_handlers.GraphOrganizationEventHandler
	emailEventHandler        *GraphEmailEventHandler
	userEventHandler         *user_event_handlers.GraphUserEventHandler
}

func NewGraphSubscriber(log logger.Logger, db *esdb.Client, repositories *repository.Repositories, cfg *config.Config) *GraphSubscriber {
	return &GraphSubscriber{
		log:                      log,
		db:                       db,
		repositories:             repositories,
		cfg:                      cfg,
		contactEventHandler:      &contact_event_handlers.GraphContactEventHandler{Repositories: repositories},
		organizationEventHandler: &organization_event_handlers.GraphOrganizationEventHandler{Repositories: repositories},
		phoneNumberEventHandler:  &phone_number_event_handlers.GraphPhoneNumberEventHandler{Repositories: repositories},
		emailEventHandler:        &GraphEmailEventHandler{Repositories: repositories},
		userEventHandler:         &user_event_handlers.GraphUserEventHandler{Repositories: repositories},
	}
}

func (s *GraphSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.GraphSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.GraphSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{},
		)
		if err != nil {
			return err
		}
		defer sub.Close()

		group.Go(s.runWorker(ctx, worker, sub, i))
	}
	return group.Wait()
}

func (consumer *GraphSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *GraphSubscriber) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {

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
			s.log.ConsumedEvent(s.cfg.Subscriptions.GraphSubscription.GroupName, event.EventAppeared.Event, workerID)

			err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
			if err != nil {
				s.log.Errorf("(GraphSubscriber.when) err: {%v}", err)

				if err := stream.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
					s.log.Errorf("(stream.Nack) err: {%v}", err)
					return errors.Wrap(err, "stream.Nack")
				}
			}

			err = stream.Ack(event.EventAppeared.Event)
			if err != nil {
				s.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			s.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (s *GraphSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "GraphSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	if strings.HasPrefix(evt.GetAggregateID(), "$") {
		return nil
	}

	switch evt.GetEventType() {

	case phone_number_events.PhoneNumberCreatedV1:
		return s.phoneNumberEventHandler.OnPhoneNumberCreate(ctx, evt)
	case phone_number_events.PhoneNumberUpdatedV1:
		return s.phoneNumberEventHandler.OnPhoneNumberUpdate(ctx, evt)

	case email_events.EmailCreatedV1:
		return s.emailEventHandler.OnEmailCreate(ctx, evt)
	case email_events.EmailUpdatedV1:
		return s.emailEventHandler.OnEmailUpdate(ctx, evt)
	case email_events.EmailValidationFailedV1:
		return s.emailEventHandler.OnEmailValidationFailed(ctx, evt)
	case email_events.EmailValidatedV1:
		return s.emailEventHandler.OnEmailValidated(ctx, evt)

	case contact_events.ContactCreatedV1:
		return s.contactEventHandler.OnContactCreate(ctx, evt)
	case contact_events.ContactUpdatedV1:
		return s.contactEventHandler.OnContactUpdate(ctx, evt)
	case contact_events.ContactPhoneNumberLinkedV1:
		return s.contactEventHandler.OnPhoneNumberLinkedToContact(ctx, evt)
	case contact_events.ContactEmailLinkedV1:
		return s.contactEventHandler.OnEmailLinkedToContact(ctx, evt)

	case organization_events.OrganizationCreatedV1:
		return s.organizationEventHandler.OnOrganizationCreate(ctx, evt)
	case organization_events.OrganizationUpdatedV1:
		return s.organizationEventHandler.OnOrganizationUpdate(ctx, evt)
	case organization_events.OrganizationPhoneNumberLinkedV1:
		return s.organizationEventHandler.OnPhoneNumberLinkedToOrganization(ctx, evt)
	case organization_events.OrganizationEmailLinkedV1:
		return s.organizationEventHandler.OnEmailLinkedToOrganization(ctx, evt)

	case user_events.UserCreatedV1:
		return s.userEventHandler.OnUserCreate(ctx, evt)
	case user_events.UserUpdatedV1:
		return s.userEventHandler.OnUserUpdate(ctx, evt)
	case user_events.UserPhoneNumberLinkedV1:
		return s.userEventHandler.OnPhoneNumberLinkedToUser(ctx, evt)
	case user_events.UserEmailLinkedV1:
		return s.userEventHandler.OnEmailLinkedToUser(ctx, evt)
	default:
		s.log.Errorf("(GraphSubscriber) Unknown EventType: {%s}", evt.EventType)
		return eventstore.ErrInvalidEventType
	}
}
