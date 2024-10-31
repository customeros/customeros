package contact

import (
	"fmt"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
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
	Username        string                        `json:"username"`
	Source          cmnmod.Source                 `json:"source"`
	CreatedAt       time.Time                     `json:"createdAt"`
	UpdatedAt       time.Time                     `json:"updatedAt"`
	PhoneNumbers    map[string]ContactPhoneNumber `json:"phoneNumbers"`
	// Deprecated
	LocationIds     []string                   `json:"locationIds,omitempty"`
	ExternalSystems []cmnmod.ExternalSystem    `json:"externalSystems"`
	TagIds          []string                   `json:"tagIds,omitempty"`
	Locations       map[string]cmnmod.Location `json:"locations,omitempty"`
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
	Primary bool `json:"primary"`
}

func (c *Contact) String() string {
	return fmt.Sprintf("Contact{ID: %s, FirstName: %s, LastName: %s, Prefix: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s}", c.ID, c.FirstName, c.LastName, c.Prefix, c.Source, c.CreatedAt, c.UpdatedAt)
}

func (c *Contact) HasPhoneNumber(phoneNumberId, label string, primary bool) bool {
	if len(c.PhoneNumbers) == 0 {
		return false
	}
	if email, ok := c.PhoneNumbers[phoneNumberId]; ok {
		return email.Label == label && email.Primary == primary
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
