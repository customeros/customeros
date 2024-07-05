package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	commentevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment"
	contactevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	contractevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	emailevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	ieevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	isevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/event"
	invoiceevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	invoicingcycleevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle"
	issueevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	jobroleevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	locationevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	logentryevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	masterplanevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	orderevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/order"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	orgplanevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	phonenumberevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	reminderevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder"
	servicelineitemevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	tenantevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	userevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"golang.org/x/sync/errgroup"

	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphSubscriber struct {
	log                            logger.Logger
	db                             *esdb.Client
	cfg                            *config.Config
	phoneNumberEventHandler        *PhoneNumberEventHandler
	contactEventHandler            *ContactEventHandler
	organizationEventHandler       *OrganizationEventHandler
	emailEventHandler              *EmailEventHandler
	userEventHandler               *UserEventHandler
	locationEventHandler           *LocationEventHandler
	jobRoleEventHandler            *JobRoleEventHandler
	interactionEventHandler        *InteractionEventHandler
	interactionSessionEventHandler *InteractionSessionEventHandler
	logEntryEventHandler           *LogEntryEventHandler
	issueEventHandler              *IssueEventHandler
	commentEventHandler            *CommentEventHandler
	opportunityEventHandler        *OpportunityEventHandler
	contractEventHandler           *ContractEventHandler
	serviceLineItemEventHandler    *ServiceLineItemEventHandler
	masterPlanEventHandler         *MasterPlanEventHandler
	invoicingCycleEventHandler     *InvoicingCycleEventHandler
	invoiceEventHandler            *InvoiceEventHandler
	tenantEventHandler             *TenantEventHandler
	organizationPlanEventHandler   *OrganizationPlanEventHandler
	bankAccountEventHandler        *BankAccountEventHandler
	reminderEventHandler           *ReminderEventHandler
	orderEventHandler              *OrderEventHandler
}

func NewGraphSubscriber(log logger.Logger, db *esdb.Client, commonServices *commonService.Services, repositories *repository.Repositories, grpcClients *grpc_client.Clients, cfg *config.Config, cache caches.Cache) *GraphSubscriber {
	return &GraphSubscriber{
		log:                            log,
		db:                             db,
		cfg:                            cfg,
		contactEventHandler:            NewContactEventHandler(log, repositories),
		organizationEventHandler:       NewOrganizationEventHandler(log, commonServices, repositories, grpcClients, cache),
		phoneNumberEventHandler:        NewPhoneNumberEventHandler(log, repositories, grpcClients),
		emailEventHandler:              NewEmailEventHandler(log, repositories, grpcClients),
		userEventHandler:               NewUserEventHandler(log, repositories),
		locationEventHandler:           NewLocationEventHandler(repositories),
		jobRoleEventHandler:            NewJobRoleEventHandler(repositories),
		interactionEventHandler:        NewInteractionEventHandler(log, repositories, grpcClients),
		interactionSessionEventHandler: NewInteractionSessionEventHandler(log, repositories, grpcClients),
		logEntryEventHandler:           NewLogEntryEventHandler(log, repositories, grpcClients),
		issueEventHandler:              NewIssueEventHandler(log, repositories, grpcClients),
		commentEventHandler:            NewCommentEventHandler(log, repositories),
		opportunityEventHandler:        NewOpportunityEventHandler(log, repositories, grpcClients),
		contractEventHandler:           NewContractEventHandler(log, repositories, grpcClients),
		serviceLineItemEventHandler:    NewServiceLineItemEventHandler(log, repositories, grpcClients),
		masterPlanEventHandler:         NewMasterPlanEventHandler(log, repositories),
		invoicingCycleEventHandler:     NewInvoicingCycleEventHandler(log, repositories),
		invoiceEventHandler:            NewInvoiceEventHandler(log, repositories, grpcClients),
		tenantEventHandler:             NewTenantEventHandler(log, repositories),
		organizationPlanEventHandler:   NewOrganizationPlanEventHandler(log, repositories),
		bankAccountEventHandler:        NewBankAccountEventHandler(log, repositories),
		reminderEventHandler:           NewReminderEventHandler(log, repositories),
		orderEventHandler:              NewOrderEventHandler(log, repositories, grpcClients),
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

	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
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
	case phonenumberevents.PhoneNumberValidateV1:
		return nil

	case emailevents.EmailCreateV1:
		return s.emailEventHandler.OnEmailCreate(ctx, evt)
	case emailevents.EmailUpdateV1:
		return s.emailEventHandler.OnEmailUpdate(ctx, evt)
	case emailevents.EmailValidationFailedV1:
		return s.emailEventHandler.OnEmailValidationFailed(ctx, evt)
	case emailevents.EmailValidatedV1:
		return s.emailEventHandler.OnEmailValidated(ctx, evt)
	case emailevents.EmailValidateV1:
		return nil

	case contactevent.ContactCreateV1:
		return s.contactEventHandler.OnContactCreate(ctx, evt)
	case contactevent.ContactUpdateV1:
		return s.contactEventHandler.OnContactUpdate(ctx, evt)
	case contactevent.ContactPhoneNumberLinkV1:
		return s.contactEventHandler.OnPhoneNumberLinkToContact(ctx, evt)
	case contactevent.ContactEmailLinkV1:
		return s.contactEventHandler.OnEmailLinkToContact(ctx, evt)
	case contactevent.ContactLocationLinkV1:
		return s.contactEventHandler.OnLocationLinkToContact(ctx, evt)
	case contactevent.ContactOrganizationLinkV1:
		return s.contactEventHandler.OnContactLinkToOrganization(ctx, evt)
	case contactevent.ContactAddSocialV1:
		return s.contactEventHandler.OnSocialAddedToContactV1(ctx, evt)
	case contactevent.ContactRemoveSocialV1:
		return s.contactEventHandler.OnSocialRemovedFromContactV1(ctx, evt)
	case contactevent.ContactAddTagV1:
		return s.contactEventHandler.OnAddTag(ctx, evt)
	case contactevent.ContactRemoveTagV1:
		return s.contactEventHandler.OnRemoveTag(ctx, evt)
	case contactevent.ContactRequestEnrichV1:
		return nil
	case contactevent.ContactAddLocationV1:
		return s.contactEventHandler.OnLocationAddedToContact(ctx, evt)

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
	case orgevents.OrganizationUnlinkDomainV1:
		return s.organizationEventHandler.OnDomainUnlinkedFromOrganization(ctx, evt)
	case orgevents.OrganizationAddSocialV1:
		return s.organizationEventHandler.OnSocialAddedToOrganization(ctx, evt)
	case orgevents.OrganizationRemoveSocialV1:
		return s.organizationEventHandler.OnSocialRemovedFromOrganization(ctx, evt)
	case orgevents.OrganizationHideV1:
		return s.organizationEventHandler.OnOrganizationHide(ctx, evt)
	case orgevents.OrganizationShowV1:
		return s.organizationEventHandler.OnOrganizationShow(ctx, evt)
	case orgevents.OrganizationRefreshLastTouchpointV1:
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
		return s.organizationEventHandler.OnUpsertCustomField(ctx, evt)
	case orgevents.OrganizationAddParentV1:
		return s.organizationEventHandler.OnLinkWithParentOrganization(ctx, evt)
	case orgevents.OrganizationRemoveParentV1:
		return s.organizationEventHandler.OnUnlinkFromParentOrganization(ctx, evt)
	case orgevents.OrganizationUpdateOnboardingStatusV1:
		return s.organizationEventHandler.OnUpdateOnboardingStatus(ctx, evt)
	case orgevents.OrganizationUpdateOwnerV1:
		return s.organizationEventHandler.OnUpdateOwner(ctx, evt)
	case orgevents.OrganizationAddTagV1:
		return s.organizationEventHandler.OnAddTag(ctx, evt)
	case orgevents.OrganizationRemoveTagV1:
		return s.organizationEventHandler.OnRemoveTag(ctx, evt)
	case orgevents.OrganizationRequestRenewalForecastV1,
		orgevents.OrganizationRequestNextCycleDateV1,
		orgevents.OrganizationUpdateRenewalLikelihoodV1,
		orgevents.OrganizationUpdateRenewalForecastV1,
		orgevents.OrganizationUpdateBillingDetailsV1,
		orgevents.OrganizationRequestScrapeByWebsiteV1,
		orgevents.OrganizationAdjustIndustryV1,
		orgevents.OrganizationRequestEnrichV1:
		return nil
	case orgevents.OrganizationCreateBillingProfileV1:
		return s.organizationEventHandler.OnCreateBillingProfile(ctx, evt)
	case orgevents.OrganizationUpdateBillingProfileV1:
		return s.organizationEventHandler.OnUpdateBillingProfile(ctx, evt)
	case orgevents.OrganizationEmailLinkToBillingProfileV1:
		return s.organizationEventHandler.OnEmailLinkedToBillingProfile(ctx, evt)
	case orgevents.OrganizationEmailUnlinkFromBillingProfileV1:
		return s.organizationEventHandler.OnEmailUnlinkedFromBillingProfile(ctx, evt)
	case orgevents.OrganizationLocationLinkToBillingProfileV1:
		return s.organizationEventHandler.OnLocationLinkedToBillingProfile(ctx, evt)
	case orgevents.OrganizationLocationUnlinkFromBillingProfileV1:
		return s.organizationEventHandler.OnLocationUnlinkedFromBillingProfile(ctx, evt)
	case orgevents.OrganizationAddLocationV1:
		return s.organizationEventHandler.OnLocationAddedToOrganization(ctx, evt)

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

	case ieevent.InteractionEventRequestSummaryV1,
		ieevent.InteractionEventRequestActionItemsV1:
		return nil
	case ieevent.InteractionEventReplaceSummaryV1:
		return s.interactionEventHandler.OnSummaryReplace(ctx, evt)
	case ieevent.InteractionEventReplaceActionItemsV1:
		return s.interactionEventHandler.OnActionItemsReplace(ctx, evt)
	case ieevent.InteractionEventCreateV1:
		return s.interactionEventHandler.OnCreate(ctx, evt)
	case ieevent.InteractionEventUpdateV1:
		return s.interactionEventHandler.OnUpdate(ctx, evt)

	case isevent.InteractionSessionCreateV1:
		return s.interactionSessionEventHandler.OnCreate(ctx, evt)

	case logentryevents.LogEntryCreateV1:
		return s.logEntryEventHandler.OnCreate(ctx, evt)
	case logentryevents.LogEntryUpdateV1:
		return s.logEntryEventHandler.OnUpdate(ctx, evt)
	case logentryevents.LogEntryAddTagV1:
		return s.logEntryEventHandler.OnAddTag(ctx, evt)
	case logentryevents.LogEntryRemoveTagV1:
		return s.logEntryEventHandler.OnRemoveTag(ctx, evt)

	case commentevent.CommentCreateV1:
		return s.commentEventHandler.OnCreate(ctx, evt)
	case commentevent.CommentUpdateV1:
		return s.commentEventHandler.OnUpdate(ctx, evt)

	case issueevent.IssueCreateV1:
		return s.issueEventHandler.OnCreate(ctx, evt)
	case issueevent.IssueUpdateV1:
		return s.issueEventHandler.OnUpdate(ctx, evt)
	case issueevent.IssueAddUserAssigneeV1:
		return s.issueEventHandler.OnAddUserAssignee(ctx, evt)
	case issueevent.IssueRemoveUserAssigneeV1:
		return s.issueEventHandler.OnRemoveUserAssignee(ctx, evt)
	case issueevent.IssueAddUserFollowerV1:
		return s.issueEventHandler.OnAddUserFollower(ctx, evt)
	case issueevent.IssueRemoveUserFollowerV1:
		return s.issueEventHandler.OnRemoveUserFollower(ctx, evt)

	case opportunityevent.OpportunityCreateV1:
		return s.opportunityEventHandler.OnCreate(ctx, evt)
	case opportunityevent.OpportunityUpdateNextCycleDateV1:
		return s.opportunityEventHandler.OnUpdateNextCycleDate(ctx, evt)
	case opportunityevent.OpportunityUpdateV1:
		return s.opportunityEventHandler.OnUpdate(ctx, evt)
	case opportunityevent.OpportunityCreateRenewalV1:
		return s.opportunityEventHandler.OnCreateRenewal(ctx, evt)
	case opportunityevent.OpportunityUpdateRenewalV1:
		return s.opportunityEventHandler.OnUpdateRenewal(ctx, evt)
	case opportunityevent.OpportunityCloseWinV1:
		return s.opportunityEventHandler.OnCloseWon(ctx, evt)
	case opportunityevent.OpportunityCloseLooseV1:
		return s.opportunityEventHandler.OnCloseLost(ctx, evt)

	case contractevent.ContractCreateV1:
		return s.contractEventHandler.OnCreate(ctx, evt)
	case contractevent.ContractUpdateV1:
		return s.contractEventHandler.OnUpdate(ctx, evt)
	case contractevent.ContractRolloutRenewalOpportunityV1:
		return s.contractEventHandler.OnRolloutRenewalOpportunity(ctx, evt)
	case contractevent.ContractUpdateStatusV1:
		return nil
	case contractevent.ContractDeleteV1:
		return s.contractEventHandler.OnDeleteV1(ctx, evt)
	case contractevent.ContractRefreshStatusV1:
		_ = s.contractEventHandler.OnRefreshStatus(ctx, evt)
		return nil
	case contractevent.ContractRefreshLtvV1:
		_ = s.contractEventHandler.OnRefreshLtv(ctx, evt)
		return nil

	case servicelineitemevent.ServiceLineItemCreateV1:
		return s.serviceLineItemEventHandler.OnCreateV1(ctx, evt)
	case servicelineitemevent.ServiceLineItemUpdateV1:
		return s.serviceLineItemEventHandler.OnUpdateV1(ctx, evt)
	case servicelineitemevent.ServiceLineItemDeleteV1:
		return s.serviceLineItemEventHandler.OnDeleteV1(ctx, evt)
	case servicelineitemevent.ServiceLineItemCloseV1:
		return s.serviceLineItemEventHandler.OnClose(ctx, evt)

	case masterplanevent.MasterPlanCreateV1:
		return s.masterPlanEventHandler.OnCreate(ctx, evt)
	case masterplanevent.MasterPlanUpdateV1:
		return s.masterPlanEventHandler.OnUpdate(ctx, evt)
	case masterplanevent.MasterPlanMilestoneCreateV1:
		return s.masterPlanEventHandler.OnCreateMilestone(ctx, evt)
	case masterplanevent.MasterPlanMilestoneUpdateV1:
		return s.masterPlanEventHandler.OnUpdateMilestone(ctx, evt)
	case masterplanevent.MasterPlanMilestoneReorderV1:
		return s.masterPlanEventHandler.OnReorderMilestones(ctx, evt)

	case invoicingcycleevent.InvoicingCycleCreateV1:
		return s.invoicingCycleEventHandler.OnCreate(ctx, evt)
	case invoicingcycleevent.InvoicingCycleUpdateV1:
		return s.invoicingCycleEventHandler.OnUpdate(ctx, evt)

	case invoiceevents.InvoiceCreateForContractV1:
		return s.invoiceEventHandler.OnInvoiceCreateForContractV1(ctx, evt)
	case invoiceevents.InvoiceFillV1:
		return s.invoiceEventHandler.OnInvoiceFillV1(ctx, evt)
	case invoiceevents.InvoicePdfGeneratedV1:
		return s.invoiceEventHandler.OnInvoicePdfGenerated(ctx, evt)
	case invoiceevents.InvoiceUpdateV1:
		return s.invoiceEventHandler.OnInvoiceUpdateV1(ctx, evt)
	case invoiceevents.InvoiceVoidV1:
		return s.invoiceEventHandler.OnInvoiceVoidV1(ctx, evt)
	case invoiceevents.InvoiceDeleteV1:
		return s.invoiceEventHandler.OnInvoiceDeleteV1(ctx, evt)
	case invoiceevents.InvoicePdfRequestedV1,
		invoiceevents.InvoicePaidV1,
		invoiceevents.InvoiceFillRequestedV1,
		invoiceevents.InvoicePayNotificationV1:
		return nil // do nothing
	case invoiceevents.InvoicePayV1:
		return nil

	case orgplanevent.OrganizationPlanCreateV1:
		return s.organizationPlanEventHandler.OnCreate(ctx, evt)
	case orgplanevent.OrganizationPlanUpdateV1:
		return s.organizationPlanEventHandler.OnUpdate(ctx, evt)
	case orgplanevent.OrganizationPlanMilestoneCreateV1:
		return s.organizationPlanEventHandler.OnCreateMilestone(ctx, evt)
	case orgplanevent.OrganizationPlanMilestoneUpdateV1:
		return s.organizationPlanEventHandler.OnUpdateMilestone(ctx, evt)
	case orgplanevent.OrganizationPlanMilestoneReorderV1:
		return s.organizationPlanEventHandler.OnReorderMilestones(ctx, evt)

	case tenantevent.TenantAddBillingProfileV1:
		return s.tenantEventHandler.OnAddBillingProfileV1(ctx, evt)
	case tenantevent.TenantUpdateBillingProfileV1:
		return s.tenantEventHandler.OnUpdateBillingProfileV1(ctx, evt)
	case tenantevent.TenantUpdateSettingsV1:
		return s.tenantEventHandler.OnUpdateTenantSettingsV1(ctx, evt)
	case tenantevent.TenantAddBankAccountV1:
		return s.bankAccountEventHandler.OnAddBankAccountV1(ctx, evt)
	case tenantevent.TenantUpdateBankAccountV1:
		return s.bankAccountEventHandler.OnUpdateBankAccountV1(ctx, evt)
	case tenantevent.TenantDeleteBankAccountV1:
		return s.bankAccountEventHandler.OnDeleteBankAccountV1(ctx, evt)

	case orgevents.OrganizationUpdateOwnerNotificationV1:
		return nil

	case reminderevents.ReminderCreateV1:
		return s.reminderEventHandler.OnCreate(ctx, evt)
	case reminderevents.ReminderUpdateV1:
		return s.reminderEventHandler.OnUpdate(ctx, evt)
	case reminderevents.ReminderNotificationV1:
		return nil

	case orderevents.OrderUpsertV1:
		return s.orderEventHandler.OnUpsertOrderV1(ctx, evt)

	default:
		s.log.Errorf("(GraphSubscriber) Unknown EventType: {%s}", evt.EventType)
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}
