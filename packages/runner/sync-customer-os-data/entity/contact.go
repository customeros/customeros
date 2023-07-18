package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ContactData struct {
	BaseData
	Prefix                        string            `json:"prefix,omitempty"`
	FirstName                     string            `json:"firstName,omitempty"`
	LastName                      string            `json:"LastName,omitempty"`
	Label                         string            `json:"label,omitempty"`
	JobTitle                      string            `json:"jobTitle,omitempty"`
	Notes                         []ContactNote     `json:"notes,omitempty"`
	ExternalUrl                   string            `json:"externalUrl,omitempty"`
	Email                         string            `json:"email,omitempty"`
	AdditionalEmails              []string          `json:"additionalEmails,omitempty"`
	PhoneNumber                   string            `json:"phoneNumber,omitempty"`
	OrganizationsExternalIds      []string          `json:"organizationsExternalIds,omitempty"`
	PrimaryOrganizationExternalId string            `json:"externalOrganizationId,omitempty"`
	UserExternalOwnerId           string            `json:"externalOwnerId,omitempty"`
	TextCustomFields              []TextCustomField `json:"textCustomFields,omitempty"`
	Tags                          []string          `json:"tags,omitempty"`
	Location                      string            `json:"location,omitempty"`
	Country                       string            `json:"country,omitempty"`
	Region                        string            `json:"region,omitempty"`
	Locality                      string            `json:"locality,omitempty"`
	Address                       string            `json:"address,omitempty"`
	Zip                           string            `json:"zip,omitempty"`
}

type ContactNote struct {
	FieldSource string `json:"fieldSource,omitempty"`
	Note        string `json:"note,omitempty"`
}

type TextCustomField struct {
	Name           string     `json:"name,omitempty"`
	Value          string     `json:"value,omitempty"`
	ExternalSystem string     `json:"externalSystem,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
}

func (c *ContactData) EmailsForUnicity() []string {
	var emailsForUnicity []string
	if len(c.Email) > 0 {
		emailsForUnicity = append(emailsForUnicity, c.Email)
	} else if len(c.AdditionalEmails) == 1 {
		emailsForUnicity = append(emailsForUnicity, c.AdditionalEmails...)
	}
	return emailsForUnicity
}

func (c *ContactData) AllEmails() []string {
	var allEmails []string
	if len(c.Email) > 0 {
		allEmails = append(allEmails, c.Email)
	}
	if len(c.AdditionalEmails) > 0 {
		allEmails = append(allEmails, c.AdditionalEmails...)
	}
	return allEmails
}

func (c *ContactData) HasPhoneNumber() bool {
	return len(c.PhoneNumber) > 0
}

func (c *ContactData) HasOrganizations() bool {
	return len(c.OrganizationsExternalIds) > 0
}

func (c *ContactData) HasNotes() bool {
	return len(c.Notes) > 0
}

func (c *ContactData) HasLocation() bool {
	return len(c.Location) > 0 || len(c.Country) > 0 || len(c.Region) > 0 || len(c.Locality) > 0 || len(c.Address) > 0 || len(c.Zip) > 0
}

func (c *ContactData) HasTextCustomFields() bool {
	return len(c.TextCustomFields) > 0
}

func (c *ContactData) HasTags() bool {
	return len(c.Tags) > 0
}

func (c *ContactData) HasOwner() bool {
	return len(c.UserExternalOwnerId) > 0
}

func (c *ContactData) FormatTimes() {
	if c.CreatedAt != nil {
		c.CreatedAt = utils.TimePtr((*c.CreatedAt).UTC())
	} else {
		c.CreatedAt = utils.TimePtr(utils.Now())
	}
	if c.UpdatedAt != nil {
		c.UpdatedAt = utils.TimePtr((*c.UpdatedAt).UTC())
	} else {
		c.UpdatedAt = utils.TimePtr(utils.Now())
	}
	for i := range c.TextCustomFields {
		if c.TextCustomFields[i].CreatedAt != nil {
			c.TextCustomFields[i].CreatedAt = utils.TimePtr((*c.TextCustomFields[i].CreatedAt).UTC())
		} else {
			c.TextCustomFields[i].CreatedAt = utils.TimePtr(utils.Now())
		}
	}
}

func (c *ContactData) Normalize() {
	c.FormatTimes()

	c.OrganizationsExternalIds = append(c.OrganizationsExternalIds, c.PrimaryOrganizationExternalId)
	c.OrganizationsExternalIds = utils.FilterEmpty(c.OrganizationsExternalIds)
	c.OrganizationsExternalIds = utils.RemoveDuplicates(c.OrganizationsExternalIds)

	c.AdditionalEmails = utils.FilterEmpty(c.AdditionalEmails)
}
