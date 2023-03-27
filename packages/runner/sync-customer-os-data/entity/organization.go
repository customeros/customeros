package entity

import "time"

type OrganizationRelation string

const (
	Subsidiary OrganizationRelation = "subsidiary"
)

type OrganizationNote struct {
	FieldSource string
	Note        string
}

type ParentOrganization struct {
	ExternalId           string
	OrganizationRelation OrganizationRelation
	Type                 string
}

type OrganizationData struct {
	Id             string
	Name           string
	Description    string
	Domains        []string
	Notes          []OrganizationNote
	Website        string
	Industry       string
	IsPublic       bool
	PhoneNumber    string
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExternalId     string
	ExternalSystem string
	ExternalUrl    string
	ExternalSyncId string

	DefaultLocationName string
	Country             string
	Region              string
	Locality            string
	Address             string
	Address2            string
	Zip                 string

	OrganizationTypeName string

	ParentOrganization *ParentOrganization
}

func (o OrganizationData) HasDomains() bool {
	return len(o.Domains) > 0
}

func (o OrganizationData) HasLocation() bool {
	return len(o.DefaultLocationName) > 0
}

func (o OrganizationData) HasNotes() bool {
	return len(o.Notes) > 0
}

func (o OrganizationData) HasOrganizationType() bool {
	return len(o.OrganizationTypeName) > 0
}

func (o OrganizationData) HasPhoneNumber() bool {
	return len(o.PhoneNumber) > 0
}

func (o OrganizationData) HasEmail() bool {
	return len(o.Email) > 0
}

func (o OrganizationData) IsSubsidiary() bool {
	return o.ParentOrganization != nil && o.ParentOrganization.OrganizationRelation == Subsidiary
}
