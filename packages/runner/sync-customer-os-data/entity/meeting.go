package entity

import (
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type MeetingData struct {
	Id                    string     `json:"id,omitempty"`
	Name                  string     `json:"name,omitempty"`
	CreatedAt             *time.Time `json:"createdAt,omitempty"`
	UpdatedAt             *time.Time `json:"updatedAt,omitempty"`
	StartedAt             *time.Time `json:"startedAt,omitempty"`
	EndedAt               *time.Time `json:"endedAt,omitempty"`
	ExternalId            string     `json:"externalId,omitempty"`
	ExternalSyncId        string     `json:"externalSyncId,omitempty"`
	ExternalSystem        string     `json:"externalSystem,omitempty"`
	Agenda                string     `json:"agenda,omitempty"`
	AgendaContentType     string     `json:"agendaContentType,omitempty"`
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

func (u *MeetingData) FormatTimes() {
	if u.CreatedAt != nil {
		u.CreatedAt = common_utils.TimePtr((*u.CreatedAt).UTC())
	}
	if u.UpdatedAt != nil {
		u.UpdatedAt = common_utils.TimePtr((*u.UpdatedAt).UTC())
	}
	if u.StartedAt != nil {
		u.StartedAt = common_utils.TimePtr((*u.StartedAt).UTC())
	}
	if u.EndedAt != nil {
		u.EndedAt = common_utils.TimePtr((*u.EndedAt).UTC())
	}
}
