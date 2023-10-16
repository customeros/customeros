package entity

import (
	"time"
)

type ContactData struct {
	BaseData
	Prefix               string                   `json:"prefix,omitempty"`
	FirstName            string                   `json:"firstName,omitempty"`
	LastName             string                   `json:"LastName,omitempty"`
	Name                 string                   `json:"name,omitempty"`
	Email                string                   `json:"email,omitempty"`
	AdditionalEmails     []string                 `json:"additionalEmails,omitempty"`
	PhoneNumbers         []PhoneNumber            `json:"phoneNumbers,omitempty"`
	UserExternalId       string                   `json:"externalUserId,omitempty"`
	UserExternalIdSecond string                   `json:"externalUserIdSecond,omitempty"`
	TextCustomFields     []TextCustomField        `json:"textCustomFields,omitempty"`
	Tags                 []string                 `json:"tags,omitempty"`
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
	OrganizationRequired bool                     `json:"organizationRequired,omitempty"`
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
