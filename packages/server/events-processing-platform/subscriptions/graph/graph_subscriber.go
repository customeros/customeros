package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	contact_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	email_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	job_role_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	location_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	organization_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	phone_number_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
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
	phoneNumberEventHandler  *GraphPhoneNumberEventHandler
	contactEventHandler      *GraphContactEventHandler
	organizationEventHandler *GraphOrganizationEventHandler
	emailEventHandler        *GraphEmailEventHandler
	userEventHandler         *GraphUserEventHandler
	locationEventHandler     *GraphLocationEventHandler
	jobRoleEventHandler      *GraphJobRoleEventHandler
}

func NewGraphSubscriber(log logger.Logger, db *esdb.Client, repositories *repository.Repositories, cfg *config.Config) *GraphSubscriber {
	return &GraphSubscriber{
		log:                      log,
		db:                       db,
		repositories:             repositories,
		cfg:                      cfg,
		contactEventHandler:      &GraphContactEventHandler{Repositories: repositories},
		organizationEventHandler: &GraphOrganizationEventHandler{Repositories: repositories},
		phoneNumberEventHandler:  &GraphPhoneNumberEventHandler{Repositories: repositories},
		emailEventHandler:        &GraphEmailEventHandler{Repositories: repositories},
		userEventHandler:         &GraphUserEventHandler{Repositories: repositories},
		locationEventHandler:     &GraphLocationEventHandler{Repositories: repositories},
		jobRoleEventHandler:      &GraphJobRoleEventHandler{Repositories: repositories},
	}
}

func (s *GraphSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.GraphSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.GraphSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.GraphSubscription.BufferSizeClient,
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
			s.log.EventAppeared(s.cfg.Subscriptions.GraphSubscription.GroupName, event.EventAppeared.Event, workerID)

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

	case
		phone_number_events.PhoneNumberCreateV1,
		phone_number_events.PhoneNumberCreateV1Legacy:
		return s.phoneNumberEventHandler.OnPhoneNumberCreate(ctx, evt)
	case
		phone_number_events.PhoneNumberUpdateV1,
		phone_number_events.PhoneNumberUpdateV1Legacy:
		return s.phoneNumberEventHandler.OnPhoneNumberUpdate(ctx, evt)
	case phone_number_events.PhoneNumberValidationFailedV1:
		return s.phoneNumberEventHandler.OnPhoneNumberValidationFailed(ctx, evt)
	case phone_number_events.PhoneNumberValidationSkippedV1:
		return nil
	case phone_number_events.PhoneNumberValidatedV1:
		return s.phoneNumberEventHandler.OnPhoneNumberValidated(ctx, evt)

	case
		email_events.EmailCreateV1,
		email_events.EmailCreateV1Legacy:
		return s.emailEventHandler.OnEmailCreate(ctx, evt)
	case
		email_events.EmailUpdateV1,
		email_events.EmailUpdateV1Legacy:
		return s.emailEventHandler.OnEmailUpdate(ctx, evt)
	case email_events.EmailValidationFailedV1:
		return s.emailEventHandler.OnEmailValidationFailed(ctx, evt)
	case email_events.EmailValidatedV1:
		return s.emailEventHandler.OnEmailValidated(ctx, evt)

	case contact_events.ContactCreateV1:
		return s.contactEventHandler.OnContactCreate(ctx, evt)
	case contact_events.ContactUpdateV1:
		return s.contactEventHandler.OnContactUpdate(ctx, evt)
	case contact_events.ContactPhoneNumberLinkV1:
		return s.contactEventHandler.OnPhoneNumberLinkToContact(ctx, evt)
	case contact_events.ContactEmailLinkV1:
		return s.contactEventHandler.OnEmailLinkToContact(ctx, evt)

	case organization_events.OrganizationCreateV1:
		return s.organizationEventHandler.OnOrganizationCreate(ctx, evt)
	case organization_events.OrganizationUpdateV1:
		return s.organizationEventHandler.OnOrganizationUpdate(ctx, evt)
	case organization_events.OrganizationPhoneNumberLinkV1:
		return s.organizationEventHandler.OnPhoneNumberLinkedToOrganization(ctx, evt)
	case organization_events.OrganizationEmailLinkV1:
		return s.organizationEventHandler.OnEmailLinkedToOrganization(ctx, evt)
	case organization_events.OrganizationLinkDomainV1:
		return s.organizationEventHandler.OnDomainLinkedToOrganization(ctx, evt)
	case organization_events.OrganizationAddSocialV1:
		return s.organizationEventHandler.OnSocialAddedToOrganization(ctx, evt)

	case user_events.UserCreateV1:
		return s.userEventHandler.OnUserCreate(ctx, evt)
	case user_events.UserUpdateV1:
		return s.userEventHandler.OnUserUpdate(ctx, evt)
	case user_events.UserPhoneNumberLinkV1:
		return s.userEventHandler.OnPhoneNumberLinkedToUser(ctx, evt)
	case user_events.UserEmailLinkV1:
		return s.userEventHandler.OnEmailLinkedToUser(ctx, evt)

	case
		location_events.LocationCreateV1Legacy,
		location_events.LocationCreateV1:
		return s.locationEventHandler.OnLocationCreate(ctx, evt)
	case
		location_events.LocationUpdateV1Legacy,
		location_events.LocationUpdateV1:
		return s.locationEventHandler.OnLocationUpdate(ctx, evt)
	case location_events.LocationValidationFailedV1:
		return s.locationEventHandler.OnLocationValidationFailed(ctx, evt)
	case location_events.LocationValidationSkippedV1:
		return nil
	case location_events.LocationValidatedV1:
		return s.locationEventHandler.OnLocationValidated(ctx, evt)
	case job_role_events.JobRoleCreateV1:
		return s.jobRoleEventHandler.OnJobRoleCreate(ctx, evt)
	case user_events.UserJobRoleLinkV1:
		return s.userEventHandler.OnJobRoleLinkedToUser(ctx, evt)
	default:
		s.log.Errorf("(GraphSubscriber) Unknown EventType: {%s}", evt.EventType)
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}
