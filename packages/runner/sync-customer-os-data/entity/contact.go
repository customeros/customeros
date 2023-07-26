package entity

import (
	local_utils "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
	"time"
)

type ContactData struct {
	BaseData
	Prefix                        string            `json:"prefix,omitempty"`
	FirstName                     string            `json:"firstName,omitempty"`
	LastName                      string            `json:"LastName,omitempty"`
	Name                          string            `json:"name,omitempty"`
	Label                         string            `json:"label,omitempty"`
	JobTitle                      string            `json:"jobTitle,omitempty"`
	Notes                         []ContactNote     `json:"notes,omitempty"`
	ExternalUrl                   string            `json:"externalUrl,omitempty"`
	Email                         string            `json:"email,omitempty"`
	AdditionalEmails              []string          `json:"additionalEmails,omitempty"`
	PhoneNumber                   string            `json:"phoneNumber,omitempty"`
	AdditionalPhoneNumbers        []string          `json:"additionalPhoneNumbers,omitempty"`
	ExternalOrganizationsIds      []string          `json:"externalOrganizationsIds,omitempty"`
	PrimaryOrganizationExternalId string            `json:"externalOrganizationId,omitempty"`
	UserExternalOwnerId           string            `json:"externalOwnerId,omitempty"`
	UserExternalUserId            string            `json:"externalUserId,omitempty"`
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
	CreatedAtStr   string     `json:"createdAt,omitempty"`
	CreatedAt      *time.Time `json:"createdAtTime,omitempty"`
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
	return len(c.ExternalOrganizationsIds) > 0
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

func (c *ContactData) HasOwnerByOwnerId() bool {
	return c.UserExternalOwnerId != ""
}

func (c *ContactData) HasOwnerByUserId() bool {
	return c.UserExternalUserId != ""
}

func (c *ContactData) SetTextCustomFieldsTimes() {
	for i := range c.TextCustomFields {
		if c.TextCustomFields[i].CreatedAtStr != "" && c.TextCustomFields[i].CreatedAt == nil {
			c.TextCustomFields[i].CreatedAt, _ = local_utils.UnmarshalDateTime(c.TextCustomFields[i].CreatedAtStr)
		}
		if c.TextCustomFields[i].CreatedAt != nil {
			c.TextCustomFields[i].CreatedAt = utils.TimePtr((*c.TextCustomFields[i].CreatedAt).UTC())
		} else {
			c.TextCustomFields[i].CreatedAt = utils.TimePtr(utils.Now())
		}
	}
}

func (c *ContactData) Normalize() {
	c.SetTimes()

	c.ExternalOrganizationsIds = append(c.ExternalOrganizationsIds, c.PrimaryOrganizationExternalId)
	c.ExternalOrganizationsIds = utils.FilterEmpty(c.ExternalOrganizationsIds)
	c.ExternalOrganizationsIds = utils.RemoveDuplicates(c.ExternalOrganizationsIds)

	c.AdditionalEmails = utils.FilterEmpty(c.AdditionalEmails)
	c.AdditionalEmails = utils.RemoveDuplicates(c.AdditionalEmails)
	utils.LowercaseStrings(c.AdditionalEmails)
	c.Email = strings.ToLower(c.Email)

	c.AdditionalPhoneNumbers = utils.FilterEmpty(c.AdditionalPhoneNumbers)
	c.AdditionalPhoneNumbers = utils.RemoveDuplicates(c.AdditionalPhoneNumbers)
}
