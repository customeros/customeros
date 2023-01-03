package entity

import "time"

type TextCustomField struct {
	Name   string
	Value  string
	Source string
}

type ContactData struct {
	Title     string
	FirstName string
	LastName  string
	Label     string
	JobTitle  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Readonly  bool

	ExternalId     string
	ExternalSystem string

	PrimaryEmail     string
	AdditionalEmails []string

	PrimaryE164 string

	OrganizationsExternalIds      []string
	PrimaryOrganizationExternalId string

	UserOwnerExternalId string

	TextCustomFields []TextCustomField

	ContactTypeName string

	Country string
	State   string
	City    string
	Address string
	Zip     string
	Fax     string
}
