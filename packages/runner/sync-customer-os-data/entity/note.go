package entity

import "time"

type NoteData struct {
	Id                       string
	Html                     string
	CreatedAt                time.Time
	UserExternalId           string
	UserExternalOwnerId      string
	ContactsExternalIds      []string
	OrganizationsExternalIds []string
	ExternalId               string
	ExternalSystem           string
}
