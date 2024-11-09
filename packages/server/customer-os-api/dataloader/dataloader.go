package dataloader

import (
	"context"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"net/http"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

const defaultDataloaderWaitTime = 32 * time.Millisecond

type loadersString string

const loadersKey = loadersString("dataloaders")

type Loaders struct {
	// tags
	TagsForOrganization *dataloader.Loader
	TagsForContact      *dataloader.Loader
	TagsForIssue        *dataloader.Loader
	TagsForLogEntry     *dataloader.Loader
	// emails
	EmailsForContact       *dataloader.Loader
	EmailsForOrganization  *dataloader.Loader
	PrimaryEmailForContact *dataloader.Loader

	LocationsForContact                         *dataloader.Loader
	LocationsForOrganization                    *dataloader.Loader
	JobRolesForContact                          *dataloader.Loader
	JobRolesForOrganization                     *dataloader.Loader
	JobRolesForUser                             *dataloader.Loader
	CalendarsForUser                            *dataloader.Loader
	DomainsForOrganization                      *dataloader.Loader
	InteractionEventsForInteractionSession      *dataloader.Loader
	InteractionSessionForInteractionEvent       *dataloader.Loader
	SentByParticipantsForInteractionEvent       *dataloader.Loader
	SentToParticipantsForInteractionEvent       *dataloader.Loader
	AttendedByParticipantsForInteractionSession *dataloader.Loader
	ReplyToInteractionEventForInteractionEvent  *dataloader.Loader
	PhoneNumbersForOrganization                 *dataloader.Loader
	PhoneNumbersForUser                         *dataloader.Loader
	PhoneNumbersForContact                      *dataloader.Loader
	UsersConnectedForContact                    *dataloader.Loader
	UsersForEmail                               *dataloader.Loader
	UsersForPhoneNumber                         *dataloader.Loader
	UsersForPlayer                              *dataloader.Loader
	UserOwnerForOrganization                    *dataloader.Loader
	UserOwnerForOpportunity                     *dataloader.Loader
	UserCreatorForOpportunity                   *dataloader.Loader
	UserCreatorForServiceLineItem               *dataloader.Loader
	UserCreatorForContract                      *dataloader.Loader
	UserAuthorForLogEntry                       *dataloader.Loader
	UserAuthorForComment                        *dataloader.Loader
	UserForFlowSender                           *dataloader.Loader
	User                                        *dataloader.Loader
	ContactsForEmail                            *dataloader.Loader
	ContactsForPhoneNumber                      *dataloader.Loader
	ContactCountForOrganization                 *dataloader.Loader
	Organization                                *dataloader.Loader
	OrganizationsForEmail                       *dataloader.Loader
	OrganizationsForPhoneNumber                 *dataloader.Loader
	SubsidiariesForOrganization                 *dataloader.Loader
	SubsidiariesOfForOrganization               *dataloader.Loader
	SuggestedMergeToForOrganization             *dataloader.Loader
	DescribedByFor                              *dataloader.Loader
	CreatedByParticipantsForMeeting             *dataloader.Loader
	AttendedByParticipantsForMeeting            *dataloader.Loader
	InteractionEventsForMeeting                 *dataloader.Loader
	InteractionEventsForIssue                   *dataloader.Loader
	CommentsForIssue                            *dataloader.Loader
	NotesForMeeting                             *dataloader.Loader
	AttachmentsForInteractionEvent              *dataloader.Loader
	AttachmentsForMeeting                       *dataloader.Loader
	AttachmentsForContract                      *dataloader.Loader
	SocialsForContact                           *dataloader.Loader
	SocialsForOrganization                      *dataloader.Loader
	ExternalSystemsForComment                   *dataloader.Loader
	ExternalSystemsForIssue                     *dataloader.Loader
	ExternalSystemsForOrganization              *dataloader.Loader
	ExternalSystemsForLogEntry                  *dataloader.Loader
	ExternalSystemsForMeeting                   *dataloader.Loader
	ExternalSystemsForInteractionEvent          *dataloader.Loader
	ExternalSystemsForContract                  *dataloader.Loader
	ExternalSystemsForOpportunity               *dataloader.Loader
	ExternalSystemsForServiceLineItem           *dataloader.Loader
	TimelineEventForTimelineEventId             *dataloader.Loader
	InboundCommsCountForOrganization            *dataloader.Loader
	OutboundCommsCountForOrganization           *dataloader.Loader
	OrganizationForJobRole                      *dataloader.Loader
	OrganizationForInvoice                      *dataloader.Loader
	OrganizationForOpportunity                  *dataloader.Loader
	OrganizationForSlackChannel                 *dataloader.Loader
	LatestOrganizationWithJobRoleForContact     *dataloader.Loader
	ContactForJobRole                           *dataloader.Loader
	IssueForInteractionEvent                    *dataloader.Loader
	MeetingForInteractionEvent                  *dataloader.Loader
	CountryForPhoneNumber                       *dataloader.Loader
	ActionsForInteractionEvent                  *dataloader.Loader
	ActionItemsForInteractionEvent              *dataloader.Loader
	SubmitterParticipantsForIssue               *dataloader.Loader
	ReporterParticipantsForIssue                *dataloader.Loader
	AssigneeParticipantsForIssue                *dataloader.Loader
	FollowerParticipantsForIssue                *dataloader.Loader
	ContractsForOrganization                    *dataloader.Loader
	ContractForInvoice                          *dataloader.Loader
	ServiceLineItemsForContract                 *dataloader.Loader
	ServiceLineItemForInvoiceLine               *dataloader.Loader
	OpportunitiesForContract                    *dataloader.Loader
	OpportunitiesForOrganization                *dataloader.Loader
	InvoiceLinesForInvoice                      *dataloader.Loader
	OrdersForOrganization                       *dataloader.Loader
	InvoicesForContract                         *dataloader.Loader
	FlowContactsForFlow                         *dataloader.Loader
	FlowActionsForFlow                          *dataloader.Loader
	FlowSendersForFlow                          *dataloader.Loader
	FlowsWithContact                            *dataloader.Loader
	FlowsWithSender                             *dataloader.Loader
}

type tagBatcher struct {
	tagService commonservice.TagService
}
type emailBatcher struct {
	emailService       service.EmailService
	commonEmailService commonservice.EmailService
}
type locationBatcher struct {
	locationService service.LocationService
}
type socialBatcher struct {
	socialService commonservice.SocialService
}
type jobRoleBatcher struct {
	jobRoleService commonservice.JobRoleService
}
type calendarBatcher struct {
	calendarService service.CalendarService
}
type domainBatcher struct {
	domainService commonservice.DomainService
}
type interactionEventBatcher struct {
	interactionEventService commonservice.InteractionEventService
}
type interactionSessionBatcher struct {
	interactionSessionService commonservice.InteractionSessionService
}
type interactionEventParticipantBatcher struct {
	interactionEventService commonservice.InteractionEventService
}
type interactionSessionParticipantBatcher struct {
	interactionSessionService commonservice.InteractionSessionService
}
type meetingParticipantBatcher struct {
	meetingService service.MeetingService
}
type phoneNumberBatcher struct {
	phoneNumberService service.PhoneNumberService
}
type notedEntityBatcher struct {
	noteService service.NoteService
}
type userBatcher struct {
	userService service.UserService
}
type contactBatcher struct {
	contactService service.ContactService
}
type organizationBatcher struct {
	organizationService       service.OrganizationService
	commonOrganizationService commonservice.OrganizationService
}
type noteBatcher struct {
	noteService service.NoteService
}
type attachmentBatcher struct {
	attachmentService commonservice.AttachmentService
}
type externalSystemBatcher struct {
	externalSystemService commonservice.ExternalSystemService
}
type timelineEventBatcher struct {
	timelineEventService service.TimelineEventService
}
type issueBatcher struct {
	issueService service.IssueService
}
type meetingBatcher struct {
	meetingService service.MeetingService
}
type countryBatcher struct {
	countryService service.CountryService
}
type actionBatcher struct {
	actionService service.ActionService
}
type actionItemBatcher struct {
	actionItemService service.ActionItemService
}
type issueParticipantBatcher struct {
	issueService service.IssueService
}
type commentBatcher struct {
	commentService service.CommentService
}
type contractBatcher struct {
	contractService service.ContractService
}
type serviceLineItemBatcher struct {
	serviceLineItemService commonservice.ServiceLineItemService
}
type opportunityBatcher struct {
	opportunityService commonservice.OpportunityService
}
type invoiceBatcher struct {
	invoiceService commonservice.InvoiceService
}
type flowBatcher struct {
	flowService commonservice.FlowService
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(services *service.Services) *Loaders {
	tagBatcher := &tagBatcher{
		tagService: services.CommonServices.TagService,
	}
	emailBatcher := &emailBatcher{
		emailService:       services.EmailService,
		commonEmailService: services.CommonServices.EmailService,
	}
	locationBatcher := &locationBatcher{
		locationService: services.LocationService,
	}
	socialBatcher := &socialBatcher{
		socialService: services.CommonServices.SocialService,
	}
	jobRoleBatcher := &jobRoleBatcher{
		jobRoleService: services.CommonServices.JobRoleService,
	}
	calendarBatcher := &calendarBatcher{
		calendarService: services.CalendarService,
	}
	domainBatcher := &domainBatcher{
		domainService: services.CommonServices.DomainService,
	}
	interactionEventBatcher := &interactionEventBatcher{
		interactionEventService: services.CommonServices.InteractionEventService,
	}
	commentBatcher := &commentBatcher{
		commentService: services.CommentService,
	}
	interactionSessionBatcher := &interactionSessionBatcher{
		interactionSessionService: services.CommonServices.InteractionSessionService,
	}
	interactionEventParticipantBatcher := &interactionEventParticipantBatcher{
		interactionEventService: services.CommonServices.InteractionEventService,
	}
	issueParticipantBatcher := &issueParticipantBatcher{
		issueService: services.IssueService,
	}
	interactionSessionParticipantBatcher := &interactionSessionParticipantBatcher{
		interactionSessionService: services.CommonServices.InteractionSessionService,
	}
	meetingParticipantBatcher := &meetingParticipantBatcher{
		meetingService: services.MeetingService,
	}
	phoneNumberBatcher := &phoneNumberBatcher{
		phoneNumberService: services.PhoneNumberService,
	}
	userBatcher := &userBatcher{
		userService: services.UserService,
	}
	contactBatcher := &contactBatcher{
		contactService: services.ContactService,
	}
	organizationBatcher := &organizationBatcher{
		organizationService:       services.OrganizationService,
		commonOrganizationService: services.CommonServices.OrganizationService,
	}
	noteBatcher := &noteBatcher{
		noteService: services.NoteService,
	}
	attachmentBatcher := attachmentBatcher{
		attachmentService: services.CommonServices.AttachmentService,
	}
	externalSystemBatcher := externalSystemBatcher{
		externalSystemService: services.ExternalSystemService,
	}
	timelineEventBatcher := timelineEventBatcher{
		timelineEventService: services.TimelineEventService,
	}
	issueBatcher := issueBatcher{
		issueService: services.IssueService,
	}
	meetingBatcher := meetingBatcher{
		meetingService: services.MeetingService,
	}
	countryBatcher := countryBatcher{
		countryService: services.CountryService,
	}
	actionBatcher := actionBatcher{
		actionService: services.ActionService,
	}
	actionItemBatcher := actionItemBatcher{
		actionItemService: services.ActionItemService,
	}
	contractBatcher := &contractBatcher{
		contractService: services.ContractService,
	}
	serviceLineItemBatcher := &serviceLineItemBatcher{
		serviceLineItemService: services.CommonServices.ServiceLineItemService,
	}
	opportunityBatcher := &opportunityBatcher{
		opportunityService: services.CommonServices.OpportunityService,
	}
	invoiceBatcher := &invoiceBatcher{
		invoiceService: services.CommonServices.InvoiceService,
	}
	flowBatcher := &flowBatcher{
		flowService: services.CommonServices.FlowService,
	}
	return &Loaders{
		TagsForOrganization:                         dataloader.NewBatchedLoader(tagBatcher.getTagsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		TagsForContact:                              dataloader.NewBatchedLoader(tagBatcher.getTagsForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		TagsForIssue:                                dataloader.NewBatchedLoader(tagBatcher.getTagsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		TagsForLogEntry:                             dataloader.NewBatchedLoader(tagBatcher.getTagsForLogEntries, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		EmailsForContact:                            dataloader.NewBatchedLoader(emailBatcher.getEmailsForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		EmailsForOrganization:                       dataloader.NewBatchedLoader(emailBatcher.getEmailsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		PrimaryEmailForContact:                      dataloader.NewBatchedLoader(emailBatcher.getPrimaryEmailForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		LocationsForContact:                         dataloader.NewBatchedLoader(locationBatcher.getLocationsForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		LocationsForOrganization:                    dataloader.NewBatchedLoader(locationBatcher.getLocationsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		JobRolesForContact:                          dataloader.NewBatchedLoader(jobRoleBatcher.getJobRolesForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		JobRolesForOrganization:                     dataloader.NewBatchedLoader(jobRoleBatcher.getJobRolesForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		JobRolesForUser:                             dataloader.NewBatchedLoader(jobRoleBatcher.getJobRolesForUsers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		CalendarsForUser:                            dataloader.NewBatchedLoader(calendarBatcher.getCalendarsForUsers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		DomainsForOrganization:                      dataloader.NewBatchedLoader(domainBatcher.getDomainsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		InteractionEventsForInteractionSession:      dataloader.NewBatchedLoader(interactionEventBatcher.getInteractionEventsForInteractionSessions, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		InteractionEventsForMeeting:                 dataloader.NewBatchedLoader(interactionEventBatcher.getInteractionEventsForMeetings, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		InteractionEventsForIssue:                   dataloader.NewBatchedLoader(interactionEventBatcher.getInteractionEventsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		InteractionSessionForInteractionEvent:       dataloader.NewBatchedLoader(interactionSessionBatcher.getInteractionSessionsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SentByParticipantsForInteractionEvent:       dataloader.NewBatchedLoader(interactionEventParticipantBatcher.getSentByParticipantsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SentToParticipantsForInteractionEvent:       dataloader.NewBatchedLoader(interactionEventParticipantBatcher.getSentToParticipantsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AttendedByParticipantsForInteractionSession: dataloader.NewBatchedLoader(interactionSessionParticipantBatcher.getAttendedByParticipantsForInteractionSessions, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		CreatedByParticipantsForMeeting:             dataloader.NewBatchedLoader(meetingParticipantBatcher.getCreatedByParticipantsForMeeting, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AttendedByParticipantsForMeeting:            dataloader.NewBatchedLoader(meetingParticipantBatcher.getAttendedByParticipantsForMeeting, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		PhoneNumbersForOrganization:                 dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		PhoneNumbersForUser:                         dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForUsers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		PhoneNumbersForContact:                      dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ReplyToInteractionEventForInteractionEvent:  dataloader.NewBatchedLoader(interactionEventBatcher.getReplyToInteractionEventsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UsersConnectedForContact:                    dataloader.NewBatchedLoader(userBatcher.getUsersConnectedForContact, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UsersForEmail:                               dataloader.NewBatchedLoader(userBatcher.getUsersForEmails, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UsersForPhoneNumber:                         dataloader.NewBatchedLoader(userBatcher.getUsersForPhoneNumbers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UsersForPlayer:                              dataloader.NewBatchedLoader(userBatcher.getUsersForPlayers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserOwnerForOrganization:                    dataloader.NewBatchedLoader(userBatcher.getUserOwnersForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserOwnerForOpportunity:                     dataloader.NewBatchedLoader(userBatcher.getUserOwnersForOpportunities, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserCreatorForOpportunity:                   dataloader.NewBatchedLoader(userBatcher.getUserCreatorsForOpportunities, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserCreatorForServiceLineItem:               dataloader.NewBatchedLoader(userBatcher.getUserCreatorsForServiceLineItems, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserCreatorForContract:                      dataloader.NewBatchedLoader(userBatcher.getUserCreatorsForContracts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserAuthorForLogEntry:                       dataloader.NewBatchedLoader(userBatcher.getUserAuthorsForLogEntries, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserAuthorForComment:                        dataloader.NewBatchedLoader(userBatcher.getUserAuthorsForComments, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserForFlowSender:                           dataloader.NewBatchedLoader(userBatcher.getUserForFlowSenders, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		User:                                        dataloader.NewBatchedLoader(userBatcher.getUsers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContactsForEmail:                            dataloader.NewBatchedLoader(contactBatcher.getContactsForEmails, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContactsForPhoneNumber:                      dataloader.NewBatchedLoader(contactBatcher.getContactsForPhoneNumbers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContactCountForOrganization:                 dataloader.NewBatchedLoader(contactBatcher.getContactCountForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		Organization:                                dataloader.NewBatchedLoader(organizationBatcher.getOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationsForEmail:                       dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForEmails, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationsForPhoneNumber:                 dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForPhoneNumbers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SubsidiariesForOrganization:                 dataloader.NewBatchedLoader(organizationBatcher.getSubsidiariesForOrganization, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SubsidiariesOfForOrganization:               dataloader.NewBatchedLoader(organizationBatcher.getSubsidiariesOfForOrganization, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(20*time.Millisecond), dataloader.WithWait(defaultDataloaderWaitTime)),
		SuggestedMergeToForOrganization:             dataloader.NewBatchedLoader(organizationBatcher.getSuggestedMergeToForOrganization, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		NotesForMeeting:                             dataloader.NewBatchedLoader(noteBatcher.getNotesForMeetings, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AttachmentsForInteractionEvent:              dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AttachmentsForMeeting:                       dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForMeetings, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AttachmentsForContract:                      dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForContracts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SocialsForContact:                           dataloader.NewBatchedLoader(socialBatcher.getSocialsForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SocialsForOrganization:                      dataloader.NewBatchedLoader(socialBatcher.getSocialsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForComment:                   dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForComments, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForIssue:                     dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForOrganization:              dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForLogEntry:                  dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForLogEntries, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForMeeting:                   dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForMeetings, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForInteractionEvent:          dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForContract:                  dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForContracts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForOpportunity:               dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForOpportunities, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForServiceLineItem:           dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForServiceLineItems, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		TimelineEventForTimelineEventId:             dataloader.NewBatchedLoader(timelineEventBatcher.getTimelineEventsForTimelineEventIds, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		InboundCommsCountForOrganization:            dataloader.NewBatchedLoader(timelineEventBatcher.getInboundCommsCountForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OutboundCommsCountForOrganization:           dataloader.NewBatchedLoader(timelineEventBatcher.getOutboundCommsCountForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationForJobRole:                      dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForJobRoles, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationForInvoice:                      dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForInvoices, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationForOpportunity:                  dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForOpportunities, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationForSlackChannel:                 dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForSlackChannels, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		LatestOrganizationWithJobRoleForContact:     dataloader.NewBatchedLoader(organizationBatcher.getLatestOrganizationWithJobRoleForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContactForJobRole:                           dataloader.NewBatchedLoader(contactBatcher.getContactsForJobRoles, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		IssueForInteractionEvent:                    dataloader.NewBatchedLoader(issueBatcher.getIssuesForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		MeetingForInteractionEvent:                  dataloader.NewBatchedLoader(meetingBatcher.getMeetingsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		CountryForPhoneNumber:                       dataloader.NewBatchedLoader(countryBatcher.getCountriesForPhoneNumbers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ActionsForInteractionEvent:                  dataloader.NewBatchedLoader(actionBatcher.getActionsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ActionItemsForInteractionEvent:              dataloader.NewBatchedLoader(actionItemBatcher.getActionItemsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SubmitterParticipantsForIssue:               dataloader.NewBatchedLoader(issueParticipantBatcher.getSubmitterParticipantsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ReporterParticipantsForIssue:                dataloader.NewBatchedLoader(issueParticipantBatcher.getReporterParticipantsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AssigneeParticipantsForIssue:                dataloader.NewBatchedLoader(issueParticipantBatcher.getAssigneeParticipantsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		FollowerParticipantsForIssue:                dataloader.NewBatchedLoader(issueParticipantBatcher.getFollowerParticipantsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		CommentsForIssue:                            dataloader.NewBatchedLoader(commentBatcher.getCommentsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContractsForOrganization:                    dataloader.NewBatchedLoader(contractBatcher.getContractsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContractForInvoice:                          dataloader.NewBatchedLoader(contractBatcher.getContractsForInvoices, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ServiceLineItemsForContract:                 dataloader.NewBatchedLoader(serviceLineItemBatcher.getServiceLineItemsForContracts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ServiceLineItemForInvoiceLine:               dataloader.NewBatchedLoader(serviceLineItemBatcher.getServiceLineItemForInvoiceLine, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OpportunitiesForContract:                    dataloader.NewBatchedLoader(opportunityBatcher.getOpportunitiesForContracts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OpportunitiesForOrganization:                dataloader.NewBatchedLoader(opportunityBatcher.getOpportunitiesForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		InvoiceLinesForInvoice:                      dataloader.NewBatchedLoader(invoiceBatcher.getInvoiceLinesForInvoice, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		InvoicesForContract:                         dataloader.NewBatchedLoader(invoiceBatcher.getInvoicesForContract, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		FlowContactsForFlow:                         dataloader.NewBatchedLoader(flowBatcher.getFlowContactsForFlow, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		FlowActionsForFlow:                          dataloader.NewBatchedLoader(flowBatcher.getFlowActionsForFlow, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		FlowSendersForFlow:                          dataloader.NewBatchedLoader(flowBatcher.getFlowSendersForFlow, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		FlowsWithContact:                            dataloader.NewBatchedLoader(flowBatcher.getFlowsWithContact, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		FlowsWithSender:                             dataloader.NewBatchedLoader(flowBatcher.getFlowsWithSender, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
	}
}

// Middleware injects a DataLoader into the request context, so it can be used later in the schema resolvers
func Middleware(loaders *Loaders, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
