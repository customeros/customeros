package models

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
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
	Source                 commonmodel.Source            `json:"source"`
	CreatedAt              time.Time                     `json:"createdAt"`
	UpdatedAt              time.Time                     `json:"updatedAt"`
	PhoneNumbers           map[string]ContactPhoneNumber `json:"phoneNumbers"`
	Emails                 map[string]ContactEmail       `json:"emails"`
	Locations              []string                      `json:"locations,omitempty"`
	ExternalSystems        []commonmodel.ExternalSystem  `json:"externalSystems"`
	JobRolesByOrganization map[string]JobRole            `json:"jobRoles,omitempty"`
}

type JobRole struct {
	JobTitle    string             `json:"jobTitle"`
	Description string             `json:"description"`
	Primary     bool               `json:"primary"`
	StartedAt   *time.Time         `json:"startedAt"`
	EndedAt     *time.Time         `json:"endedAt"`
	CreatedAt   time.Time          `json:"createdAt"`
	Source      commonmodel.Source `json:"source"`
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

func (c *Contact) SameData(fields ContactDataFields, externalSystem commonmodel.ExternalSystem) bool {
	if externalSystem.Available() && !c.HasExternalSystem(externalSystem) {
		return false
	}
	if c.Name == fields.Name &&
		c.FirstName == fields.FirstName &&
		c.LastName == fields.LastName &&
		c.Prefix == fields.Prefix &&
		c.Description == fields.Description &&
		c.Timezone == fields.Timezone &&
		c.ProfilePhotoUrl == fields.ProfilePhotoUrl {
		return true
	}
	return false
}

func (c *Contact) HasExternalSystem(externalSystem commonmodel.ExternalSystem) bool {
	for _, es := range c.ExternalSystems {
		if es.ExternalSystemId == externalSystem.ExternalSystemId &&
			es.ExternalId == externalSystem.ExternalId &&
			es.ExternalSource == externalSystem.ExternalSource &&
			es.ExternalUrl == externalSystem.ExternalUrl &&
			es.ExternalIdSecond == externalSystem.ExternalIdSecond {
			return true
		}
	}
	return false
}

func (c *Contact) HasLocation(locationId string) bool {
	for _, location := range c.Locations {
		if location == locationId {
			return true
		}
	}
	return false
}

func (c *Contact) HasJobRoleInOrganization(organizationId string, jobRoleFields JobRole) bool {
	if c.JobRolesByOrganization == nil {
		return false
	}
	if jobRoles, ok := c.JobRolesByOrganization[organizationId]; ok {
		return jobRoles.JobTitle == jobRoleFields.JobTitle &&
			jobRoles.Description == jobRoleFields.Description &&
			jobRoles.Primary == jobRoleFields.Primary &&
			utils.IsEqualTimePtr(jobRoles.StartedAt, jobRoleFields.StartedAt) &&
			utils.IsEqualTimePtr(jobRoles.EndedAt, jobRoleFields.EndedAt)
	}
	return false
}
