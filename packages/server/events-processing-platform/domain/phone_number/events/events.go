package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

const (
	PhoneNumberCreateV1            = "V1_PHONE_NUMBER_CREATE"
	PhoneNumberUpdateV1            = "V1_PHONE_NUMBER_UPDATE"
	PhoneNumberValidationFailedV1  = "V1_PHONE_NUMBER_VALIDATION_FAILED"
	PhoneNumberValidationSkippedV1 = "V1_PHONE_NUMBER_VALIDATION_SKIPPED"
	PhoneNumberValidatedV1         = "V1_PHONE_NUMBER_VALIDATED"
)

type PhoneNumberCreateEvent struct {
	Tenant         string        `json:"tenant" validate:"required"`
	RawPhoneNumber string        `json:"rawPhoneNumber" validate:"required"`
	Source         string        `json:"source"`        //Deprecated
	SourceOfTruth  string        `json:"sourceOfTruth"` //Deprecated
	AppSource      string        `json:"appSource"`     //Deprecated
	SourceFields   cmnmod.Source `json:"sourceFields"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}

func NewPhoneNumberCreateEvent(aggregate eventstore.Aggregate, tenant, rawPhoneNumber string, source cmnmod.Source, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := PhoneNumberCreateEvent{
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
		SourceFields:   source,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate PhoneNumberCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for PhoneNumberCreateEvent")
	}
	return event, nil
}

type PhoneNumberUpdatedEvent struct {
	Tenant    string    `json:"tenant"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewPhoneNumberUpdateEvent(aggregate eventstore.Aggregate, tenant, source string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := PhoneNumberUpdatedEvent{
		Tenant:    tenant,
		Source:    source,
		UpdatedAt: updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate PhoneNumberUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for PhoneNumberUpdateEvent")
	}
	return event, nil
}

type PhoneNumberFailedValidationEvent struct {
	Tenant          string    `json:"tenant" validate:"required"`
	ValidationError string    `json:"validationError" validate:"required"`
	RawPhoneNumber  string    `json:"rawPhoneNumber" validate:"required"`
	CountryCodeA2   string    `json:"countryCodeA2UsedForValidation"`
	ValidatedAt     time.Time `json:"validatedAt" validate:"required"`
}

func NewPhoneNumberFailedValidationEvent(aggregate eventstore.Aggregate, tenant, rawPhoneNumber, countryCodeA2, validationError string) (eventstore.Event, error) {
	eventData := PhoneNumberFailedValidationEvent{
		Tenant:          tenant,
		ValidationError: validationError,
		RawPhoneNumber:  rawPhoneNumber,
		CountryCodeA2:   countryCodeA2,
		ValidatedAt:     utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate PhoneNumberFailedValidationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberValidationFailedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for PhoneNumberFailedValidationEvent")
	}
	return event, nil
}

type PhoneNumberSkippedValidationEvent struct {
	Tenant         string `json:"tenant" validate:"required"`
	RawPhoneNumber string `json:"rawPhoneNumber" validate:"required"`
	CountryCodeA2  string `json:"countryCodeA2UsedForValidation"`
	Reason         string `json:"validationSkipReason" validate:"required"`
}

func NewPhoneNumberSkippedValidationEvent(aggregate eventstore.Aggregate, tenant, rawPhoneNumber, countryCodeA2, validationSkipReason string) (eventstore.Event, error) {
	eventData := PhoneNumberSkippedValidationEvent{
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
		CountryCodeA2:  countryCodeA2,
		Reason:         validationSkipReason,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate PhoneNumberSkippedValidationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberValidationSkippedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for PhoneNumberSkippedValidationEvent")
	}
	return event, nil
}

type PhoneNumberValidatedEvent struct {
	Tenant         string    `json:"tenant" validate:"required"`
	RawPhoneNumber string    `json:"rawPhoneNumber" validate:"required"`
	E164           string    `json:"e164" validate:"required,e164"`
	CountryCodeA2  string    `json:"countryCodeA2"`
	ValidatedAt    time.Time `json:"validatedAt" validate:"required"`
}

func NewPhoneNumberValidatedEvent(aggregate eventstore.Aggregate, tenant, rawPhoneNumber, e164, countryCodeA2 string) (eventstore.Event, error) {
	eventData := PhoneNumberValidatedEvent{
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
		E164:           e164,
		CountryCodeA2:  countryCodeA2,
		ValidatedAt:    utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate PhoneNumberValidatedEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberValidatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for PhoneNumberValidatedEvent")
	}
	return event, nil
}
