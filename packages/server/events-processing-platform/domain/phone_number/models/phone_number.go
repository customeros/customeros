package models

import (
	"fmt"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
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
	Source                commonModels.Source   `json:"source"`
	CreatedAt             time.Time             `json:"createdAt"`
	UpdatedAt             time.Time             `json:"updatedAt"`
	PhoneNumberValidation PhoneNumberValidation `json:"phoneNumberValidation"`
}

func (p *PhoneNumber) String() string {
	return fmt.Sprintf("PhoneNumber{ID: %s, RawPhoneNumber: %s, E164: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s}", p.ID, p.RawPhoneNumber, p.E164, p.Source, p.CreatedAt, p.UpdatedAt)
}

func NewPhoneNumber() *PhoneNumber {
	return &PhoneNumber{}
}
