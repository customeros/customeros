package dto

import "time"

type AddContactToOrganization struct {
	ContactId      string     `json:"contactId"`
	OrganizationId string     `json:"organizationId"`
	JobTitle       string     `json:"jobTitle"`
	Description    string     `json:"description"`
	Primary        bool       `json:"primary"`
	StartedAt      *time.Time `json:"startedAt,omitempty"`
	EndedAt        *time.Time `json:"endedAt,omitempty"`
}
