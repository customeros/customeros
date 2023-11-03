package models

import (
	"fmt"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Contact struct {
	ID                     string                        `json:"id"`
	FirstName              string                        `json:"firstName"`
	LastName               string                        `json:"lastName"`
	Name                   string                        `json:"name"`
	Prefix                 string                        `json:"prefix"`
	Description            string                        `json:"description"`
	Timezone               string                        `json:"timezone"`
	ProfilePhotoUrl        string                        `json:"profilePhotoUrl"`
	Source                 cmnmod.Source                 `json:"source"`
	CreatedAt              time.Time                     `json:"createdAt"`
	UpdatedAt              time.Time                     `json:"updatedAt"`
	PhoneNumbers           map[string]ContactPhoneNumber `json:"phoneNumbers"`
	Emails                 map[string]ContactEmail       `json:"emails"`
	Locations              []string                      `json:"locations,omitempty"`
	ExternalSystems        []cmnmod.ExternalSystem       `json:"externalSystems"`
	JobRolesByOrganization map[string]JobRole            `json:"jobRoles,omitempty"`
}

type JobRole struct {
	JobTitle    string        `json:"jobTitle"`
	Description string        `json:"description"`
	Primary     bool          `json:"primary"`
	StartedAt   *time.Time    `json:"startedAt"`
	EndedAt     *time.Time    `json:"endedAt"`
	CreatedAt   time.Time     `json:"createdAt"`
	Source      cmnmod.Source `json:"source"`
}

type ContactPhoneNumber struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

type ContactEmail struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

func (c *Contact) String() string {
	return fmt.Sprintf("Contact{ID: %s, FirstName: %s, LastName: %s, Prefix: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s}", c.ID, c.FirstName, c.LastName, c.Prefix, c.Source, c.CreatedAt, c.UpdatedAt)
}

func (c *Contact) HasEmail(emailId, label string, primary bool) bool {
	if len(c.Emails) == 0 {
		return false
	}
	if email, ok := c.Emails[emailId]; ok {
		return email.Label == label && email.Primary == primary
	}
	return false
}

func (c *Contact) HasPhoneNumber(phoneNumberId, label string, primary bool) bool {
	if len(c.Emails) == 0 {
		return false
	}
	if email, ok := c.PhoneNumbers[phoneNumberId]; ok {
		return email.Label == label && email.Primary == primary
	}
	return false
}
