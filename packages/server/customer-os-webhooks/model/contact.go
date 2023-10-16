package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
	"time"
)

type ContactData struct {
	BaseData
	Prefix               string                   `json:"prefix,omitempty"`
	FirstName            string                   `json:"firstName,omitempty"`
	LastName             string                   `json:"lastName,omitempty"`
	Name                 string                   `json:"name,omitempty"`
	Description          string                   `json:"description,omitempty"`
	Email                string                   `json:"email,omitempty"`
	AdditionalEmails     []string                 `json:"additionalEmails,omitempty"`
	PhoneNumbers         []PhoneNumber            `json:"phoneNumbers,omitempty"`
	UserExternalId       string                   `json:"externalUserId,omitempty"`       // TODO not used in webhooks yet
	UserExternalIdSecond string                   `json:"externalUserIdSecond,omitempty"` // TODO not used in webhooks yet to be changed with referenced user
	TextCustomFields     []TextCustomField        `json:"textCustomFields,omitempty"`     // TODO not used in webhooks yet
	Tags                 []string                 `json:"tags,omitempty"`                 // TODO not used in webhooks yet
	LocationName         string                   `json:"locationName,omitempty"`
	Country              string                   `json:"country,omitempty"`
	Region               string                   `json:"region,omitempty"`
	Locality             string                   `json:"locality,omitempty"`
	Street               string                   `json:"street,omitempty"`
	Address              string                   `json:"address,omitempty"`
	Zip                  string                   `json:"zip,omitempty"`
	PostalCode           string                   `json:"postalCode,omitempty"`
	Timezone             string                   `json:"timezone,omitempty"`
	ProfilePhotoUrl      string                   `json:"profilePhotoUrl,omitempty"`
	Organizations        []ReferencedOrganization `json:"organizations,omitempty"`
	// if no valid organization provided sync is skipped
	OrganizationRequired bool `json:"organizationRequired,omitempty"`
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

func (c *ContactData) AllEmailAddresses() []string {
	var allEmails []string
	if len(c.Email) > 0 {
		allEmails = append(allEmails, c.Email)
	}
	if len(c.AdditionalEmails) > 0 {
		allEmails = append(allEmails, c.AdditionalEmails...)
	}
	return allEmails
}

func (c *ContactData) HasPhoneNumbers() bool {
	return len(c.PhoneNumbers) > 0
}

func (c *ContactData) PrimaryPhoneNumber() string {
	for _, phoneNumber := range c.PhoneNumbers {
		if phoneNumber.Primary {
			return phoneNumber.Number
		}
	}
	return ""
}

func (c *ContactData) HasOrganizations() bool {
	found := false
	for _, org := range c.Organizations {
		if org.Available() {
			found = true
			break
		}
	}
	return found
}

func (c *ContactData) HasPrimaryEmail() bool {
	return c.Email != ""
}

func (c *ContactData) HasAdditionalEmails() bool {
	return len(c.AdditionalEmails) > 0
}

func (c *ContactData) HasLocation() bool {
	return c.LocationName != "" || c.Country != "" || c.Region != "" || c.Locality != "" || c.Address != "" || c.Zip != "" || c.PostalCode != "" || c.Street != ""
}

func (c *ContactData) HasTextCustomFields() bool {
	return len(c.TextCustomFields) > 0
}

func (c *ContactData) HasTags() bool {
	return len(c.Tags) > 0
}

func (c *ContactData) SetTextCustomFieldsTimes() {
	for i := range c.TextCustomFields {
		if c.TextCustomFields[i].CreatedAtStr != "" && c.TextCustomFields[i].CreatedAt == nil {
			c.TextCustomFields[i].CreatedAt, _ = utils.UnmarshalDateTime(c.TextCustomFields[i].CreatedAtStr)
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
	c.BaseData.Normalize()

	c.AdditionalEmails = utils.FilterEmpty(c.AdditionalEmails)
	c.AdditionalEmails = utils.RemoveDuplicates(c.AdditionalEmails)
	utils.LowercaseStrings(c.AdditionalEmails)
	c.Email = strings.ToLower(c.Email)
	c.AdditionalEmails = utils.RemoveFromList(c.AdditionalEmails, c.Email)

	for _, phoneNumber := range c.PhoneNumbers {
		phoneNumber.Normalize()
	}
	c.PhoneNumbers = GetNonEmptyPhoneNumbers(c.PhoneNumbers)
	c.PhoneNumbers = RemoveDuplicatedPhoneNumbers(c.PhoneNumbers)
}
