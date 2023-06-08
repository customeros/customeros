package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	LocationCreateV1            = "V1_LOCATION_CREATE"
	LocationUpdateV1            = "V1_LOCATION_UPDATE"
	LocationValidationFailedV1  = "V1_LOCATION_VALIDATION_FAILED"
	LocationValidationSkippedV1 = "V1_LOCATION_VALIDATION_SKIPPED"
	LocationValidatedV1         = "V1_LOCATION_VALIDATED"
)

type LocationCreateEvent struct {
	Tenant          string                 `json:"tenant" validate:"required"`
	Source          string                 `json:"source"`
	SourceOfTruth   string                 `json:"sourceOfTruth"`
	AppSource       string                 `json:"appSource"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	Name            string                 `json:"name"`
	RawAddress      string                 `json:"rawAddress"`
	LocationAddress models.LocationAddress `json:"address"`
}

func NewLocationCreateEvent(aggregate eventstore.Aggregate, tenant, name, rawAddress, source, sourceOfTruth, appSource string, createdAt, updatedAt time.Time, locationAddress models.LocationAddress) (eventstore.Event, error) {
	eventData := LocationCreateEvent{
		Tenant:          tenant,
		Source:          source,
		SourceOfTruth:   sourceOfTruth,
		AppSource:       appSource,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Name:            name,
		RawAddress:      rawAddress,
		LocationAddress: locationAddress,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, LocationCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type LocationUpdateEvent struct {
	Tenant          string                 `json:"tenant"`
	SourceOfTruth   string                 `json:"sourceOfTruth"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	Name            string                 `json:"name"`
	RawAddress      string                 `json:"rawAddress"`
	LocationAddress models.LocationAddress `json:"address"`
}

func NewLocationUpdateEvent(aggregate eventstore.Aggregate, tenant, name, rawAddress, sourceOfTruth string, updatedAt time.Time, locationAddress models.LocationAddress) (eventstore.Event, error) {
	eventData := LocationUpdateEvent{
		Tenant:          tenant,
		SourceOfTruth:   sourceOfTruth,
		UpdatedAt:       updatedAt,
		Name:            name,
		RawAddress:      rawAddress,
		LocationAddress: locationAddress,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, LocationUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type LocationFailedValidationEvent struct {
	Tenant          string    `json:"tenant" validate:"required"`
	ValidationError string    `json:"validationError" validate:"required"`
	RawAddress      string    `json:"rawAddress" validate:"required"`
	Country         string    `json:"country" `
	ValidatedAt     time.Time `json:"validatedAt" validate:"required"`
}

func NewLocationFailedValidationEvent(aggregate eventstore.Aggregate, tenant, rawAddress, country, validationError string) (eventstore.Event, error) {
	eventData := LocationFailedValidationEvent{
		Tenant:          tenant,
		ValidationError: validationError,
		RawAddress:      rawAddress,
		Country:         country,
		ValidatedAt:     utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, LocationValidationFailedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type LocationSkippedValidationEvent struct {
	Tenant      string    `json:"tenant" validate:"required"`
	RawAddress  string    `json:"rawAddress"`
	Reason      string    `json:"validationSkipReason" validate:"required"`
	ValidatedAt time.Time `json:"validatedAt" validate:"required"`
}

func NewLocationSkippedValidationEvent(aggregate eventstore.Aggregate, tenant, rawAddress, validationSkipReason string) (eventstore.Event, error) {
	eventData := LocationSkippedValidationEvent{
		Tenant:      tenant,
		RawAddress:  rawAddress,
		Reason:      validationSkipReason,
		ValidatedAt: utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, LocationValidationSkippedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type LocationValidatedEvent struct {
	Tenant               string                 `json:"tenant" validate:"required"`
	RawAddress           string                 `json:"rawAddress" validate:"required"`
	CountryForValidation string                 `json:"countryForValidation" `
	ValidatedAt          time.Time              `json:"validatedAt" validate:"required"`
	LocationAddress      models.LocationAddress `json:"address"`
}

func NewLocationValidatedEvent(aggregate eventstore.Aggregate, tenant, rawAddress, countryForValidation string, locationAddress models.LocationAddress) (eventstore.Event, error) {
	eventData := LocationValidatedEvent{
		Tenant:               tenant,
		RawAddress:           rawAddress,
		CountryForValidation: countryForValidation,
		LocationAddress:      locationAddress,
		ValidatedAt:          utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, LocationValidatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
