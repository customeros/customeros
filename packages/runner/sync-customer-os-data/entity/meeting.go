package entity

import (
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

/*
{
  "name": "Project Review",
  "startedAt": "2023-03-01T14:00:00Z",
  "endedAt": "2023-03-01T15:00:00Z",
  "agenda": "Review project status and milestones",
  "contentType": "text/plain",
  "meetingUrl": "https://meetings.com/meeting/1234",
  "location": "Board Room A",
  "conferenceUrl": "https://meetings.com/dialin/1234",
  "contactsExternalIds": [
    "contact-123",
    "contact-456"
  ],
  "externalUserId": "user-123",

  "skip": false,
  "skipReason": "draft data",
  "id": "1234",
  "externalId": "abcd1234",
  "externalSystem": "HubSpot",
  "createdAt": "2022-02-28T19:52:05Z",
  "updatedAt": "2022-03-01T11:23:45Z",
  "syncId": "sync_1234"
}
*/

type MeetingData struct {
	BaseData
	Name                  string     `json:"name,omitempty"`
	StartedAt             *time.Time `json:"startedAt,omitempty"`
	EndedAt               *time.Time `json:"endedAt,omitempty"`
	Agenda                string     `json:"agenda,omitempty"`
	ContentType           string     `json:"contentType,omitempty"`
	MeetingUrl            string     `json:"meetingUrl,omitempty"`
	Location              string     `json:"location,omitempty"`
	ConferenceUrl         string     `json:"conferenceUrl,omitempty"`
	ContactsExternalIds   []string   `json:"contactsExternalIds,omitempty"`
	CreatorUserExternalId string     `json:"externalUserId,omitempty"`
}

func (m *MeetingData) HasContacts() bool {
	return len(m.ContactsExternalIds) > 0
}

func (m *MeetingData) HasUserCreator() bool {
	return len(m.CreatorUserExternalId) > 0
}

func (m *MeetingData) HasLocation() bool {
	return len(m.Location) > 0
}

func (m *MeetingData) FormatTimes() {
	if m.CreatedAt != nil {
		m.CreatedAt = common_utils.TimePtr((*m.CreatedAt).UTC())
	} else {
		m.CreatedAt = common_utils.TimePtr(common_utils.Now())
	}
	if m.UpdatedAt != nil {
		m.UpdatedAt = common_utils.TimePtr((*m.UpdatedAt).UTC())
	} else {
		m.UpdatedAt = common_utils.TimePtr(common_utils.Now())
	}
	if m.StartedAt != nil {
		m.StartedAt = common_utils.TimePtr((*m.StartedAt).UTC())
	}
	if m.EndedAt != nil {
		m.EndedAt = common_utils.TimePtr((*m.EndedAt).UTC())
	}
}

func (m *MeetingData) Normalize() {
	m.FormatTimes()
}
