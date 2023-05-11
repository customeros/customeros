package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	PhoneNumberCreatedV1           = "V1_PHONE_NUMBER_CREATED"
	PhoneNumberUpdatedV1           = "V1_PHONE_NUMBER_UPDATED"
	PhoneNumberValidationFailedV1  = "V1_PHONE_NUMBER_VALIDATION_FAILED"
	PhoneNumberValidationSkippedV1 = "V1_PHONE_NUMBER_VALIDATION_SKIPPED"
	PhoneNumberValidatedV1         = "V1_PHONE_NUMBER_VALIDATED"
)

type PhoneNumberCreatedEvent struct {
	Tenant         string    `json:"tenant" validate:"required"`
	RawPhoneNumber string    `json:"rawPhoneNumber" validate:"required"`
	Source         string    `json:"source"`
	SourceOfTruth  string    `json:"sourceOfTruth"`
	AppSource      string    `json:"appSource"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func NewPhoneNumberCreatedEvent(aggregate eventstore.Aggregate, tenant, rawPhoneNumber, source, sourceOfTruth, appSource string, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := PhoneNumberCreatedEvent{
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
		Source:         source,
		SourceOfTruth:  sourceOfTruth,
		AppSource:      appSource,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberCreatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type PhoneNumberUpdatedEvent struct {
	Tenant        string    `json:"tenant"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewPhoneNumberUpdatedEvent(aggregate eventstore.Aggregate, tenant, sourceOfTruth string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := PhoneNumberUpdatedEvent{
		Tenant:        tenant,
		SourceOfTruth: sourceOfTruth,
		UpdatedAt:     updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberUpdatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type PhoneNumberFailedValidationEvent struct {
	Tenant          string    `json:"tenant" validate:"required"`
	ValidationError string    `json:"validationError" validate:"required"`
	ValidatedAt     time.Time `json:"validatedAt" validate:"required"`
}

func NewPhoneNumberFailedValidationEvent(aggregate eventstore.Aggregate, tenant, validationError string) (eventstore.Event, error) {
	eventData := PhoneNumberFailedValidationEvent{
		Tenant:          tenant,
		ValidationError: validationError,
		ValidatedAt:     utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberValidationFailedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type PhoneNumberSkippedValidationEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	Reason string `json:"validationSkipReason" validate:"required"`
}

func NewPhoneNumberSkippedValidationEvent(aggregate eventstore.Aggregate, tenant, validationSkipReason string) (eventstore.Event, error) {
	eventData := PhoneNumberSkippedValidationEvent{
		Tenant: tenant,
		Reason: validationSkipReason,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberValidationSkippedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type PhoneNumberValidatedEvent struct {
	Tenant      string    `json:"tenant" validate:"required"`
	PhoneNumber string    `json:"phoneNumber" validate:"required"`
	E164        string    `json:"e164" validate:"required,e164"`
	ValidatedAt time.Time `json:"validatedAt" validate:"required"`
}

func NewPhoneNumberValidatedEvent(aggregate eventstore.Aggregate, tenant, phoneNumber, e164 string) (eventstore.Event, error) {
	eventData := PhoneNumberValidatedEvent{
		Tenant:      tenant,
		PhoneNumber: phoneNumber,
		E164:        e164,
		ValidatedAt: utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, PhoneNumberValidatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
