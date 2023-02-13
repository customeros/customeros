package entity

import "time"

type TextCustomField struct {
	Name           string
	Value          string
	ExternalSystem string
}

type ContactData struct {
	Title     string
	FirstName string
	LastName  string
	Label     string
	JobTitle  string
	CreatedAt time.Time
	UpdatedAt time.Time

	ExternalId     string
	ExternalSystem string

	PrimaryEmail     string
	AdditionalEmails []string

	PrimaryE164 string

	OrganizationsExternalIds      []string
	PrimaryOrganizationExternalId string

	UserExternalOwnerId string

	TextCustomFields []TextCustomField

	TagName string

	DefaultLocationName string
	Country             string
	Region              string
	Locality            string
	Address             string
	Zip                 string
}
