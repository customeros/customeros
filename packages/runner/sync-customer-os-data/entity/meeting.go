package entity

import "time"

type MeetingData struct {
	Id                 string
	Name               string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	StartedAt          time.Time
	EndedAt            time.Time
	ExternalId         string
	ExternalSyncId     string
	ExternalSystem     string
	Agenda             string
	AgendaContentType  string
	MeetingExternalUrl string
	Location           string
	ConferenceUrl      string

	ContactsExternalIds   []string
	UserCreatorExternalId string
}

func (m MeetingData) HasContacts() bool {
	return len(m.ContactsExternalIds) > 0
}

func (m MeetingData) HasUserCreator() bool {
	return len(m.UserCreatorExternalId) > 0
}

func (m MeetingData) HasLocation() bool {
	return len(m.Location) > 0
}
