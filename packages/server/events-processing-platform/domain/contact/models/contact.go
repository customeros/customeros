package models

import (
	"fmt"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type Contact struct {
	ID           string                `json:"id"`
	FirstName    string                `json:"firstName"`
	LastName     string                `json:"lastName"`
	Name         string                `json:"name"`
	Prefix       string                `json:"prefix"`
	Source       commonModels.Source   `json:"source"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
	PhoneNumbers []ContactPhoneNumbers `json:"phoneNumbers"`
}

type ContactPhoneNumbers struct {
	PhoneNumberID string `json:"id"`
	Primary       bool   `json:"primary"`
	Label         string `json:"label"`
}

func (contact *Contact) String() string {
	return fmt.Sprintf("Contact{ID: %s, FirstName: %s, LastName: %s, Name: %s, Prefix: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s}", contact.ID, contact.FirstName, contact.LastName, contact.Name, contact.Prefix, contact.Source, contact.CreatedAt, contact.UpdatedAt)
}

func NewContact() *Contact {
	return &Contact{}
}
