package entity

import "time"

type OrganizationData struct {
	Id             string
	Name           string
	Description    string
	Domains        []string
	NoteContent    string
	Website        string
	Industry       string
	IsPublic       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExternalId     string
	ExternalSystem string

	DefaultLocationName string
	Country             string
	Region              string
	Locality            string
	Address             string
	Address2            string
	Zip                 string
	Phone               string

	OrganizationTypeName string

	ExternalSyncId string
}
