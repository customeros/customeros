package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain"
	contactevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	emailevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	interactionevtevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/events"
	jobroleevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	locationevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	logentryevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/events"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	phonenumberevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	userevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
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
	organizationEventHandler *OrganizationEventHandler
	emailEventHandler        *GraphEmailEventHandler
	userEventHandler         *GraphUserEventHandler
	locationEventHandler     *GraphLocationEventHandler
	jobRoleEventHandler      *GraphJobRoleEventHandler
	interactionEventHandler  *GraphInteractionEventHandler
	logEntryEventHandler     *GraphLogEntryEventHandler
}

func NewGraphSubscriber(log logger.Logger, db *esdb.Client, repositories *repository.Repositories, commands *domain.Commands, cfg *config.Config) *GraphSubscriber {
	return &GraphSubscriber{
		log:                      log,
		db:                       db,
		repositories:             repositories,
		cfg:                      cfg,
		contactEventHandler:      &GraphContactEventHandler{Repositories: repositories},
		organizationEventHandler: &OrganizationEventHandler{log: log, repositories: repositories, organizationCommands: commands.OrganizationCommands},
		phoneNumberEventHandler:  &GraphPhoneNumberEventHandler{Repositories: repositories},
		emailEventHandler:        &GraphEmailEventHandler{Repositories: repositories},
		userEventHandler:         &GraphUserEventHandler{repositories: repositories, log: log},
		locationEventHandler:     &GraphLocationEventHandler{Repositories: repositories},
		jobRoleEventHandler:      &GraphJobRoleEventHandler{Repositories: repositories},
		interactionEventHandler:  &GraphInteractionEventHandler{Repositories: repositories, Log: log},
		logEntryEventHandler:     &GraphLogEntryEventHandler{Repositories: repositories, organizationCommands: commands.OrganizationCommands, log: log},
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

			if event.EventAppeared.Event.Event == nil {
				s.log.Errorf("(GraphSubscriber) event.EventAppeared.Event.Event is nil")
			} else {
				err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
				if err != nil {
					s.log.Errorf("(GraphSubscriber.when) err: {%v}", err)

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

func (s *GraphSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "GraphSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	if strings.HasPrefix(evt.GetAggregateID(), "$") {
		return nil
	}

	switch evt.GetEventType() {

	case phonenumberevents.PhoneNumberCreateV1:
		return s.phoneNumberEventHandler.OnPhoneNumberCreate(ctx, evt)
	case phonenumberevents.PhoneNumberUpdateV1:
		return s.phoneNumberEventHandler.OnPhoneNumberUpdate(ctx, evt)
	case phonenumberevents.PhoneNumberValidationFailedV1:
		return s.phoneNumberEventHandler.OnPhoneNumberValidationFailed(ctx, evt)
	case phonenumberevents.PhoneNumberValidationSkippedV1:
		return nil
	case phonenumberevents.PhoneNumberValidatedV1:
		return s.phoneNumberEventHandler.OnPhoneNumberValidated(ctx, evt)

	case emailevents.EmailCreateV1:
		return s.emailEventHandler.OnEmailCreate(ctx, evt)
	case emailevents.EmailUpdateV1:
		return s.emailEventHandler.OnEmailUpdate(ctx, evt)
	case emailevents.EmailValidationFailedV1:
		return s.emailEventHandler.OnEmailValidationFailed(ctx, evt)
	case emailevents.EmailValidatedV1:
		return s.emailEventHandler.OnEmailValidated(ctx, evt)

	case contactevents.ContactCreateV1:
		return s.contactEventHandler.OnContactCreate(ctx, evt)
	case contactevents.ContactUpdateV1:
		return s.contactEventHandler.OnContactUpdate(ctx, evt)
	case contactevents.ContactPhoneNumberLinkV1:
		return s.contactEventHandler.OnPhoneNumberLinkToContact(ctx, evt)
	case contactevents.ContactEmailLinkV1:
		return s.contactEventHandler.OnEmailLinkToContact(ctx, evt)

	case orgevents.OrganizationCreateV1:
		return s.organizationEventHandler.OnOrganizationCreate(ctx, evt)
	case orgevents.OrganizationUpdateV1:
		return s.organizationEventHandler.OnOrganizationUpdate(ctx, evt)
	case orgevents.OrganizationPhoneNumberLinkV1:
		return s.organizationEventHandler.OnPhoneNumberLinkedToOrganization(ctx, evt)
	case orgevents.OrganizationEmailLinkV1:
		return s.organizationEventHandler.OnEmailLinkedToOrganization(ctx, evt)
	case orgevents.OrganizationLocationLinkV1:
		return s.organizationEventHandler.OnLocationLinkedToOrganization(ctx, evt)
	case orgevents.OrganizationLinkDomainV1:
		return s.organizationEventHandler.OnDomainLinkedToOrganization(ctx, evt)
	case orgevents.OrganizationAddSocialV1:
		return s.organizationEventHandler.OnSocialAddedToOrganization(ctx, evt)
	case orgevents.OrganizationUpdateRenewalLikelihoodV1:
		return s.organizationEventHandler.OnRenewalLikelihoodUpdate(ctx, evt)
	case orgevents.OrganizationUpdateRenewalForecastV1:
		return s.organizationEventHandler.OnRenewalForecastUpdate(ctx, evt)
	case orgevents.OrganizationUpdateBillingDetailsV1:
		return s.organizationEventHandler.OnBillingDetailsUpdate(ctx, evt)
	case orgevents.OrganizationHideV1:
		return s.organizationEventHandler.OnOrganizationHide(ctx, evt)
	case orgevents.OrganizationShowV1:
		return s.organizationEventHandler.OnOrganizationShow(ctx, evt)
	case orgevents.OrganizationRefreshLastTouchpointV1:
		return s.organizationEventHandler.OnRefreshLastTouchpoint(ctx, evt)
	case orgevents.OrganizationUpsertCustomFieldV1:
		return s.organizationEventHandler.OnUpsertCustomField(ctx, evt)
	case orgevents.OrganizationRequestRenewalForecastV1,
		orgevents.OrganizationRequestNextCycleDateV1,
		orgevents.OrganizationRequestScrapeByWebsiteV1:
		return nil

	case userevents.UserCreateV1:
		return s.userEventHandler.OnUserCreate(ctx, evt)
	case userevents.UserUpdateV1:
		return s.userEventHandler.OnUserUpdate(ctx, evt)
	case userevents.UserPhoneNumberLinkV1:
		return s.userEventHandler.OnPhoneNumberLinkedToUser(ctx, evt)
	case userevents.UserEmailLinkV1:
		return s.userEventHandler.OnEmailLinkedToUser(ctx, evt)
	case userevents.UserJobRoleLinkV1:
		return s.userEventHandler.OnJobRoleLinkedToUser(ctx, evt)
	case userevents.UserAddPlayerV1:
		return s.userEventHandler.OnAddPlayer(ctx, evt)
	case userevents.UserAddRoleV1:
		return s.userEventHandler.OnAddRole(ctx, evt)
	case userevents.UserRemoveRoleV1:
		return s.userEventHandler.OnRemoveRole(ctx, evt)

	case locationevents.LocationCreateV1:
		return s.locationEventHandler.OnLocationCreate(ctx, evt)
	case locationevents.LocationUpdateV1:
		return s.locationEventHandler.OnLocationUpdate(ctx, evt)
	case locationevents.LocationValidationFailedV1:
		return s.locationEventHandler.OnLocationValidationFailed(ctx, evt)
	case locationevents.LocationValidationSkippedV1:
		return nil
	case locationevents.LocationValidatedV1:
		return s.locationEventHandler.OnLocationValidated(ctx, evt)
	case jobroleevents.JobRoleCreateV1:
		return s.jobRoleEventHandler.OnJobRoleCreate(ctx, evt)

	case interactionevtevents.InteractionEventRequestSummaryV1,
		interactionevtevents.InteractionEventRequestActionItemsV1:
		return nil
	case interactionevtevents.InteractionEventReplaceSummaryV1:
		return s.interactionEventHandler.OnSummaryReplace(ctx, evt)
	case interactionevtevents.InteractionEventReplaceActionItemsV1:
		return s.interactionEventHandler.OnActionItemsReplace(ctx, evt)

	case logentryevents.LogEntryCreateV1:
		return s.logEntryEventHandler.OnCreate(ctx, evt)
	case logentryevents.LogEntryUpdateV1:
		return s.logEntryEventHandler.OnUpdate(ctx, evt)
	case logentryevents.LogEntryAddTagV1:
		return s.logEntryEventHandler.OnAddTag(ctx, evt)
	case logentryevents.LogEntryRemoveTagV1:
		return s.logEntryEventHandler.OnRemoveTag(ctx, evt)
	default:
		s.log.Errorf("(GraphSubscriber) Unknown EventType: {%s}", evt.EventType)
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}
