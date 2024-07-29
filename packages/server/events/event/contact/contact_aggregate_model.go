package contact

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"reflect"
	"time"
)

type Contact struct {
	ID              string                        `json:"id"`
	FirstName       string                        `json:"firstName"`
	LastName        string                        `json:"lastName"`
	Name            string                        `json:"name"`
	Prefix          string                        `json:"prefix"`
	Description     string                        `json:"description"`
	Timezone        string                        `json:"timezone"`
	ProfilePhotoUrl string                        `json:"profilePhotoUrl"`
	Source          cmnmod.Source                 `json:"source"`
	CreatedAt       time.Time                     `json:"createdAt"`
	UpdatedAt       time.Time                     `json:"updatedAt"`
	PhoneNumbers    map[string]ContactPhoneNumber `json:"phoneNumbers"`
	Emails          map[string]ContactEmail       `json:"emails"`
	Socials         map[string]cmnmod.Social      `json:"socials,omitempty"`
	// Deprecated
	LocationIds            []string                   `json:"locationIds,omitempty"`
	ExternalSystems        []cmnmod.ExternalSystem    `json:"externalSystems"`
	JobRolesByOrganization map[string]JobRole         `json:"jobRoles,omitempty"`
	TagIds                 []string                   `json:"tagIds,omitempty"`
	Locations              map[string]cmnmod.Location `json:"locations,omitempty"`
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

func (c *Contact) SameData(fields event.ContactDataFields, externalSystem cmnmod.ExternalSystem) bool {
	if !externalSystem.Available() {
		return false
	}

	if externalSystem.Available() && !c.HasExternalSystem(externalSystem) {
		return false
	}

	if c.Source.SourceOfTruth == externalSystem.ExternalSystemId {
		if c.Name == fields.Name &&
			c.FirstName == fields.FirstName &&
			c.LastName == fields.LastName &&
			c.Prefix == fields.Prefix &&
			c.Description == fields.Description &&
			c.Timezone == fields.Timezone &&
			c.ProfilePhotoUrl == fields.ProfilePhotoUrl {
			return true
		}
	} else {
		if (c.Name != "" || c.Name == fields.Name) &&
			(c.FirstName != "" || c.FirstName == fields.FirstName) &&
			(c.LastName != "" || c.LastName == fields.LastName) &&
			(c.Prefix != "" || c.Prefix == fields.Prefix) &&
			(c.Description != "" || c.Description == fields.Description) &&
			(c.Timezone != "" || c.Timezone == fields.Timezone) &&
			(c.ProfilePhotoUrl != "" || c.ProfilePhotoUrl == fields.ProfilePhotoUrl) {
			return true
		}
	}

	return false
}

func (c *Contact) HasExternalSystem(externalSystem cmnmod.ExternalSystem) bool {
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
	for _, location := range c.LocationIds {
		if location == locationId {
			return true
		}
	}
	return false
}

func (c *Contact) HasSocialUrl(url string) bool {
	if c.Socials == nil {
		return false
	}
	for _, social := range c.Socials {
		if social.Url == url {
			return true
		}
	}
	return false
}

func (c *Contact) GetSocialIdForUrl(url string) string {
	if c.Socials == nil {
		return ""
	}
	for key, social := range c.Socials {
		if social.Url == url {
			return key
		}
	}
	return ""
}

func (c *Contact) HasJobRoleInOrganization(organizationId string, jobRoleFields JobRole, sourceFields cmnmod.Source) bool {
	if c.JobRolesByOrganization == nil {
		return false
	}
	if jobRole, ok := c.JobRolesByOrganization[organizationId]; ok {
		found := jobRole.JobTitle == jobRoleFields.JobTitle &&
			jobRole.Description == jobRoleFields.Description &&
			jobRole.Primary == jobRoleFields.Primary &&
			(utils.IsEqualTimePtr(jobRole.StartedAt, jobRoleFields.StartedAt) || jobRoleFields.StartedAt == nil) &&
			(utils.IsEqualTimePtr(jobRole.EndedAt, jobRoleFields.EndedAt) || jobRoleFields.EndedAt == nil)
		if found {
			return true
		}
		if sourceFields.Source != jobRole.Source.SourceOfTruth && jobRole.Source.SourceOfTruth == events2.SourceOpenline {
			return !(jobRole.JobTitle == "" && jobRoleFields.JobTitle != "") &&
				!(jobRole.Description == "" && jobRoleFields.Description != "") &&
				!(jobRole.StartedAt == nil && jobRoleFields.StartedAt != nil) &&
				!(jobRole.EndedAt == nil && jobRoleFields.EndedAt != nil)
		}
	}
	return false
}

func (c *Contact) GetLocationIdForDetails(location cmnmod.Location) string {
	for id, orgLocation := range c.Locations {
		if locationMatchesExcludingName(orgLocation, location) {
			return id
		}
	}
	return ""
}

func locationMatchesExcludingName(contactLocation, inputLocation cmnmod.Location) bool {
	// Create copies of the locations to avoid modifying the original structs
	contactCopy := contactLocation
	inputCopy := inputLocation

	// Set Name to empty string for both locations to exclude it from comparison
	contactCopy.Name = ""
	inputCopy.Name = ""

	// Compare all fields except Name
	return reflect.DeepEqual(contactCopy, inputCopy)
}
