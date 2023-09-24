package service

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/calendar/v3"
	"time"
)

type meetingService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type MeetingService interface {
	ReadNewCalendarEventsForUsername(gCalService *calendar.Service, tenant, username string) error
}

func (s *meetingService) ReadNewCalendarEventsForUsername(gCalService *calendar.Service, tenant, username string) error {
	calendarList, err := gCalService.CalendarList.List().Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve calendar list: %v", err)
	}

	for _, cal := range calendarList.Items {

		if cal.AccessRole == "owner" {

			importState, err := s.repositories.UserGCalImportStateRepository.GetGCalImportStateUsername(tenant, username, cal.Id)
			if err != nil {
				return fmt.Errorf("unable to retrieve history id for username: %v", err)
			}

			//incremental sync based on sync token
			//https://developers.google.com/calendar/api/guides/sync
			var syncToken string
			var pageToken string
			var timeMin time.Time
			var timeMax time.Time
			var maxResults int64

			if importState != nil {
				syncToken = importState.SyncToken
				pageToken = importState.PageToken
				maxResults = importState.MaxResults
				timeMin = importState.TimeMin
				timeMax = importState.TimeMax
			} else {
				syncToken = ""
				pageToken = ""
				maxResults = s.cfg.SyncData.GoogleCalendarReadBatchSize
				timeMin, err = time.Parse(time.DateOnly, s.cfg.SyncData.GoogleCalendarSyncStartDate)
				if err != nil {
					return fmt.Errorf("unable to parse start date: %v", err)
				}
				timeMax, err = time.Parse(time.DateOnly, s.cfg.SyncData.GoogleCalendarSyncStopDate)
				if err != nil {
					return fmt.Errorf("unable to parse start date: %v", err)
				}
			}

			getEventsWith := gCalService.Events.List(cal.Id)
			getEventsWith = getEventsWith.SingleEvents(true)
			getEventsWith = getEventsWith.MaxResults(maxResults)
			getEventsWith = getEventsWith.PageToken(pageToken)

			if syncToken != "" {
				logrus.Infof("Sync token is not empty. Performing incremental sync with token: %s", syncToken)
				getEventsWith = getEventsWith.SyncToken(syncToken)
			} else {
				getEventsWith = getEventsWith.TimeMin(timeMin.Format(time.RFC3339))
				getEventsWith = getEventsWith.TimeMax(timeMax.Format(time.RFC3339))
			}

			events, err := getEventsWith.Do()
			if err != nil {
				return fmt.Errorf("unable to retrieve events: %v", err)
			}
			if len(events.Items) > 0 {

				for _, item := range events.Items {

					if item.Status == "cancelled" {
						existingCalendarEvent, err := s.repositories.RawCalendarEventRepository.GetByProviderId("gcal", tenant, username, cal.Id, item.Id)
						if err != nil {
							return fmt.Errorf("failed to get existing calendar event: %v", err)
						}
						var calendarEventRawData CalendarEventRawData
						err = json.Unmarshal([]byte(existingCalendarEvent.Data), &calendarEventRawData)
						if err != nil {
							return fmt.Errorf("failed to unmarshal existing calendar event: %v", err)
						}
						calendarEventRawData.Status = "cancelled"

						jc, err := JSONMarshal(calendarEventRawData)
						if err != nil {
							return fmt.Errorf("failed to marshal calendar event content: %v", err)
						}
						existingCalendarEvent.Data = string(jc)

						err = s.repositories.RawCalendarEventRepository.Update(*existingCalendarEvent)
						if err != nil {
							return fmt.Errorf("failed to update existing calendar event: %v", err)
						}

						continue
					}

					if item.EventType != "default" {
						continue
					}

					var calendarEventRawData = CalendarEventRawData{
						Id:      item.Id,
						ICalUID: item.ICalUID,

						Created:           item.Created,
						Creator:           ConvertEventCreatorToRawData(item.Creator),
						Organizer:         ConvertEventOrganizerToRawData(item.Organizer),
						Attendees:         ConvertEventAttendeesToRawData(item.Attendees),
						EventType:         item.EventType,
						Description:       item.Description,
						OriginalStartTime: ConvertEventDateTimeToRawData(item.OriginalStartTime),
						Start:             ConvertEventDateTimeToRawData(item.Start),
						End:               ConvertEventDateTimeToRawData(item.End),
						HangoutLink:       item.HangoutLink,
						HtmlLink:          item.HtmlLink,
						Location:          item.Location,
						RecurringEventId:  item.RecurringEventId,
						Recurrence:        item.Recurrence,
						Status:            item.Status,
						Summary:           item.Summary,
						Transparency:      item.Transparency,
						Updated:           item.Updated,
						Visibility:        item.Visibility,
					}

					jsonContent, err := JSONMarshal(calendarEventRawData)
					if err != nil {
						return fmt.Errorf("failed to marshal calendar event content: %v", err)
					}

					err = s.repositories.RawCalendarEventRepository.SaveOrUpdate("gcal", tenant, username, cal.Id, calendarEventRawData.Id, calendarEventRawData.ICalUID, string(jsonContent))
					if err != nil {
						return fmt.Errorf("unable to store calendar event: %v", err)
					}

				}

			}

			if syncToken != events.NextSyncToken && events.NextSyncToken != "" {
				logrus.Infof("Sync token has changed. Updating sync token from %s to %s", syncToken, events.NextSyncToken)
				syncToken = events.NextSyncToken
			}

			err = s.repositories.UserGCalImportStateRepository.UpdateGCalImportStateForUsername(tenant, username, cal.Id, syncToken, events.NextPageToken, maxResults, timeMin, timeMax)
			if err != nil {
				return fmt.Errorf("unable to update the gcal page token for username: %v", err)
			}

		}

	}

	return nil

	//email, err := gmailService.Users.Messages.Get(username, providerMessageId).Format("full").Do()
	//if err != nil {
	//	return nil, fmt.Errorf("unable to retrieve email: %v", err)
	//}
	//
	//messageId := ""
	//emailSubject := ""
	//emailHtml := ""
	//emailText := ""
	//emailSentDate := ""
	//
	//from := ""
	//to := ""
	//cc := ""
	//bcc := ""
	//
	//threadId := email.ThreadId
	//references := ""
	//inReplyTo := ""
	//
	//emailHeaders := make(map[string]string)
	//
	//for i := range email.Payload.Headers {
	//	headerName := strings.ToLower(email.Payload.Headers[i].Name)
	//	emailHeaders[email.Payload.Headers[i].Name] = email.Payload.Headers[i].Value
	//	if headerName == "message-id" {
	//		messageId = email.Payload.Headers[i].Value
	//	} else if headerName == "subject" {
	//		emailSubject = email.Payload.Headers[i].Value
	//		if emailSubject == "" {
	//			emailSubject = "No Subject"
	//		} else if strings.HasPrefix(emailSubject, "Re: ") {
	//			emailSubject = emailSubject[4:]
	//		}
	//	} else if headerName == "from" {
	//		from = email.Payload.Headers[i].Value
	//	} else if headerName == "to" {
	//		to = email.Payload.Headers[i].Value
	//	} else if headerName == "cc" {
	//		cc = email.Payload.Headers[i].Value
	//	} else if headerName == "bcc" {
	//		bcc = email.Payload.Headers[i].Value
	//	} else if headerName == "references" {
	//		references = email.Payload.Headers[i].Value
	//	} else if headerName == "in-reply-to" {
	//		inReplyTo = email.Payload.Headers[i].Value
	//	} else if headerName == "date" {
	//		emailSentDate = email.Payload.Headers[i].Value
	//	}
	//}
	//
	//for i := range email.Payload.Parts {
	//	if email.Payload.Parts[i].MimeType == "text/html" {
	//		emailHtmlBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Body.Data)
	//		emailHtml = fmt.Sprintf("%s", emailHtmlBytes)
	//	} else if email.Payload.Parts[i].MimeType == "text/plain" {
	//		emailTextBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Body.Data)
	//		emailText = fmt.Sprintf("%s", string(emailTextBytes))
	//	} else if strings.HasPrefix(email.Payload.Parts[i].MimeType, "multipart") {
	//		for j := range email.Payload.Parts[i].Parts {
	//			if email.Payload.Parts[i].Parts[j].MimeType == "text/html" {
	//				emailHtmlBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Parts[j].Body.Data)
	//				emailHtml = fmt.Sprintf("%s", emailHtmlBytes)
	//			} else if email.Payload.Parts[i].Parts[j].MimeType == "text/plain" {
	//				emailTextBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Parts[j].Body.Data)
	//				emailText = fmt.Sprintf("%s", string(emailTextBytes))
	//			}
	//		}
	//	}
	//}
	//
	//rawEmailData := &EmailRawData{
	//	ProviderMessageId: providerMessageId,
	//	MessageId:         messageId,
	//	Sent:              emailSentDate,
	//	Subject:           emailSubject,
	//	From:              from,
	//	To:                to,
	//	Cc:                cc,
	//	Bcc:               bcc,
	//	Html:              emailHtml,
	//	Text:              emailText,
	//	ThreadId:          threadId,
	//	InReplyTo:         inReplyTo,
	//	Reference:         references,
	//	Headers:           emailHeaders,
	//}
	//
	//return rawEmailData, nil
	return nil
}

func ConvertEventDateTimeToRawData(eventDateTime *calendar.EventDateTime) *CalendarEventDateTime {
	if eventDateTime == nil {
		return nil
	}

	return &CalendarEventDateTime{
		DateTime: eventDateTime.DateTime,
		Date:     eventDateTime.Date,
		TimeZone: eventDateTime.TimeZone,
	}
}

func ConvertEventCreatorToRawData(eventCreator *calendar.EventCreator) *CalendarEventAttendee {
	if eventCreator == nil {
		return nil
	}

	return &CalendarEventAttendee{
		Id:             eventCreator.Id,
		Email:          eventCreator.Email,
		DisplayName:    eventCreator.DisplayName,
		Optional:       false,
		Organizer:      false,
		ResponseStatus: "accepted",
	}
}

func ConvertEventOrganizerToRawData(eventCreator *calendar.EventOrganizer) *CalendarEventAttendee {
	if eventCreator == nil {
		return nil
	}

	return &CalendarEventAttendee{
		Id:             eventCreator.Id,
		Email:          eventCreator.Email,
		DisplayName:    eventCreator.DisplayName,
		Optional:       false,
		Organizer:      true,
		ResponseStatus: "accepted",
	}
}

func ConvertEventAttendeesToRawData(eventAttendees []*calendar.EventAttendee) []*CalendarEventAttendee {
	if eventAttendees == nil {
		return nil
	}
	var result = make([]*CalendarEventAttendee, len(eventAttendees))
	for i := range eventAttendees {
		result[i] = ConvertEventAttendeeToRawData(eventAttendees[i])
	}
	return result
}

func ConvertEventAttendeeToRawData(eventCreator *calendar.EventAttendee) *CalendarEventAttendee {
	if eventCreator == nil {
		return nil
	}

	return &CalendarEventAttendee{
		Id:             eventCreator.Id,
		Email:          eventCreator.Email,
		DisplayName:    eventCreator.DisplayName,
		Optional:       eventCreator.Optional,
		Organizer:      eventCreator.Organizer,
		ResponseStatus: eventCreator.ResponseStatus,
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
