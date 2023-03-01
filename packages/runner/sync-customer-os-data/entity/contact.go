package entity

import "time"

type ContactNote struct {
	FieldSource string
	Note        string
}

type ContactData struct {
	Id        string
	Title     string
	FirstName string
	LastName  string
	Label     string
	JobTitle  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Notes     []ContactNote

	ExternalId     string
	ExternalSystem string
	ExternalUrl    string
	ExternalSyncId string

	PrimaryEmail     string
	AdditionalEmails []string

	PhoneNumber string

	OrganizationsExternalIds      []string
	PrimaryOrganizationExternalId string

	UserExternalOwnerId string

	TextCustomFields []TextCustomField
	Tags             []string

	DefaultLocationName string
	Country             string
	Region              string
	Locality            string
	Address             string
	Zip                 string
}

func (c ContactData) AllEmails() []string {
	var allEmails []string
	if len(c.PrimaryEmail) > 0 {
		allEmails = append(allEmails, c.PrimaryEmail)
	}
	if len(c.AdditionalEmails) > 0 {
		allEmails = append(allEmails, c.AdditionalEmails...)
	}
	return allEmails
}

func (c ContactData) HasPhoneNumber() bool {
	return len(c.PhoneNumber) > 0
}

func (c ContactData) HasOrganizations() bool {
	return len(c.OrganizationsExternalIds) > 0
}

func (c ContactData) HasNotes() bool {
	return len(c.Notes) > 0
}

func (c ContactData) HasLocation() bool {
	return len(c.DefaultLocationName) > 0
}

func (c ContactData) HasTextCustomFields() bool {
	return len(c.TextCustomFields) > 0
}

func (c ContactData) HasTags() bool {
	return len(c.Tags) > 0
}
