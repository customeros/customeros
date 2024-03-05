package models

import (
	"fmt"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type PhoneNumberValidation struct {
	ValidationError string `json:"validationError"`
	SkipReason      string `json:"skipReason"`
}

type PhoneNumber struct {
	ID                    string                `json:"id"`
	RawPhoneNumber        string                `json:"rawPhoneNumber"`
	E164                  string                `json:"e164"`
	Source                cmnmod.Source         `json:"source"`
	CreatedAt             time.Time             `json:"createdAt"`
	UpdatedAt             time.Time             `json:"updatedAt"`
	PhoneNumberValidation PhoneNumberValidation `json:"phoneNumberValidation"`
	CountryCodeA2         string                `json:"countryCodeA2"`
}

func (p *PhoneNumber) String() string {
	return fmt.Sprintf("PhoneNumber{ID: %s, RawPhoneNumber: %s, E164: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s}", p.ID, p.RawPhoneNumber, p.E164, p.Source, p.CreatedAt, p.UpdatedAt)
}

func NewPhoneNumber() *PhoneNumber {
	return &PhoneNumber{}
}
