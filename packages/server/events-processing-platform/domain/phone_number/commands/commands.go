package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreatePhoneNumberCommand struct {
	eventstore.BaseCommand
	Tenant         string
	RawPhoneNumber string
	Source         models.Source
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

type UpsertPhoneNumberCommand struct {
	eventstore.BaseCommand
	Tenant         string
	RawPhoneNumber string
	Source         models.Source
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

type FailedPhoneNumberValidationCommand struct {
	eventstore.BaseCommand
	Tenant          string
	RawPhoneNumber  string
	ValidationError string
	CountryCodeA2   string
}

type SkippedPhoneNumberValidationCommand struct {
	eventstore.BaseCommand
	Tenant               string
	RawPhoneNumber       string
	ValidationSkipReason string
	CountryCodeA2        string
}

type PhoneNumberValidatedCommand struct {
	eventstore.BaseCommand
	Tenant         string
	RawPhoneNumber string
	E164           string
	CountryCodeA2  string
}

func NewCreatePhoneNumberCommand(baseAggregateId, tenant, rawPhoneNumber, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *CreatePhoneNumberCommand {
	return &CreatePhoneNumberCommand{
		BaseCommand:    eventstore.NewBaseCommand(baseAggregateId),
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
		Source: models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewUpsertPhoneNumberCommand(baseAggregateId, tenant, rawPhoneNumber, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *UpsertPhoneNumberCommand {
	return &UpsertPhoneNumberCommand{
		BaseCommand:    eventstore.NewBaseCommand(baseAggregateId),
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
		Source: models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewFailedPhoneNumberValidationCommand(baseAggregateId, tenant, rawPhoneNumber, countryCodeA2, validationError string) *FailedPhoneNumberValidationCommand {
	return &FailedPhoneNumberValidationCommand{
		BaseCommand:     eventstore.NewBaseCommand(baseAggregateId),
		Tenant:          tenant,
		RawPhoneNumber:  rawPhoneNumber,
		ValidationError: validationError,
		CountryCodeA2:   countryCodeA2,
	}
}

func NewSkippedPhoneNumberValidationCommand(baseAggregateId, tenant, rawPhoneNumber, countryCodeA2, validationSkipReason string) *SkippedPhoneNumberValidationCommand {
	return &SkippedPhoneNumberValidationCommand{
		BaseCommand:          eventstore.NewBaseCommand(baseAggregateId),
		Tenant:               tenant,
		RawPhoneNumber:       rawPhoneNumber,
		ValidationSkipReason: validationSkipReason,
		CountryCodeA2:        countryCodeA2,
	}
}

func NewPhoneNumberValidatedCommand(baseAggregateId, tenant, rawPhoneNumber, e164, countryCodeA2 string) *PhoneNumberValidatedCommand {
	return &PhoneNumberValidatedCommand{
		BaseCommand:    eventstore.NewBaseCommand(baseAggregateId),
		Tenant:         tenant,
		E164:           e164,
		RawPhoneNumber: rawPhoneNumber,
		CountryCodeA2:  countryCodeA2,
	}
}
