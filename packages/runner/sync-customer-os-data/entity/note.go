package entity

import "time"

type NoteData struct {
	Id                  string
	Html                string
	CreatedAt           time.Time
	UserExternalId      string
	ContactsExternalIds []string
	ExternalId          string
	ExternalSystem      string
}
