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
	PhoneNumber    string
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

	OrganizationTypeName string

	ExternalSyncId string
}

func (o OrganizationData) HasDomains() bool {
	return len(o.Domains) > 0
}

func (o OrganizationData) HasLocation() bool {
	return len(o.DefaultLocationName) > 0
}

func (o OrganizationData) HasNotes() bool {
	return len(o.NoteContent) > 0
}

func (o OrganizationData) HasOrganizationType() bool {
	return len(o.OrganizationTypeName) > 0
}

func (o OrganizationData) HasPhoneNumber() bool {
	return len(o.PhoneNumber) > 0
}
