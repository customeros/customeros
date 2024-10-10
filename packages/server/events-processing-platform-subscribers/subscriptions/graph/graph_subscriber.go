package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	contactevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	emailevents "github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/generic"
	reminderevents "github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder/event"
	"github.com/opentracing/opentracing-go"
	"strings"
	"time"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	commentevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment"
	contractevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	invoiceevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	issueevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	jobroleevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	locationevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	logentryevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	masterplanevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	orgplanevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	phonenumberevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	servicelineitemevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	tenantevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	userevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
)

type GraphSubscriber struct {
	log                          logger.Logger
	db                           *esdb.Client
	cfg                          *config.Config
	genericEventHandler          *GenericEventHandler
	phoneNumberEventHandler      *PhoneNumberEventHandler
	contactEventHandler          *ContactEventHandler
	organizationEventHandler     *OrganizationEventHandler
	emailEventHandler            *EmailEventHandler
	userEventHandler             *UserEventHandler
	locationEventHandler         *LocationEventHandler
	jobRoleEventHandler          *JobRoleEventHandler
	logEntryEventHandler         *LogEntryEventHandler
	issueEventHandler            *IssueEventHandler
	commentEventHandler          *CommentEventHandler
	opportunityEventHandler      *OpportunityEventHandler
	contractEventHandler         *ContractEventHandler
	serviceLineItemEventHandler  *ServiceLineItemEventHandler
	masterPlanEventHandler       *MasterPlanEventHandler
	invoiceEventHandler          *InvoiceEventHandler
	tenantEventHandler           *TenantEventHandler
	organizationPlanEventHandler *OrganizationPlanEventHandler
	bankAccountEventHandler      *BankAccountEventHandler
	reminderEventHandler         *ReminderEventHandler
}

func NewGraphSubscriber(log logger.Logger, db *esdb.Client, services *service.Services, grpcClients *grpc_client.Clients, cfg *config.Config, cache caches.Cache) *GraphSubscriber {
	return &GraphSubscriber{
		log:                          log,
		db:                           db,
		cfg:                          cfg,
		genericEventHandler:          NewGenericEventHandler(log, services, grpcClients),
		contactEventHandler:          NewContactEventHandler(log, services, grpcClients),
		organizationEventHandler:     NewOrganizationEventHandler(log, services, grpcClients, cache),
		phoneNumberEventHandler:      NewPhoneNumberEventHandler(log, services, grpcClients),
		emailEventHandler:            NewEmailEventHandler(log, services, grpcClients),
		userEventHandler:             NewUserEventHandler(log, services),
		locationEventHandler:         NewLocationEventHandler(services),
		jobRoleEventHandler:          NewJobRoleEventHandler(services),
		logEntryEventHandler:         NewLogEntryEventHandler(log, services, grpcClients),
		issueEventHandler:            NewIssueEventHandler(log, services, grpcClients),
		commentEventHandler:          NewCommentEventHandler(log, services),
		opportunityEventHandler:      NewOpportunityEventHandler(log, services, grpcClients),
		contractEventHandler:         NewContractEventHandler(log, services, grpcClients),
		serviceLineItemEventHandler:  NewServiceLineItemEventHandler(log, services, grpcClients),
		masterPlanEventHandler:       NewMasterPlanEventHandler(log, services),
		invoiceEventHandler:          NewInvoiceEventHandler(log, services, grpcClients),
		tenantEventHandler:           NewTenantEventHandler(log, services),
		organizationPlanEventHandler: NewOrganizationPlanEventHandler(log, services),
		bankAccountEventHandler:      NewBankAccountEventHandler(log, services),
		reminderEventHandler:         NewReminderEventHandler(log, services),
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

func (s *GraphSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
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
			span, _ := opentracing.StartSpanFromContext(ctx, "GraphSubscriber.ProcessEvents")
			defer span.Finish()
			wrappedErr := errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
			tracing.TraceErr(span, wrappedErr)
			s.log.Errorf(wrappedErr.Error())
			return wrappedErr
		}

		if event.EventAppeared != nil {
			s.log.EventAppeared(s.cfg.Subscriptions.GraphSubscription.GroupName, event.EventAppeared.Event, workerID)

			if event.EventAppeared.Event.Event == nil {
				span, _ := opentracing.StartSpanFromContext(ctx, "GraphSubscriber.ProcessEvents")
				defer span.Finish()
				err := errors.Wrap(errors.New("event.EventAppeared.Event.Event is nil"), "GraphSubscriber")
				tracing.TraceErr(span, err)
				s.log.Errorf(err.Error())
			} else {
				err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
				if err != nil {
					span, _ := opentracing.StartSpanFromContext(ctx, "GraphSubscriber.ProcessEvents")
					defer span.Finish()
					tracing.TraceErr(span, err)
					s.log.Errorf("(GraphSubscriber.when) err: {%v}", err)

					if err := stream.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("(stream.Nack) err: {%v}", err)
						return errors.Wrap(err, "stream.Nack")
					}
				}
			}

			err := stream.Ack(event.EventAppeared.Event)
			if err != nil {
				span, _ := opentracing.StartSpanFromContext(ctx, "GraphSubscriber.ProcessEvents")
				defer span.Finish()
				tracing.TraceErr(span, err)
				s.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			s.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (s *GraphSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
		return nil
	}
	switch evt.GetEventType() {
	case "V1_EVENT_COMPLETED",
		phonenumberevents.PhoneNumberValidateV1,
		emailevents.EmailValidationFailedV1,
		emailevents.EmailValidatedV1,
		emailevents.EmailValidateV1,
		emailevents.EmailUpsertV1,
		orgevents.OrganizationRefreshLastTouchpointV1,
		phonenumberevents.PhoneNumberValidationSkippedV1,
		contactevent.ContactRequestEnrichV1,
		orgevents.OrganizationRequestRenewalForecastV1,
		orgevents.OrganizationRequestNextCycleDateV1,
		orgevents.OrganizationUpdateRenewalLikelihoodV1,
		orgevents.OrganizationUpdateRenewalForecastV1,
		orgevents.OrganizationUpdateBillingDetailsV1,
		orgevents.OrganizationRequestScrapeByWebsiteV1,
		orgevents.OrganizationAdjustIndustryV1,
		orgevents.OrganizationRequestEnrichV1,
		contractevent.ContractUpdateStatusV1,
		locationevents.LocationValidationSkippedV1,
		reminderevents.ReminderNotificationV1,
		orgevents.OrganizationUpdateOwnerNotificationV1,
		invoiceevents.InvoicePdfRequestedV1,
		invoiceevents.InvoicePaidV1,
		invoiceevents.InvoiceFillRequestedV1,
		invoiceevents.InvoicePayNotificationV1,
		invoiceevents.InvoicePayV1,
		invoiceevents.InvoiceRemindNotificationV1,
		emailevents.EmailCreateV1,
		emailevents.EmailUpdateV1:
		return nil
	}

	ctx, span := tracing.StartProjectionTracerSpan(ctx, "GraphSubscriber.When", evt)
	defer span.Finish()

	// set 25 sec context deadline
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	switch evt.GetEventType() {
	case generic.LinkEntityWithEntityV1:
		_ = s.genericEventHandler.OnLinkEntityWithEntityV1(ctx, evt)
		return nil

	case phonenumberevents.PhoneNumberCreateV1:
		_ = s.phoneNumberEventHandler.OnPhoneNumberCreate(ctx, evt)
		return nil
	case phonenumberevents.PhoneNumberUpdateV1:
		_ = s.phoneNumberEventHandler.OnPhoneNumberUpdate(ctx, evt)
		return nil
	case phonenumberevents.PhoneNumberValidationFailedV1:
		_ = s.phoneNumberEventHandler.OnPhoneNumberValidationFailed(ctx, evt)
		return nil
	case phonenumberevents.PhoneNumberValidatedV1:
		_ = s.phoneNumberEventHandler.OnPhoneNumberValidated(ctx, evt)
		return nil

	case emailevents.EmailValidatedV2:
		_ = s.emailEventHandler.OnEmailValidatedV2(ctx, evt)
		return nil
	case emailevents.EmailDeleteV1:
		_ = s.emailEventHandler.OnEmailDelete(ctx, evt)
		return nil

	case contactevent.ContactCreateV1:
		_ = s.contactEventHandler.OnContactCreate(ctx, evt)
		return nil
	case contactevent.ContactUpdateV1:
		_ = s.contactEventHandler.OnContactUpdate(ctx, evt)
		return nil
	case contactevent.ContactPhoneNumberLinkV1:
		_ = s.contactEventHandler.OnPhoneNumberLinkToContact(ctx, evt)
		return nil
	case contactevent.ContactEmailLinkV1:
		_ = s.contactEventHandler.OnEmailLinkToContact(ctx, evt)
		return nil
	case contactevent.ContactEmailUnlinkV1:
		_ = s.contactEventHandler.OnEmailUnlinkFromContact(ctx, evt)
		return nil
	case contactevent.ContactLocationLinkV1:
		_ = s.contactEventHandler.OnLocationLinkToContact(ctx, evt)
		return nil
	case contactevent.ContactOrganizationLinkV1:
		_ = s.contactEventHandler.OnContactLinkToOrganization(ctx, evt)
		return nil
	case contactevent.ContactAddSocialV1:
		_ = s.contactEventHandler.OnSocialAddedToContactV1(ctx, evt)
		return nil
	case contactevent.ContactRemoveSocialV1:
		_ = s.contactEventHandler.OnSocialRemovedFromContactV1(ctx, evt)
		return nil
	case contactevent.ContactAddTagV1:
		_ = s.contactEventHandler.OnAddTag(ctx, evt)
		return nil
	case contactevent.ContactRemoveTagV1:
		_ = s.contactEventHandler.OnRemoveTag(ctx, evt)
		return nil
	case contactevent.ContactAddLocationV1:
		_ = s.contactEventHandler.OnLocationAddedToContact(ctx, evt)
		return nil
	case contactevent.ContactHideV1:
		_ = s.contactEventHandler.OnContactHide(ctx, evt)
		return nil
	case contactevent.ContactShowV1:
		_ = s.contactEventHandler.OnContactShow(ctx, evt)
		return nil

	case orgevents.OrganizationCreateV1:
		_ = s.organizationEventHandler.OnOrganizationCreate(ctx, evt)
		return nil
	case orgevents.OrganizationUpdateV1:
		_ = s.organizationEventHandler.OnOrganizationUpdate(ctx, evt)
		return nil
	case orgevents.OrganizationPhoneNumberLinkV1:
		_ = s.organizationEventHandler.OnPhoneNumberLinkedToOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationEmailLinkV1:
		_ = s.organizationEventHandler.OnEmailLinkToOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationEmailUnlinkV1:
		_ = s.organizationEventHandler.OnEmailUnlinkFromOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationLocationLinkV1:
		_ = s.organizationEventHandler.OnLocationLinkedToOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationLinkDomainV1:
		_ = s.organizationEventHandler.OnDomainLinkedToOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationUnlinkDomainV1:
		_ = s.organizationEventHandler.OnDomainUnlinkedFromOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationAddSocialV1:
		_ = s.organizationEventHandler.OnSocialAddedToOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationRemoveSocialV1:
		_ = s.organizationEventHandler.OnSocialRemovedFromOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationHideV1:
		_ = s.organizationEventHandler.OnOrganizationHide(ctx, evt)
		return nil
	case orgevents.OrganizationShowV1:
		_ = s.organizationEventHandler.OnOrganizationShow(ctx, evt)
		return nil
	case orgevents.OrganizationRefreshArrV1:
		_ = s.organizationEventHandler.OnRefreshArr(ctx, evt)
		return nil
	case orgevents.OrganizationRefreshRenewalSummaryV1:
		_ = s.organizationEventHandler.OnRefreshRenewalSummaryV1(ctx, evt)
		return nil
	case orgevents.OrganizationRefreshDerivedDataV1:
		_ = s.organizationEventHandler.OnRefreshDerivedDataV1(ctx, evt)
		return nil
	case orgevents.OrganizationUpsertCustomFieldV1:
		_ = s.organizationEventHandler.OnUpsertCustomField(ctx, evt)
		return nil
	case orgevents.OrganizationAddParentV1:
		_ = s.organizationEventHandler.OnLinkWithParentOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationRemoveParentV1:
		_ = s.organizationEventHandler.OnUnlinkFromParentOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationUpdateOnboardingStatusV1:
		_ = s.organizationEventHandler.OnUpdateOnboardingStatus(ctx, evt)
		return nil
	case orgevents.OrganizationUpdateOwnerV1:
		_ = s.organizationEventHandler.OnUpdateOwner(ctx, evt)
		return nil
	case orgevents.OrganizationAddTagV1:
		_ = s.organizationEventHandler.OnAddTag(ctx, evt)
		return nil
	case orgevents.OrganizationRemoveTagV1:
		_ = s.organizationEventHandler.OnRemoveTag(ctx, evt)
		return nil
	case orgevents.OrganizationCreateBillingProfileV1:
		_ = s.organizationEventHandler.OnCreateBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationUpdateBillingProfileV1:
		_ = s.organizationEventHandler.OnUpdateBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationEmailLinkToBillingProfileV1:
		_ = s.organizationEventHandler.OnEmailLinkedToBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationEmailUnlinkFromBillingProfileV1:
		_ = s.organizationEventHandler.OnEmailUnlinkedFromBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationLocationLinkToBillingProfileV1:
		_ = s.organizationEventHandler.OnLocationLinkedToBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationLocationUnlinkFromBillingProfileV1:
		_ = s.organizationEventHandler.OnLocationUnlinkedFromBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationAddLocationV1:
		_ = s.organizationEventHandler.OnLocationAddedToOrganization(ctx, evt)
		return nil

	case userevents.UserCreateV1:
		_ = s.userEventHandler.OnUserCreate(ctx, evt)
		return nil
	case userevents.UserUpdateV1:
		_ = s.userEventHandler.OnUserUpdate(ctx, evt)
		return nil
	case userevents.UserPhoneNumberLinkV1:
		_ = s.userEventHandler.OnPhoneNumberLinkedToUser(ctx, evt)
		return nil
	case userevents.UserEmailLinkV1:
		_ = s.userEventHandler.OnEmailLinkedToUser(ctx, evt)
		return nil
	case userevents.UserEmailUnlinkV1:
		_ = s.userEventHandler.OnEmailUnlinkedFromUser(ctx, evt)
		return nil
	case userevents.UserJobRoleLinkV1:
		_ = s.userEventHandler.OnJobRoleLinkedToUser(ctx, evt)
		return nil
	case userevents.UserAddRoleV1:
		_ = s.userEventHandler.OnAddRole(ctx, evt)
		return nil
	case userevents.UserRemoveRoleV1:
		_ = s.userEventHandler.OnRemoveRole(ctx, evt)
		return nil

	case locationevents.LocationCreateV1:
		_ = s.locationEventHandler.OnLocationCreate(ctx, evt)
		return nil
	case locationevents.LocationUpdateV1:
		_ = s.locationEventHandler.OnLocationUpdate(ctx, evt)
		return nil
	case locationevents.LocationValidationFailedV1:
		_ = s.locationEventHandler.OnLocationValidationFailed(ctx, evt)
		return nil
	case locationevents.LocationValidatedV1:
		_ = s.locationEventHandler.OnLocationValidated(ctx, evt)
		return nil
	case jobroleevents.JobRoleCreateV1:
		_ = s.jobRoleEventHandler.OnJobRoleCreate(ctx, evt)
		return nil

	case logentryevents.LogEntryCreateV1:
		_ = s.logEntryEventHandler.OnCreate(ctx, evt)
		return nil
	case logentryevents.LogEntryUpdateV1:
		_ = s.logEntryEventHandler.OnUpdate(ctx, evt)
		return nil
	case logentryevents.LogEntryAddTagV1:
		_ = s.logEntryEventHandler.OnAddTag(ctx, evt)
		return nil
	case logentryevents.LogEntryRemoveTagV1:
		_ = s.logEntryEventHandler.OnRemoveTag(ctx, evt)
		return nil

	case commentevent.CommentCreateV1:
		_ = s.commentEventHandler.OnCreate(ctx, evt)
		return nil
	case commentevent.CommentUpdateV1:
		_ = s.commentEventHandler.OnUpdate(ctx, evt)
		return nil

	case issueevent.IssueCreateV1:
		_ = s.issueEventHandler.OnCreate(ctx, evt)
		return nil
	case issueevent.IssueUpdateV1:
		_ = s.issueEventHandler.OnUpdate(ctx, evt)
		return nil
	case issueevent.IssueAddUserAssigneeV1:
		_ = s.issueEventHandler.OnAddUserAssignee(ctx, evt)
		return nil
	case issueevent.IssueRemoveUserAssigneeV1:
		_ = s.issueEventHandler.OnRemoveUserAssignee(ctx, evt)
		return nil
	case issueevent.IssueAddUserFollowerV1:
		_ = s.issueEventHandler.OnAddUserFollower(ctx, evt)
		return nil
	case issueevent.IssueRemoveUserFollowerV1:
		_ = s.issueEventHandler.OnRemoveUserFollower(ctx, evt)
		return nil

	case opportunityevent.OpportunityCreateV1:
		_ = s.opportunityEventHandler.OnCreate(ctx, evt)
		return nil
	case opportunityevent.OpportunityUpdateNextCycleDateV1:
		_ = s.opportunityEventHandler.OnUpdateNextCycleDate(ctx, evt)
		return nil
	case opportunityevent.OpportunityUpdateV1:
		_ = s.opportunityEventHandler.OnUpdate(ctx, evt)
		return nil
	case opportunityevent.OpportunityCreateRenewalV1:
		_ = s.opportunityEventHandler.OnCreateRenewal(ctx, evt)
		return nil
	case opportunityevent.OpportunityUpdateRenewalV1:
		_ = s.opportunityEventHandler.OnUpdateRenewal(ctx, evt)
		return nil
	case opportunityevent.OpportunityCloseWinV1:
		_ = s.opportunityEventHandler.OnCloseWon(ctx, evt)
		return nil
	case opportunityevent.OpportunityCloseLooseV1:
		_ = s.opportunityEventHandler.OnCloseLost(ctx, evt)
		return nil

	case contractevent.ContractCreateV1:
		_ = s.contractEventHandler.OnCreate(ctx, evt)
		return nil
	case contractevent.ContractUpdateV1:
		_ = s.contractEventHandler.OnUpdate(ctx, evt)
		return nil
	case contractevent.ContractRolloutRenewalOpportunityV1:
		_ = s.contractEventHandler.OnRolloutRenewalOpportunity(ctx, evt)
		return nil
	case contractevent.ContractDeleteV1:
		_ = s.contractEventHandler.OnDeleteV1(ctx, evt)
		return nil
	case contractevent.ContractRefreshStatusV1:
		_ = s.contractEventHandler.OnRefreshStatus(ctx, evt)
		return nil
	case contractevent.ContractRefreshLtvV1:
		_ = s.contractEventHandler.OnRefreshLtv(ctx, evt)
		return nil

	case servicelineitemevent.ServiceLineItemCreateV1:
		_ = s.serviceLineItemEventHandler.OnCreateV1(ctx, evt)
		return nil
	case servicelineitemevent.ServiceLineItemUpdateV1:
		_ = s.serviceLineItemEventHandler.OnUpdateV1(ctx, evt)
		return nil
	case servicelineitemevent.ServiceLineItemDeleteV1:
		_ = s.serviceLineItemEventHandler.OnDeleteV1(ctx, evt)
		return nil
	case servicelineitemevent.ServiceLineItemCloseV1:
		_ = s.serviceLineItemEventHandler.OnClose(ctx, evt)
		return nil
	case servicelineitemevent.ServiceLineItemPauseV1:
		_ = s.serviceLineItemEventHandler.OnPause(ctx, evt)
		return nil
	case servicelineitemevent.ServiceLineItemResumeV1:
		_ = s.serviceLineItemEventHandler.OnResume(ctx, evt)
		return nil

	case masterplanevent.MasterPlanCreateV1:
		_ = s.masterPlanEventHandler.OnCreate(ctx, evt)
		return nil
	case masterplanevent.MasterPlanUpdateV1:
		_ = s.masterPlanEventHandler.OnUpdate(ctx, evt)
		return nil
	case masterplanevent.MasterPlanMilestoneCreateV1:
		_ = s.masterPlanEventHandler.OnCreateMilestone(ctx, evt)
		return nil
	case masterplanevent.MasterPlanMilestoneUpdateV1:
		_ = s.masterPlanEventHandler.OnUpdateMilestone(ctx, evt)
		return nil
	case masterplanevent.MasterPlanMilestoneReorderV1:
		_ = s.masterPlanEventHandler.OnReorderMilestones(ctx, evt)
		return nil

	case invoiceevents.InvoiceCreateForContractV1:
		_ = s.invoiceEventHandler.OnInvoiceCreateForContractV1(ctx, evt)
		return nil
	case invoiceevents.InvoiceFillV1:
		_ = s.invoiceEventHandler.OnInvoiceFillV1(ctx, evt)
		return nil
	case invoiceevents.InvoicePdfGeneratedV1:
		_ = s.invoiceEventHandler.OnInvoicePdfGenerated(ctx, evt)
		return nil
	case invoiceevents.InvoiceUpdateV1:
		_ = s.invoiceEventHandler.OnInvoiceUpdateV1(ctx, evt)
		return nil
	case invoiceevents.InvoiceVoidV1:
		_ = s.invoiceEventHandler.OnInvoiceVoidV1(ctx, evt)
		return nil
	case invoiceevents.InvoiceDeleteV1:
		_ = s.invoiceEventHandler.OnInvoiceDeleteV1(ctx, evt)
		return nil

	case orgplanevent.OrganizationPlanCreateV1:
		_ = s.organizationPlanEventHandler.OnCreate(ctx, evt)
		return nil
	case orgplanevent.OrganizationPlanUpdateV1:
		_ = s.organizationPlanEventHandler.OnUpdate(ctx, evt)
		return nil
	case orgplanevent.OrganizationPlanMilestoneCreateV1:
		_ = s.organizationPlanEventHandler.OnCreateMilestone(ctx, evt)
		return nil
	case orgplanevent.OrganizationPlanMilestoneUpdateV1:
		_ = s.organizationPlanEventHandler.OnUpdateMilestone(ctx, evt)
		return nil
	case orgplanevent.OrganizationPlanMilestoneReorderV1:
		_ = s.organizationPlanEventHandler.OnReorderMilestones(ctx, evt)
		return nil

	case tenantevent.TenantAddBillingProfileV1:
		_ = s.tenantEventHandler.OnAddBillingProfileV1(ctx, evt)
		return nil
	case tenantevent.TenantUpdateBillingProfileV1:
		_ = s.tenantEventHandler.OnUpdateBillingProfileV1(ctx, evt)
		return nil
	case tenantevent.TenantUpdateSettingsV1:
		_ = s.tenantEventHandler.OnUpdateTenantSettingsV1(ctx, evt)
		return nil
	case tenantevent.TenantAddBankAccountV1:
		_ = s.bankAccountEventHandler.OnAddBankAccountV1(ctx, evt)
		return nil
	case tenantevent.TenantUpdateBankAccountV1:
		_ = s.bankAccountEventHandler.OnUpdateBankAccountV1(ctx, evt)
		return nil
	case tenantevent.TenantDeleteBankAccountV1:
		_ = s.bankAccountEventHandler.OnDeleteBankAccountV1(ctx, evt)
		return nil

	case reminderevents.ReminderCreateV1:
		_ = s.reminderEventHandler.OnCreate(ctx, evt)
		return nil
	case reminderevents.ReminderUpdateV1:
		_ = s.reminderEventHandler.OnUpdate(ctx, evt)
		return nil

	default:
		s.log.Errorf("(GraphSubscriber) Unknown EventType: {%s}", evt.EventType)
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		tracing.TraceErr(span, err)
		return err
	}
}
