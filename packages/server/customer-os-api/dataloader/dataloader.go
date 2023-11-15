package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"net/http"
	"time"
)

const defaultDataloaderWaitTime = 32 * time.Millisecond

type loadersString string

const loadersKey = loadersString("dataloaders")

type Loaders struct {
	TagsForOrganization                         *dataloader.Loader
	TagsForContact                              *dataloader.Loader
	TagsForIssue                                *dataloader.Loader
	TagsForLogEntry                             *dataloader.Loader
	EmailsForContact                            *dataloader.Loader
	EmailsForOrganization                       *dataloader.Loader
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
	NotedEntitiesForNote                        *dataloader.Loader
	UsersForEmail                               *dataloader.Loader
	UsersForPhoneNumber                         *dataloader.Loader
	UsersForPlayer                              *dataloader.Loader
	UserOwnerForOrganization                    *dataloader.Loader
	UserAuthorForLogEntry                       *dataloader.Loader
	UserAuthorForComment                        *dataloader.Loader
	User                                        *dataloader.Loader
	ContactsForEmail                            *dataloader.Loader
	ContactsForPhoneNumber                      *dataloader.Loader
	OrganizationsForEmail                       *dataloader.Loader
	OrganizationsForPhoneNumber                 *dataloader.Loader
	SubsidiariesForOrganization                 *dataloader.Loader
	SubsidiariesOfForOrganization               *dataloader.Loader
	SuggestedMergeToForOrganization             *dataloader.Loader
	DescribesForAnalysis                        *dataloader.Loader
	DescribedByFor                              *dataloader.Loader
	CreatedByParticipantsForMeeting             *dataloader.Loader
	AttendedByParticipantsForMeeting            *dataloader.Loader
	InteractionEventsForMeeting                 *dataloader.Loader
	InteractionEventsForIssue                   *dataloader.Loader
	CommentsForIssue                            *dataloader.Loader
	NotesForMeeting                             *dataloader.Loader
	AttachmentsForInteractionEvent              *dataloader.Loader
	AttachmentsForInteractionSession            *dataloader.Loader
	AttachmentsForMeeting                       *dataloader.Loader
	SocialsForContact                           *dataloader.Loader
	SocialsForOrganization                      *dataloader.Loader
	ExternalSystemsForComment                   *dataloader.Loader
	ExternalSystemsForIssue                     *dataloader.Loader
	ExternalSystemsForOrganization              *dataloader.Loader
	ExternalSystemsForLogEntry                  *dataloader.Loader
	ExternalSystemsForMeeting                   *dataloader.Loader
	ExternalSystemsForInteractionEvent          *dataloader.Loader
	TimelineEventForTimelineEventId             *dataloader.Loader
	OrganizationForJobRole                      *dataloader.Loader
	ContactForJobRole                           *dataloader.Loader
	IssueForInteractionEvent                    *dataloader.Loader
	MeetingForInteractionEvent                  *dataloader.Loader
	CountryForPhoneNumber                       *dataloader.Loader
	ActionItemsForInteractionEvent              *dataloader.Loader
	SubmitterParticipantsForIssue               *dataloader.Loader
	ReporterParticipantsForIssue                *dataloader.Loader
	AssigneeParticipantsForIssue                *dataloader.Loader
	FollowerParticipantsForIssue                *dataloader.Loader
	ContractsForOrganization                    *dataloader.Loader
	ServiceLineItemsForContract                 *dataloader.Loader
}

type tagBatcher struct {
	tagService service.TagService
}
type emailBatcher struct {
	emailService service.EmailService
}
type locationBatcher struct {
	locationService service.LocationService
}
type socialBatcher struct {
	socialService service.SocialService
}
type jobRoleBatcher struct {
	jobRoleService service.JobRoleService
}
type calendarBatcher struct {
	calendarService service.CalendarService
}
type domainBatcher struct {
	domainService service.DomainService
}
type interactionEventBatcher struct {
	interactionEventService service.InteractionEventService
}
type interactionSessionBatcher struct {
	interactionSessionService service.InteractionSessionService
}
type interactionEventParticipantBatcher struct {
	interactionEventService service.InteractionEventService
}
type interactionSessionParticipantBatcher struct {
	interactionSessionService service.InteractionSessionService
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
	organizationService service.OrganizationService
}
type analysisBatcher struct {
	analysisService service.AnalysisService
}
type noteBatcher struct {
	noteService service.NoteService
}
type attachmentBatcher struct {
	attachmentService service.AttachmentService
}
type externalSystemBatcher struct {
	externalSystemService service.ExternalSystemService
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
	serviceLineItemService service.ServiceLineItemService
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(services *service.Services) *Loaders {
	tagBatcher := &tagBatcher{
		tagService: services.TagService,
	}
	emailBatcher := &emailBatcher{
		emailService: services.EmailService,
	}
	locationBatcher := &locationBatcher{
		locationService: services.LocationService,
	}
	socialBatcher := &socialBatcher{
		socialService: services.SocialService,
	}
	jobRoleBatcher := &jobRoleBatcher{
		jobRoleService: services.JobRoleService,
	}
	calendarBatcher := &calendarBatcher{
		calendarService: services.CalendarService,
	}
	domainBatcher := &domainBatcher{
		domainService: services.DomainService,
	}
	interactionEventBatcher := &interactionEventBatcher{
		interactionEventService: services.InteractionEventService,
	}
	commentBatcher := &commentBatcher{
		commentService: services.CommentService,
	}
	interactionSessionBatcher := &interactionSessionBatcher{
		interactionSessionService: services.InteractionSessionService,
	}
	interactionEventParticipantBatcher := &interactionEventParticipantBatcher{
		interactionEventService: services.InteractionEventService,
	}
	issueParticipantBatcher := &issueParticipantBatcher{
		issueService: services.IssueService,
	}
	interactionSessionParticipantBatcher := &interactionSessionParticipantBatcher{
		interactionSessionService: services.InteractionSessionService,
	}
	meetingParticipantBatcher := &meetingParticipantBatcher{
		meetingService: services.MeetingService,
	}
	phoneNumberBatcher := &phoneNumberBatcher{
		phoneNumberService: services.PhoneNumberService,
	}
	notedEntityBatcher := &notedEntityBatcher{
		noteService: services.NoteService,
	}
	userBatcher := &userBatcher{
		userService: services.UserService,
	}
	contactBatcher := &contactBatcher{
		contactService: services.ContactService,
	}
	organizationBatcher := &organizationBatcher{
		organizationService: services.OrganizationService,
	}
	analysisBatcher := &analysisBatcher{
		analysisService: services.AnalysisService,
	}
	noteBatcher := &noteBatcher{
		noteService: services.NoteService,
	}
	attachmentBatcher := attachmentBatcher{
		attachmentService: services.AttachmentService,
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
	actionItemBatcher := actionItemBatcher{
		actionItemService: services.ActionItemService,
	}
	contractBatcher := &contractBatcher{
		contractService: services.ContractService,
	}
	serviceLineItemBatcher := &serviceLineItemBatcher{
		serviceLineItemService: services.ServiceLineItemService,
	}
	return &Loaders{
		TagsForOrganization:                         dataloader.NewBatchedLoader(tagBatcher.getTagsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		TagsForContact:                              dataloader.NewBatchedLoader(tagBatcher.getTagsForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		TagsForIssue:                                dataloader.NewBatchedLoader(tagBatcher.getTagsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		TagsForLogEntry:                             dataloader.NewBatchedLoader(tagBatcher.getTagsForLogEntries, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		EmailsForContact:                            dataloader.NewBatchedLoader(emailBatcher.getEmailsForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		EmailsForOrganization:                       dataloader.NewBatchedLoader(emailBatcher.getEmailsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
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
		NotedEntitiesForNote:                        dataloader.NewBatchedLoader(notedEntityBatcher.getNotedEntitiesForNotes, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UsersForEmail:                               dataloader.NewBatchedLoader(userBatcher.getUsersForEmails, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UsersForPhoneNumber:                         dataloader.NewBatchedLoader(userBatcher.getUsersForPhoneNumbers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UsersForPlayer:                              dataloader.NewBatchedLoader(userBatcher.getUsersForPlayers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserOwnerForOrganization:                    dataloader.NewBatchedLoader(userBatcher.getUserOwnersForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserAuthorForLogEntry:                       dataloader.NewBatchedLoader(userBatcher.getUserAuthorsForLogEntries, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		UserAuthorForComment:                        dataloader.NewBatchedLoader(userBatcher.getUserAuthorsForComments, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		User:                                        dataloader.NewBatchedLoader(userBatcher.getUsers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContactsForEmail:                            dataloader.NewBatchedLoader(contactBatcher.getContactsForEmails, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContactsForPhoneNumber:                      dataloader.NewBatchedLoader(contactBatcher.getContactsForPhoneNumbers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationsForEmail:                       dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForEmails, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationsForPhoneNumber:                 dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForPhoneNumbers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SubsidiariesForOrganization:                 dataloader.NewBatchedLoader(organizationBatcher.getSubsidiariesForOrganization, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SubsidiariesOfForOrganization:               dataloader.NewBatchedLoader(organizationBatcher.getSubsidiariesOfForOrganization, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(20*time.Millisecond), dataloader.WithWait(defaultDataloaderWaitTime)),
		SuggestedMergeToForOrganization:             dataloader.NewBatchedLoader(organizationBatcher.getSuggestedMergeToForOrganization, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		DescribesForAnalysis:                        dataloader.NewBatchedLoader(analysisBatcher.getDescribesForAnalysis, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		DescribedByFor:                              dataloader.NewBatchedLoader(analysisBatcher.getDescribedByFor, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		NotesForMeeting:                             dataloader.NewBatchedLoader(noteBatcher.getNotesForMeetings, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AttachmentsForInteractionEvent:              dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AttachmentsForInteractionSession:            dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForInteractionSessions, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AttachmentsForMeeting:                       dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForMeetings, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SocialsForContact:                           dataloader.NewBatchedLoader(socialBatcher.getSocialsForContacts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SocialsForOrganization:                      dataloader.NewBatchedLoader(socialBatcher.getSocialsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForComment:                   dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForComments, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForIssue:                     dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForOrganization:              dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForLogEntry:                  dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForLogEntries, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForMeeting:                   dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForMeetings, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ExternalSystemsForInteractionEvent:          dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		TimelineEventForTimelineEventId:             dataloader.NewBatchedLoader(timelineEventBatcher.getTimelineEventsForTimelineEventIds, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		OrganizationForJobRole:                      dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForJobRoles, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContactForJobRole:                           dataloader.NewBatchedLoader(contactBatcher.getContactsForJobRoles, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		IssueForInteractionEvent:                    dataloader.NewBatchedLoader(issueBatcher.getIssuesForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		MeetingForInteractionEvent:                  dataloader.NewBatchedLoader(meetingBatcher.getMeetingsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		CountryForPhoneNumber:                       dataloader.NewBatchedLoader(countryBatcher.getCountriesForPhoneNumbers, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ActionItemsForInteractionEvent:              dataloader.NewBatchedLoader(actionItemBatcher.getActionItemsForInteractionEvents, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		SubmitterParticipantsForIssue:               dataloader.NewBatchedLoader(issueParticipantBatcher.getSubmitterParticipantsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ReporterParticipantsForIssue:                dataloader.NewBatchedLoader(issueParticipantBatcher.getReporterParticipantsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		AssigneeParticipantsForIssue:                dataloader.NewBatchedLoader(issueParticipantBatcher.getAssigneeParticipantsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		FollowerParticipantsForIssue:                dataloader.NewBatchedLoader(issueParticipantBatcher.getFollowerParticipantsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		CommentsForIssue:                            dataloader.NewBatchedLoader(commentBatcher.getCommentsForIssues, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ContractsForOrganization:                    dataloader.NewBatchedLoader(contractBatcher.getContractsForOrganizations, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
		ServiceLineItemsForContract:                 dataloader.NewBatchedLoader(serviceLineItemBatcher.getServiceLineItemsForContracts, dataloader.WithClearCacheOnBatch(), dataloader.WithWait(defaultDataloaderWaitTime)),
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
