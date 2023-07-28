package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"net/http"
)

type loadersString string

const loadersKey = loadersString("dataloaders")

type Loaders struct {
	TagsForOrganization                         *dataloader.Loader
	TagsForContact                              *dataloader.Loader
	TagsForIssue                                *dataloader.Loader
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
	MentionedEntitiesForNote                    *dataloader.Loader
	UsersForEmail                               *dataloader.Loader
	UsersForPhoneNumber                         *dataloader.Loader
	UsersForPlayer                              *dataloader.Loader
	UserOwnerForOrganization                    *dataloader.Loader
	ContactsForEmail                            *dataloader.Loader
	ContactsForPhoneNumber                      *dataloader.Loader
	OrganizationsForEmail                       *dataloader.Loader
	OrganizationsForPhoneNumber                 *dataloader.Loader
	SubsidiariesForOrganization                 *dataloader.Loader
	SubsidiariesOfForOrganization               *dataloader.Loader
	DescribesForAnalysis                        *dataloader.Loader
	DescribedByForMeeting                       *dataloader.Loader
	DescribedByForInteractionSession            *dataloader.Loader
	CreatedByParticipantsForMeeting             *dataloader.Loader
	AttendedByParticipantsForMeeting            *dataloader.Loader
	InteractionEventsForMeeting                 *dataloader.Loader
	InteractionEventsForIssue                   *dataloader.Loader
	MentionedByNotesForIssue                    *dataloader.Loader
	NotesForMeeting                             *dataloader.Loader
	AttachmentsForInteractionEvent              *dataloader.Loader
	AttachmentsForInteractionSession            *dataloader.Loader
	AttachmentsForMeeting                       *dataloader.Loader
	SocialsForContact                           *dataloader.Loader
	SocialsForOrganization                      *dataloader.Loader
	RelationshipsForOrganization                *dataloader.Loader
	RelationshipStagesForOrganization           *dataloader.Loader
	ExternalSystemsForEntity                    *dataloader.Loader
	TimelineEventForTimelineEventId             *dataloader.Loader
	OrganizationForJobRole                      *dataloader.Loader
	IssueForInteractionEvent                    *dataloader.Loader
	MeetingForInteractionEvent                  *dataloader.Loader
	HealthIndicatorForOrganization              *dataloader.Loader
	CountryForPhoneNumber                       *dataloader.Loader
	ActionItemsForInteractionEvent              *dataloader.Loader
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
type mentionedEntityBatcher struct {
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
type relationshipBatcher struct {
	organizationRelationshipService service.OrganizationRelationshipService
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
type healthIndicatorBatcher struct {
	healthIndicatorService service.HealthIndicatorService
}
type countryBatcher struct {
	countryService service.CountryService
}
type actionItemBatcher struct {
	actionItemService service.ActionItemService
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
	interactionSessionBatcher := &interactionSessionBatcher{
		interactionSessionService: services.InteractionSessionService,
	}
	interactionEventParticipantBatcher := &interactionEventParticipantBatcher{
		interactionEventService: services.InteractionEventService,
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
	mentionedEntityBatcher := &mentionedEntityBatcher{
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
	relationshipBatcher := relationshipBatcher{
		organizationRelationshipService: services.OrganizationRelationshipService,
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
	healthIndicatorBatcher := healthIndicatorBatcher{
		healthIndicatorService: services.HealthIndicatorService,
	}
	countryBatcher := countryBatcher{
		countryService: services.CountryService,
	}
	actionItemBatcher := actionItemBatcher{
		actionItemService: services.ActionItemService,
	}
	return &Loaders{
		TagsForOrganization:                         dataloader.NewBatchedLoader(tagBatcher.getTagsForOrganizations, dataloader.WithClearCacheOnBatch()),
		TagsForContact:                              dataloader.NewBatchedLoader(tagBatcher.getTagsForContacts, dataloader.WithClearCacheOnBatch()),
		TagsForIssue:                                dataloader.NewBatchedLoader(tagBatcher.getTagsForIssues, dataloader.WithClearCacheOnBatch()),
		EmailsForContact:                            dataloader.NewBatchedLoader(emailBatcher.getEmailsForContacts, dataloader.WithClearCacheOnBatch()),
		EmailsForOrganization:                       dataloader.NewBatchedLoader(emailBatcher.getEmailsForOrganizations, dataloader.WithClearCacheOnBatch()),
		LocationsForContact:                         dataloader.NewBatchedLoader(locationBatcher.getLocationsForContacts, dataloader.WithClearCacheOnBatch()),
		LocationsForOrganization:                    dataloader.NewBatchedLoader(locationBatcher.getLocationsForOrganizations, dataloader.WithClearCacheOnBatch()),
		JobRolesForContact:                          dataloader.NewBatchedLoader(jobRoleBatcher.getJobRolesForContacts, dataloader.WithClearCacheOnBatch()),
		JobRolesForOrganization:                     dataloader.NewBatchedLoader(jobRoleBatcher.getJobRolesForOrganizations, dataloader.WithClearCacheOnBatch()),
		JobRolesForUser:                             dataloader.NewBatchedLoader(jobRoleBatcher.getJobRolesForUsers, dataloader.WithClearCacheOnBatch()),
		CalendarsForUser:                            dataloader.NewBatchedLoader(calendarBatcher.getCalendarsForUsers, dataloader.WithClearCacheOnBatch()),
		DomainsForOrganization:                      dataloader.NewBatchedLoader(domainBatcher.getDomainsForOrganizations, dataloader.WithClearCacheOnBatch()),
		InteractionEventsForInteractionSession:      dataloader.NewBatchedLoader(interactionEventBatcher.getInteractionEventsForInteractionSessions, dataloader.WithClearCacheOnBatch()),
		InteractionEventsForMeeting:                 dataloader.NewBatchedLoader(interactionEventBatcher.getInteractionEventsForMeetings, dataloader.WithClearCacheOnBatch()),
		InteractionEventsForIssue:                   dataloader.NewBatchedLoader(interactionEventBatcher.getInteractionEventsForIssues, dataloader.WithClearCacheOnBatch()),
		InteractionSessionForInteractionEvent:       dataloader.NewBatchedLoader(interactionSessionBatcher.getInteractionSessionsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		SentByParticipantsForInteractionEvent:       dataloader.NewBatchedLoader(interactionEventParticipantBatcher.getSentByParticipantsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		SentToParticipantsForInteractionEvent:       dataloader.NewBatchedLoader(interactionEventParticipantBatcher.getSentToParticipantsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		AttendedByParticipantsForInteractionSession: dataloader.NewBatchedLoader(interactionSessionParticipantBatcher.getAttendedByParticipantsForInteractionSessions, dataloader.WithClearCacheOnBatch()),
		CreatedByParticipantsForMeeting:             dataloader.NewBatchedLoader(meetingParticipantBatcher.getCreatedByParticipantsForMeeting, dataloader.WithClearCacheOnBatch()),
		AttendedByParticipantsForMeeting:            dataloader.NewBatchedLoader(meetingParticipantBatcher.getAttendedByParticipantsForMeeting, dataloader.WithClearCacheOnBatch()),
		PhoneNumbersForOrganization:                 dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForOrganizations, dataloader.WithClearCacheOnBatch()),
		PhoneNumbersForUser:                         dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForUsers, dataloader.WithClearCacheOnBatch()),
		PhoneNumbersForContact:                      dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForContacts, dataloader.WithClearCacheOnBatch()),
		ReplyToInteractionEventForInteractionEvent:  dataloader.NewBatchedLoader(interactionEventBatcher.getReplyToInteractionEventsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		NotedEntitiesForNote:                        dataloader.NewBatchedLoader(notedEntityBatcher.getNotedEntitiesForNotes, dataloader.WithClearCacheOnBatch()),
		MentionedEntitiesForNote:                    dataloader.NewBatchedLoader(mentionedEntityBatcher.getMentionedEntitiesForNotes, dataloader.WithClearCacheOnBatch()),
		UsersForEmail:                               dataloader.NewBatchedLoader(userBatcher.getUsersForEmails, dataloader.WithClearCacheOnBatch()),
		UsersForPhoneNumber:                         dataloader.NewBatchedLoader(userBatcher.getUsersForPhoneNumbers, dataloader.WithClearCacheOnBatch()),
		UsersForPlayer:                              dataloader.NewBatchedLoader(userBatcher.getUsersForPlayers, dataloader.WithClearCacheOnBatch()),
		UserOwnerForOrganization:                    dataloader.NewBatchedLoader(userBatcher.getUserOwnersForOrganizations, dataloader.WithClearCacheOnBatch()),
		ContactsForEmail:                            dataloader.NewBatchedLoader(contactBatcher.getContactsForEmails, dataloader.WithClearCacheOnBatch()),
		ContactsForPhoneNumber:                      dataloader.NewBatchedLoader(contactBatcher.getContactsForPhoneNumbers, dataloader.WithClearCacheOnBatch()),
		OrganizationsForEmail:                       dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForEmails, dataloader.WithClearCacheOnBatch()),
		OrganizationsForPhoneNumber:                 dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForPhoneNumbers, dataloader.WithClearCacheOnBatch()),
		SubsidiariesForOrganization:                 dataloader.NewBatchedLoader(organizationBatcher.getSubsidiariesForOrganization, dataloader.WithClearCacheOnBatch()),
		SubsidiariesOfForOrganization:               dataloader.NewBatchedLoader(organizationBatcher.getSubsidiariesOfForOrganization, dataloader.WithClearCacheOnBatch()),
		DescribesForAnalysis:                        dataloader.NewBatchedLoader(analysisBatcher.getDescribesForAnalysis, dataloader.WithClearCacheOnBatch()),
		MentionedByNotesForIssue:                    dataloader.NewBatchedLoader(noteBatcher.getMentionedByNotesForIssue, dataloader.WithClearCacheOnBatch()),
		DescribedByForInteractionSession:            dataloader.NewBatchedLoader(analysisBatcher.getDescribedByForInteractionSession, dataloader.WithClearCacheOnBatch()),
		DescribedByForMeeting:                       dataloader.NewBatchedLoader(analysisBatcher.getDescribedByForMeeting, dataloader.WithClearCacheOnBatch()),
		NotesForMeeting:                             dataloader.NewBatchedLoader(noteBatcher.getNotesForMeetings, dataloader.WithClearCacheOnBatch()),
		AttachmentsForInteractionEvent:              dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		AttachmentsForInteractionSession:            dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForInteractionSessions, dataloader.WithClearCacheOnBatch()),
		AttachmentsForMeeting:                       dataloader.NewBatchedLoader(attachmentBatcher.getAttachmentsForMeetings, dataloader.WithClearCacheOnBatch()),
		SocialsForContact:                           dataloader.NewBatchedLoader(socialBatcher.getSocialsForContacts, dataloader.WithClearCacheOnBatch()),
		SocialsForOrganization:                      dataloader.NewBatchedLoader(socialBatcher.getSocialsForOrganizations, dataloader.WithClearCacheOnBatch()),
		RelationshipsForOrganization:                dataloader.NewBatchedLoader(relationshipBatcher.getRelationshipsForOrganizations, dataloader.WithClearCacheOnBatch()),
		RelationshipStagesForOrganization:           dataloader.NewBatchedLoader(relationshipBatcher.getRelationshipStagesForOrganizations, dataloader.WithClearCacheOnBatch()),
		ExternalSystemsForEntity:                    dataloader.NewBatchedLoader(externalSystemBatcher.getExternalSystemsForEntities, dataloader.WithClearCacheOnBatch()),
		TimelineEventForTimelineEventId:             dataloader.NewBatchedLoader(timelineEventBatcher.getTimelineEventsForTimelineEventIds, dataloader.WithClearCacheOnBatch()),
		OrganizationForJobRole:                      dataloader.NewBatchedLoader(organizationBatcher.getOrganizationsForJobRoles, dataloader.WithClearCacheOnBatch()),
		IssueForInteractionEvent:                    dataloader.NewBatchedLoader(issueBatcher.getIssuesForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		MeetingForInteractionEvent:                  dataloader.NewBatchedLoader(meetingBatcher.getMeetingsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		HealthIndicatorForOrganization:              dataloader.NewBatchedLoader(healthIndicatorBatcher.getHealthIndicatorsForOrganizations, dataloader.WithClearCacheOnBatch()),
		CountryForPhoneNumber:                       dataloader.NewBatchedLoader(countryBatcher.getCountriesForPhoneNumbers, dataloader.WithClearCacheOnBatch()),
		ActionItemsForInteractionEvent:              dataloader.NewBatchedLoader(actionItemBatcher.getActionItemsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
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
