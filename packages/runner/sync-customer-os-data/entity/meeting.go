package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type MeetingData struct {
	BaseData
	Name                  string     `json:"name,omitempty"`
	StartedAtStr          string     `json:"startedAt,omitempty"`
	EndedAtStr            string     `json:"endedAt,omitempty"`
	StartedAt             *time.Time `json:"startedAtTime,omitempty"`
	EndedAt               *time.Time `json:"endedAtTime,omitempty"`
	Agenda                string     `json:"agenda,omitempty"`
	ContentType           string     `json:"contentType,omitempty"`
	MeetingUrl            string     `json:"meetingUrl,omitempty"`
	Location              string     `json:"location,omitempty"`
	ConferenceUrl         string     `json:"conferenceUrl,omitempty"`
	ExternalContactsIds   []string   `json:"externalContactsIds,omitempty"`
	CreatorUserExternalId string     `json:"externalUserId,omitempty"`
}

func (m *MeetingData) HasContacts() bool {
	return len(m.ExternalContactsIds) > 0
}

func (m *MeetingData) HasUserCreator() bool {
	return len(m.CreatorUserExternalId) > 0
}

func (m *MeetingData) HasLocation() bool {
	return len(m.Location) > 0
}

func (m *MeetingData) SetMeetingTimes() {
	if m.StartedAtStr != "" && m.StartedAt == nil {
		m.StartedAt, _ = utils.UnmarshalDateTime(m.StartedAtStr)
	}
	if m.StartedAt != nil {
		m.StartedAt = common_utils.TimePtr((*m.StartedAt).UTC())
	}

	if m.EndedAtStr != "" && m.EndedAt == nil {
		m.EndedAt, _ = utils.UnmarshalDateTime(m.EndedAtStr)
	}
	if m.EndedAt != nil {
		m.EndedAt = common_utils.TimePtr((*m.EndedAt).UTC())
	}
}

func (m *MeetingData) Normalize() {
	m.SetTimes()
	m.SetMeetingTimes()
}
