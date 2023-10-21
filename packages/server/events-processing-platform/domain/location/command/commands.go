package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertLocationCommand struct {
	eventstore.BaseCommand
	IsCreateCommand       bool
	Source                cmnmod.Source
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

func NewUpsertLocationCommand(locationId, tenant, userId, name, rawAddress string, addressFields models.LocationAddressFields, source cmnmod.Source, createdAt, updatedAt *time.Time) *UpsertLocationCommand {
	return &UpsertLocationCommand{
		BaseCommand:           eventstore.NewBaseCommand(locationId, tenant, userId),
		RawAddress:            rawAddress,
		Name:                  name,
		LocationAddressFields: addressFields,
		Source:                source,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}
}

func NewFailedLocationValidationCommand(locationId, tenant, userId, rawAddress, country, validationError string) *FailedLocationValidationCommand {
	return &FailedLocationValidationCommand{
		BaseCommand:     eventstore.NewBaseCommand(locationId, tenant, userId),
		RawAddress:      rawAddress,
		Country:         country,
		ValidationError: validationError,
	}
}

func NewSkippedLocationValidationCommand(locationId, tenant, userId, rawAddress, validationSkipReason string) *SkippedLocationValidationCommand {
	return &SkippedLocationValidationCommand{
		BaseCommand:          eventstore.NewBaseCommand(locationId, tenant, userId),
		RawAddress:           rawAddress,
		ValidationSkipReason: validationSkipReason,
	}
}

func NewLocationValidatedCommand(locationId, tenant, userId, rawAddress, countryForValidation string, addressFields models.LocationAddressFields) *LocationValidatedCommand {
	return &LocationValidatedCommand{
		BaseCommand:           eventstore.NewBaseCommand(locationId, tenant, userId),
		RawAddress:            rawAddress,
		CountryForValidation:  countryForValidation,
		LocationAddressFields: addressFields,
	}
}
