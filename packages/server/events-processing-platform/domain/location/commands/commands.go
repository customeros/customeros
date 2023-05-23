package commands

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertLocationCommand struct {
	eventstore.BaseCommand
	Source                common_models.Source
	CreatedAt             *time.Time
	UpdatedAt             *time.Time
	Name                  string
	RawAddress            string
	LocationAddressFields models.LocationAddressFields
}

type FailedLocationValidationCommand struct {
	eventstore.BaseCommand
	RawAddress      string
	Country         string
	ValidationError string
}

type SkippedLocationValidationCommand struct {
	eventstore.BaseCommand
	RawAddress           string
	ValidationSkipReason string
}

type LocationValidatedCommand struct {
	eventstore.BaseCommand
	RawAddress            string
	CountryForValidation  string
	LocationAddressFields models.LocationAddressFields
}

func NewUpsertLocationCommand(objectId, tenant, name, rawAddress string, addressFields models.LocationAddressFields, source common_models.Source, createdAt, updatedAt *time.Time) *UpsertLocationCommand {
	return &UpsertLocationCommand{
		BaseCommand:           eventstore.NewBaseCommand(objectId, tenant),
		RawAddress:            rawAddress,
		Name:                  name,
		LocationAddressFields: addressFields,
		Source:                source,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}
}

func NewFailedLocationValidationCommand(objectId, tenant, rawAddress, country, validationError string) *FailedLocationValidationCommand {
	return &FailedLocationValidationCommand{
		BaseCommand:     eventstore.NewBaseCommand(objectId, tenant),
		RawAddress:      rawAddress,
		Country:         country,
		ValidationError: validationError,
	}
}

func NewSkippedLocationValidationCommand(objectId, tenant, rawAddress, validationSkipReason string) *SkippedLocationValidationCommand {
	return &SkippedLocationValidationCommand{
		BaseCommand:          eventstore.NewBaseCommand(objectId, tenant),
		RawAddress:           rawAddress,
		ValidationSkipReason: validationSkipReason,
	}
}

func NewLocationValidatedCommand(objectId, tenant, rawAddress, countryForValidation string, addressFields models.LocationAddressFields) *LocationValidatedCommand {
	return &LocationValidatedCommand{
		BaseCommand:           eventstore.NewBaseCommand(objectId, tenant),
		RawAddress:            rawAddress,
		CountryForValidation:  countryForValidation,
		LocationAddressFields: addressFields,
	}
}
