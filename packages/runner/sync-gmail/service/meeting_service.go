package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"time"
)

const GCalSource = "gcal"

type meetingService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type MeetingService interface {
	SyncCalendarEvents(externalSystemId, tenant string, personalEmailProviderList []commonEntity.PersonalEmailProvider, organizationAllowedForImport []commonEntity.WhitelistDomain)
}

func (s *meetingService) SyncCalendarEvents(externalSystemId, tenant string, personalEmailProviderList []commonEntity.PersonalEmailProvider, organizationAllowedForImport []commonEntity.WhitelistDomain) {
	calendarEventsIdsForSync, err := s.repositories.RawCalendarEventRepository.GetCalendarEventsIdsForSync(externalSystemId, tenant)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
	}

	s.syncCalendarEvents(externalSystemId, tenant, calendarEventsIdsForSync, personalEmailProviderList, organizationAllowedForImport)
}

func (s *meetingService) syncCalendarEvents(externalSystemId, tenant string, calendarEvents []entity.RawCalendarEvent, personalEmailProviderList []commonEntity.PersonalEmailProvider, organizationAllowedForImport []commonEntity.WhitelistDomain) {
	for _, calendarEvent := range calendarEvents {
		state, reason, err := s.syncCalendarEvent(externalSystemId, tenant, calendarEvent.ID, personalEmailProviderList, organizationAllowedForImport)

		var errMessage *string
		if err != nil {
			s2 := err.Error()
			errMessage = &s2
		}

		err = s.repositories.RawCalendarEventRepository.MarkSentToEventStore(calendarEvent.ID, state, reason, errMessage)
		if err != nil {
			logrus.Errorf("unable to mark raw calendar event as sent to event store: %v", err)
		}

		fmt.Println("raw calendar event processed: " + calendarEvent.ID.String())
	}
}

func (s *meetingService) syncCalendarEvent(externalSystemId, tenant string, rawCalendarId uuid.UUID, personalEmailProviderList []commonEntity.PersonalEmailProvider, whitelistDomainList []commonEntity.WhitelistDomain) (entity.RawState, *string, error) {
	ctx := context.Background()

	rawCalendarIdString := rawCalendarId.String()

	calendarEvent, err := s.repositories.RawCalendarEventRepository.GetCalendarEventForSync(rawCalendarId)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
		return entity.ERROR, nil, err
	}

	rawCalendarEventData := CalendarEventRawData{}
	err = json.Unmarshal([]byte(calendarEvent.Data), &rawCalendarEventData)
	if err != nil {
		logrus.Errorf("failed to unmarshal raw calendar event data: %v", err)
		return entity.ERROR, nil, err
	}

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Neo4jDriver)
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		logrus.Errorf("failed to start transaction for calendar event with id %v: %v", rawCalendarIdString, err)
		return entity.ERROR, nil, err
	}

	existingMeetingNode, err := s.repositories.MeetingRepository.GetByExternalId(ctx, tx, tenant, externalSystemId, calendarEvent.ProviderId)
	if err != nil {
		logrus.Errorf("failed to check if meeting exists for external id %v for tenant %v :%v", calendarEvent.ProviderId, tenant, err)
		return entity.ERROR, nil, err
	}

	if existingMeetingNode != nil {
		//todo update / delete
		reason := "implement update / delete"
		return entity.SKIPPED, &reason, nil
	} else {

		now := time.Now().UTC()

		createdAt, err := s.services.SyncService.ConvertToUTC(rawCalendarEventData.Created)
		if err != nil {
			logrus.Errorf("failed to convert created date to utc for email with id %v: %v", rawCalendarIdString, err)
			return entity.ERROR, nil, err
		}

		startedAt, err := getDate(rawCalendarEventData.Start)
		if err != nil {
			logrus.Errorf("failed to convert start date to utc for email with id %v: %v", rawCalendarIdString, err)
			return entity.ERROR, nil, err
		}

		endedAt, err := getDate(rawCalendarEventData.End)
		if err != nil {
			logrus.Errorf("failed to convert end date to utc for email with id %v: %v", rawCalendarIdString, err)
			return entity.ERROR, nil, err
		}

		status := getMeetingStatus(rawCalendarEventData.Status)
		contentType := "text/html"
		meetingForCustomerOS := entity.MeetingEntity{
			Name:               &rawCalendarEventData.Summary,
			CreatedAt:          createdAt,
			UpdatedAt:          createdAt,
			StartedAt:          startedAt,
			EndedAt:            endedAt,
			ConferenceUrl:      &rawCalendarEventData.HangoutLink,
			MeetingExternalUrl: &rawCalendarEventData.HangoutLink, //what is this and what's the difference with ConfereceUrl
			Agenda:             &rawCalendarEventData.Description,
			AgendaContentType:  &contentType,
			AppSource:          AppSource,
			Source:             GCalSource,
			SourceOfTruth:      "openline",
			Status:             &status,
		}

		meetingNode, err := s.repositories.MeetingRepository.Create(ctx, tx, tenant, externalSystemId, rawCalendarEventData.Id, &meetingForCustomerOS, now)
		if err != nil {
			logrus.Errorf("failed merge meeting for raw calendar id %v :%v", rawCalendarIdString, err)
			return entity.ERROR, nil, err
		}
		meetingId := utils.GetStringPropOrNil(meetingNode.Props, "id")

		//link meeting with creator
		creatorEmailId, err := s.GetAttendeeEmailIdAndType(tx, tenant, *meetingId, rawCalendarEventData.Creator.Email, whitelistDomainList, personalEmailProviderList, now)
		if err != nil {
			logrus.Errorf("failed to get creator email id for raw email id %v :%v", rawCalendarIdString, err)
			return entity.ERROR, nil, err
		}
		err = s.repositories.MeetingRepository.LinkWithEmailInTx(ctx, tx, tenant, *meetingId, *creatorEmailId, entity.CREATED_BY)
		if err != nil {
			logrus.Errorf("failed to link creator with meeting for raw email id %v :%v", rawCalendarIdString, err)
			return entity.ERROR, nil, err
		}

		//link meeting with attendees
		for _, attendee := range rawCalendarEventData.Attendees {

			attendeeEmailId, err := s.GetAttendeeEmailIdAndType(tx, tenant, *meetingId, attendee.Email, whitelistDomainList, personalEmailProviderList, now)
			if err != nil {
				logrus.Errorf("failed to get attendee email id for raw email id %v :%v", rawCalendarIdString, err)
				return entity.ERROR, nil, err
			}
			err = s.repositories.MeetingRepository.LinkWithEmailInTx(ctx, tx, tenant, *meetingId, *attendeeEmailId, entity.ATTENDED_BY)
			if err != nil {
				logrus.Errorf("failed to link attendee with meeting for raw email id %v :%v", rawCalendarIdString, err)
				return entity.ERROR, nil, err
			}

		}

		//organizationIdList, err := s.repositories.OrganizationRepository.GetOrganizationsLinkedToEmailsInTx(ctx, tx, tenant, emailidList)
		//if err != nil {
		//	logrus.Errorf("unable to retrieve organization id list for tenant: %v", err)
		//	return entity.ERROR, nil, err
		//}
		//
		//for _, organizationId := range organizationIdList {
		//	lastTouchpointAt, lastTouchpointId, err := s.repositories.TimelineEventRepository.CalculateAndGetLastTouchpointInTx(ctx, tx, tenant, organizationId)
		//	if err != nil {
		//		logrus.Errorf("unable to calculate last touchpoint for organization: %v", err)
		//		return entity.ERROR, nil, err
		//	}
		//
		//	if lastTouchpointAt != nil && lastTouchpointId != "" {
		//		err := s.repositories.OrganizationRepository.UpdateLastTouchpointInTx(ctx, tx, tenant, organizationId, *lastTouchpointAt, lastTouchpointId)
		//		if err != nil {
		//			logrus.Errorf("unable to update last touchpoint for organization: %v", err)
		//			return entity.ERROR, nil, err
		//		}
		//	}
		//}

		err = tx.Commit(ctx)
		if err != nil {
			logrus.Errorf("failed to commit transaction: %v", err)
			return entity.ERROR, nil, err
		}

	}

	return entity.SENT, nil, err
}

func (s *meetingService) GetAttendeeEmailIdAndType(tx neo4j.ManagedTransaction, tenant, meetingId, emailAddress string, whitelistDomainList []commonEntity.WhitelistDomain, personalEmailProviderList []commonEntity.PersonalEmailProvider, now time.Time) (*string, error) {
	ctx := context.Background()

	fromEmailId, err := s.repositories.EmailRepository.GetEmailIdInTx(ctx, tx, tenant, emailAddress)
	if err != nil {
		logrus.Errorf("failed to get email id for email %v: %v", emailAddress, err)
		return nil, err
	}
	if fromEmailId == "" {
		fromEmailId, err = s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, meetingId, emailAddress, s.services.SyncService.GetWhitelistedDomain(extractDomain(emailAddress), whitelistDomainList), personalEmailProviderList, now, GCalSource)
		if err != nil {
			logrus.Errorf("unable to retrieve email id for tenant: %v", err)
			return nil, err
		}
	}

	return &fromEmailId, nil
}

func getDate(calendarDate *CalendarEventDateTime) (*time.Time, error) {
	if calendarDate == nil {
		return nil, nil
	}

	if calendarDate.Date != "" {
		layouts := []string{
			time.DateOnly,
		}
		var parsedTime time.Time
		var err error

		// Try parsing with each layout until successful
		for _, layout := range layouts {
			parsedTime, err = time.Parse(layout, calendarDate.Date)
			if err == nil {
				break
			}
		}

		if parsedTime.IsZero() {
			return nil, errors.New(fmt.Sprintf("unable to parse date: %v", calendarDate.Date))
		} else {
			parsedTime = parsedTime.UTC()
		}

		return &parsedTime, nil
	} else if calendarDate.DateTime != "" {
		// Parse the dateTime string into a time.Time value
		parsedTime, err := time.Parse(time.RFC3339, calendarDate.DateTime)
		if err != nil {
			return nil, err
		}

		if parsedTime.IsZero() {
			return nil, errors.New(fmt.Sprintf("unable to parse date: %v", calendarDate.DateTime))
		}

		if calendarDate.TimeZone != "" {
			// Set the time zone for the parsed time
			location, err := time.LoadLocation(calendarDate.TimeZone)
			if err != nil {
				return nil, err
			}
			parsedTime = parsedTime.In(location)
		}

		parsedTime = parsedTime.UTC()

		return &parsedTime, nil
	}

	return nil, nil
}

func getMeetingStatus(calendarEventStatus string) entity.MeetingStatus {
	switch calendarEventStatus {
	case "confirmed":
		return entity.MeetingStatusAccepted
	case "cancelled":
		return entity.MeetingStatusCanceled
	default:
		return entity.MeetingStatusUndefined
	}
}

type CalendarEventAttendee struct {

	// Id: The creator's Profile ID, if available.
	Id string `json:"id,omitempty"`

	// Email: The creator's email address, if available.
	Email string `json:"email,omitempty"`

	// DisplayName: The creator's name, if available.
	DisplayName string `json:"displayName,omitempty"`

	// Optional: Whether this is an optional attendee. Optional. The default
	// is False.
	Optional bool `json:"optional,omitempty"`

	// Organizer: Whether the attendee is the organizer of the event.
	// Read-only. The default is False.
	Organizer bool `json:"organizer,omitempty"`

	// ResponseStatus: The attendee's response status. Possible values are:
	//
	// - "needsAction" - The attendee has not responded to the invitation
	// (recommended for new events).
	// - "declined" - The attendee has declined the invitation.
	// - "tentative" - The attendee has tentatively accepted the invitation.
	//
	// - "accepted" - The attendee has accepted the invitation.  Warning: If
	// you add an event using the values declined, tentative, or accepted,
	// attendees with the "Add invitations to my calendar" setting set to
	// "When I respond to invitation in email" won't see an event on their
	// calendar unless they choose to change their invitation response in
	// the event invitation email.
	ResponseStatus string `json:"responseStatus,omitempty"`
}

type CalendarEventDateTime struct {
	// Date: The date, in the format "yyyy-mm-dd", if this is an all-day
	// event.
	Date string `json:"date,omitempty"`

	// DateTime: The time, as a combined date-time value (formatted
	// according to RFC3339). A time zone offset is required unless a time
	// zone is explicitly specified in timeZone.
	DateTime string `json:"dateTime,omitempty"`

	// TimeZone: The time zone in which the time is specified. (Formatted as
	// an IANA Time Zone Database name, e.g. "Europe/Zurich".) For recurring
	// events this field is required and specifies the time zone in which
	// the recurrence is expanded. For single events this field is optional
	// and indicates a custom time zone for the event start/end.
	TimeZone string `json:"timeZone,omitempty"`
}

type CalendarEventRawData struct {
	// Attendees: The attendees of the event. See the Events with attendees
	// guide for more information on scheduling events with other calendar
	// users. Service accounts need to use domain-wide delegation of
	// authority to populate the attendee list.
	Attendees []*CalendarEventAttendee `json:"attendees,omitempty"`

	// ConferenceData: The conference-related information, such as details
	// of a Google Meet conference. To create new conference details use the
	// createRequest field. To persist your changes, remember to set the
	// conferenceDataVersion request parameter to 1 for all event
	// modification requests.
	//ConferenceData *CalendarConferenceData `json:"conferenceData,omitempty"`

	// Created: Creation time of the event (as a RFC3339 timestamp).
	// Read-only.
	Created string `json:"created,omitempty"`

	// Creator: The creator of the event. Read-only.
	Creator *CalendarEventAttendee `json:"creator,omitempty"`

	// Organizer: The organizer of the event. If the organizer is also an
	// attendee, this is indicated with a separate entry in attendees with
	// the organizer field set to True. To change the organizer, use the
	// move operation. Read-only, except when importing an event.
	Organizer *CalendarEventAttendee `json:"organizer,omitempty"`

	// EventType: Specific type of the event. Read-only. Possible values
	// are:
	// - "default" - A regular event or not further specified.
	// - "outOfOffice" - An out-of-office event.
	// - "focusTime" - A focus-time event.
	// - "workingLocation" - A working location event. Developer Preview.
	EventType string `json:"eventType,omitempty"`

	// Description: Description of the event. Can contain HTML. Optional.
	Description string `json:"description,omitempty"`

	// OriginalStartTime: For an instance of a recurring event, this is the
	// time at which this event would start according to the recurrence data
	// in the recurring event identified by recurringEventId. It uniquely
	// identifies the instance within the recurring event series even if the
	// instance was moved to a different time. Immutable.
	OriginalStartTime *CalendarEventDateTime `json:"originalStartTime,omitempty"`

	// Start: The (inclusive) start time of the event. For a recurring
	// event, this is the start time of the first instance.
	Start *CalendarEventDateTime `json:"start,omitempty"`

	// End: The (exclusive) end time of the event. For a recurring event,
	// this is the end time of the first instance.
	End *CalendarEventDateTime `json:"end,omitempty"`

	// HangoutLink: An absolute link to the Google Hangout associated with
	// this event. Read-only.
	HangoutLink string `json:"hangoutLink,omitempty"`

	// HtmlLink: An absolute link to this event in the Google Calendar Web
	// UI. Read-only.
	HtmlLink string `json:"htmlLink,omitempty"`

	// ICalUID: Event unique identifier as defined in RFC5545. It is used to
	// uniquely identify events accross calendaring systems and must be
	// supplied when importing events via the import method.
	// Note that the iCalUID and the id are not identical and only one of
	// them should be supplied at event creation time. One difference in
	// their semantics is that in recurring events, all occurrences of one
	// event have different ids while they all share the same iCalUIDs. To
	// retrieve an event using its iCalUID, call the events.list method
	// using the iCalUID parameter. To retrieve an event using its id, call
	// the events.get method.
	ICalUID string `json:"iCalUID,omitempty"`

	// Id: Opaque identifier of the event. When creating new single or
	// recurring events, you can specify their IDs. Provided IDs must follow
	// these rules:
	// - characters allowed in the ID are those used in base32hex encoding,
	// i.e. lowercase letters a-v and digits 0-9, see section 3.1.2 in
	// RFC2938
	// - the length of the ID must be between 5 and 1024 characters
	// - the ID must be unique per calendar  Due to the globally distributed
	// nature of the system, we cannot guarantee that ID collisions will be
	// detected at event creation time. To minimize the risk of collisions
	// we recommend using an established UUID algorithm such as one
	// described in RFC4122.
	// If you do not specify an ID, it will be automatically generated by
	// the server.
	// Note that the icalUID and the id are not identical and only one of
	// them should be supplied at event creation time. One difference in
	// their semantics is that in recurring events, all occurrences of one
	// event have different ids while they all share the same icalUIDs.
	Id string `json:"id,omitempty"`

	// Location: Geographic location of the event as free-form text.
	// Optional.
	Location string `json:"location,omitempty"`

	// Recurrence: List of RRULE, EXRULE, RDATE and EXDATE lines for a
	// recurring event, as specified in RFC5545. Note that DTSTART and DTEND
	// lines are not allowed in this field; event start and end times are
	// specified in the start and end fields. This field is omitted for
	// single events or instances of recurring events.
	Recurrence []string `json:"recurrence,omitempty"`

	// RecurringEventId: For an instance of a recurring event, this is the
	// id of the recurring event to which this instance belongs. Immutable.
	RecurringEventId string `json:"recurringEventId,omitempty"`

	// Status: Status of the event. Optional. Possible values are:
	// - "confirmed" - The event is confirmed. This is the default status.
	//
	// - "tentative" - The event is tentatively confirmed.
	// - "cancelled" - The event is cancelled (deleted). The list method
	// returns cancelled events only on incremental sync (when syncToken or
	// updatedMin are specified) or if the showDeleted flag is set to true.
	// The get method always returns them.
	// A cancelled status represents two different states depending on the
	// event type:
	// - Cancelled exceptions of an uncancelled recurring event indicate
	// that this instance should no longer be presented to the user. Clients
	// should store these events for the lifetime of the parent recurring
	// event.
	// Cancelled exceptions are only guaranteed to have values for the id,
	// recurringEventId and originalStartTime fields populated. The other
	// fields might be empty.
	// - All other cancelled events represent deleted events. Clients should
	// remove their locally synced copies. Such cancelled events will
	// eventually disappear, so do not rely on them being available
	// indefinitely.
	// Deleted events are only guaranteed to have the id field populated.
	// On the organizer's calendar, cancelled events continue to expose
	// event details (summary, location, etc.) so that they can be restored
	// (undeleted). Similarly, the events to which the user was invited and
	// that they manually removed continue to provide details. However,
	// incremental sync requests with showDeleted set to false will not
	// return these details.
	// If an event changes its organizer (for example via the move
	// operation) and the original organizer is not on the attendee list, it
	// will leave behind a cancelled event where only the id field is
	// guaranteed to be populated.
	Status string `json:"status,omitempty"`

	// Summary: Title of the event.
	Summary string `json:"summary,omitempty"`

	// Transparency: Whether the event blocks time on the calendar.
	// Optional. Possible values are:
	// - "opaque" - Default value. The event does block time on the
	// calendar. This is equivalent to setting Show me as to Busy in the
	// Calendar UI.
	// - "transparent" - The event does not block time on the calendar. This
	// is equivalent to setting Show me as to Available in the Calendar UI.
	Transparency string `json:"transparency,omitempty"`

	// Updated: Last modification time of the event (as a RFC3339
	// timestamp). Read-only.
	Updated string `json:"updated,omitempty"`

	// Visibility: Visibility of the event. Optional. Possible values are:
	//
	// - "default" - Uses the default visibility for events on the calendar.
	// This is the default value.
	// - "public" - The event is public and event details are visible to all
	// readers of the calendar.
	// - "private" - The event is private and only event attendees may view
	// event details.
	// - "confidential" - The event is private. This value is provided for
	// compatibility reasons.
	Visibility string `json:"visibility,omitempty"`
}

func NewMeetingService(cfg *config.Config, repositories *repository.Repositories, services *Services) MeetingService {
	return &meetingService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
