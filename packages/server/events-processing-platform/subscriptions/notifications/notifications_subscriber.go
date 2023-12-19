package notifications

import (
	"context"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	commentevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/event"
	contactevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	contractevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	emailevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	ieevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	issueevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	jobroleevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	locationevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	logentryevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	phonenumberevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	servicelineitemevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	userevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"golang.org/x/sync/errgroup"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type NotificationsSubscriber struct {
	log              logger.Logger
	db               *esdb.Client
	cfg              *config.Config
	userEventHandler *UserEventHandler
}

func NewNotificationsSubscriber(log logger.Logger, db *esdb.Client, repositories *repository.Repositories, grpcClients *grpc_client.Clients, cfg *config.Config) *NotificationsSubscriber {
	return &NotificationsSubscriber{
		log:              log,
		db:               db,
		cfg:              cfg,
		userEventHandler: NewUserEventHandler(log, repositories, cfg),
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

	switch evt.GetEventType() {

	case phonenumberevents.PhoneNumberCreateV1,
		phonenumberevents.PhoneNumberUpdateV1,
		phonenumberevents.PhoneNumberValidationFailedV1,
		phonenumberevents.PhoneNumberValidationSkippedV1,
		phonenumberevents.PhoneNumberValidatedV1:
		return nil

	case emailevents.EmailCreateV1,
		emailevents.EmailUpdateV1,
		emailevents.EmailValidationFailedV1,
		emailevents.EmailValidatedV1:
		return nil

	case contactevent.ContactCreateV1,
		contactevent.ContactUpdateV1,
		contactevent.ContactPhoneNumberLinkV1,
		contactevent.ContactEmailLinkV1,
		contactevent.ContactLocationLinkV1,
		contactevent.ContactOrganizationLinkV1:
		return nil

	case orgevents.OrganizationCreateV1:
		return nil
	case orgevents.OrganizationUpdateV1:
		return nil
	case orgevents.OrganizationPhoneNumberLinkV1,
		orgevents.OrganizationEmailLinkV1,
		orgevents.OrganizationLocationLinkV1,
		orgevents.OrganizationLinkDomainV1,
		orgevents.OrganizationAddSocialV1,
		orgevents.OrganizationHideV1,
		orgevents.OrganizationShowV1,
		orgevents.OrganizationRefreshLastTouchpointV1,
		orgevents.OrganizationRefreshArrV1,
		orgevents.OrganizationRefreshRenewalSummaryV1,
		orgevents.OrganizationUpsertCustomFieldV1,
		orgevents.OrganizationAddParentV1,
		orgevents.OrganizationRemoveParentV1,
		orgevents.OrganizationUpdateOnboardingStatusV1,
		orgevents.OrganizationRequestRenewalForecastV1,
		orgevents.OrganizationRequestNextCycleDateV1,
		orgevents.OrganizationUpdateRenewalLikelihoodV1,
		orgevents.OrganizationUpdateRenewalForecastV1,
		orgevents.OrganizationUpdateBillingDetailsV1,
		orgevents.OrganizationRequestScrapeByWebsiteV1:
		return nil

	case userevents.UserUpdateV1:
		return s.userEventHandler.OnUserUpdate(ctx, evt)
	case userevents.UserJobRoleLinkV1:
		return s.userEventHandler.OnJobRoleLinkedToUser(ctx, evt)
	case userevents.UserAddRoleV1:
		return s.userEventHandler.OnAddRole(ctx, evt)
	case userevents.UserCreateV1,
		userevents.UserPhoneNumberLinkV1,
		userevents.UserEmailLinkV1,
		userevents.UserAddPlayerV1,
		userevents.UserRemoveRoleV1:
		return nil

	case locationevents.LocationCreateV1,
		locationevents.LocationUpdateV1,
		locationevents.LocationValidationFailedV1,
		locationevents.LocationValidationSkippedV1,
		locationevents.LocationValidatedV1,
		jobroleevents.JobRoleCreateV1:
		return nil

	case ieevent.InteractionEventRequestSummaryV1,
		ieevent.InteractionEventRequestActionItemsV1,
		ieevent.InteractionEventReplaceSummaryV1,
		ieevent.InteractionEventReplaceActionItemsV1,
		ieevent.InteractionEventCreateV1,
		ieevent.InteractionEventUpdateV1:
		return nil

	case logentryevents.LogEntryCreateV1,
		logentryevents.LogEntryUpdateV1,
		logentryevents.LogEntryAddTagV1,
		logentryevents.LogEntryRemoveTagV1:
		return nil

	case commentevent.CommentCreateV1,
		commentevent.CommentUpdateV1:
		return nil

	case issueevent.IssueCreateV1,
		issueevent.IssueUpdateV1,
		issueevent.IssueAddUserAssigneeV1,
		issueevent.IssueRemoveUserAssigneeV1,
		issueevent.IssueAddUserFollowerV1,
		issueevent.IssueRemoveUserFollowerV1:
		return nil

	case opportunityevent.OpportunityCreateV1,
		opportunityevent.OpportunityUpdateNextCycleDateV1,
		opportunityevent.OpportunityUpdateV1,
		opportunityevent.OpportunityCreateRenewalV1,
		opportunityevent.OpportunityUpdateRenewalV1,
		opportunityevent.OpportunityCloseWinV1,
		opportunityevent.OpportunityCloseLooseV1:
		return nil

	case contractevent.ContractCreateV1,
		contractevent.ContractUpdateV1,
		contractevent.ContractRolloutRenewalOpportunityV1,
		contractevent.ContractUpdateStatusV1:
		return nil

	case servicelineitemevent.ServiceLineItemCreateV1,
		servicelineitemevent.ServiceLineItemUpdateV1,
		servicelineitemevent.ServiceLineItemDeleteV1,
		servicelineitemevent.ServiceLineItemCloseV1:
		return nil

	default:
		s.log.Errorf("(NotificationSubscriber) Unknown EventType: {%s}", evt.EventType)
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}
