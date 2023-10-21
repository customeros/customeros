package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertPhoneNumberCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	RawPhoneNumber  string
	Source          cmnmod.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

type FailedPhoneNumberValidationCommand struct {
	eventstore.BaseCommand
	RawPhoneNumber  string
	ValidationError string
	CountryCodeA2   string
}

type SkippedPhoneNumberValidationCommand struct {
	eventstore.BaseCommand
	RawPhoneNumber       string
	ValidationSkipReason string
	CountryCodeA2        string
}

type PhoneNumberValidatedCommand struct {
	eventstore.BaseCommand
	RawPhoneNumber string
	E164           string
	CountryCodeA2  string
}

func NewUpsertPhoneNumberCommand(objectId, tenant, loggedInUserId, rawPhoneNumber string, source cmnmod.Source, createdAt, updatedAt *time.Time) *UpsertPhoneNumberCommand {
	return &UpsertPhoneNumberCommand{
		BaseCommand:    eventstore.NewBaseCommand(objectId, tenant, loggedInUserId),
		RawPhoneNumber: rawPhoneNumber,
		Source:         source,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func NewFailedPhoneNumberValidationCommand(objectId, tenant, rawPhoneNumber, countryCodeA2, validationError string) *FailedPhoneNumberValidationCommand {
	return &FailedPhoneNumberValidationCommand{
		BaseCommand:     eventstore.NewBaseCommand(objectId, tenant, ""),
		RawPhoneNumber:  rawPhoneNumber,
		ValidationError: validationError,
		CountryCodeA2:   countryCodeA2,
	}
}

func NewSkippedPhoneNumberValidationCommand(objectId, tenant, rawPhoneNumber, countryCodeA2, validationSkipReason string) *SkippedPhoneNumberValidationCommand {
	return &SkippedPhoneNumberValidationCommand{
		BaseCommand:          eventstore.NewBaseCommand(objectId, tenant, ""),
		RawPhoneNumber:       rawPhoneNumber,
		ValidationSkipReason: validationSkipReason,
		CountryCodeA2:        countryCodeA2,
	}
}

func NewPhoneNumberValidatedCommand(objectId, tenant, rawPhoneNumber, e164, countryCodeA2 string) *PhoneNumberValidatedCommand {
	return &PhoneNumberValidatedCommand{
		BaseCommand:    eventstore.NewBaseCommand(objectId, tenant, ""),
		E164:           e164,
		RawPhoneNumber: rawPhoneNumber,
		CountryCodeA2:  countryCodeA2,
	}
}
