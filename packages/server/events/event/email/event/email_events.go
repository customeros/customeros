package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	EmailUpsertV1    = "V1_EMAIL_UPSERT"
	EmailValidateV1  = "V1_EMAIL_VALIDATE"
	EmailValidatedV2 = "V2_EMAIL_VALIDATED"
	// Deprecated
	EmailCreateV1 = "V1_EMAIL_CREATE"
	// Deprecated
	EmailUpdateV1 = "V1_EMAIL_UPDATE"
	//Deprecated
	EmailValidatedV1 = "V1_EMAIL_VALIDATED"
	//Deprecated
	EmailValidationFailedV1 = "V1_EMAIL_VALIDATION_FAILED"
)

type EmailCreateEvent struct {
	Tenant        string        `json:"tenant" validate:"required"`
	RawEmail      string        `json:"rawEmail"`
	Source        string        `json:"source"`        //Deprecated
	SourceOfTruth string        `json:"sourceOfTruth"` //Deprecated
	AppSource     string        `json:"appSource"`     //Deprecated
	SourceFields  common.Source `json:"sourceFields"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`

	LinkWithType *string `json:"linkWith"`
	LinkWithId   *string `json:"linkWithId"`
}

func NewEmailCreateEvent(aggregate eventstore.Aggregate, tenant, rawEmail string, source common.Source, createdAt, updatedAt time.Time, linkWithType, linkWithId *string) (eventstore.Event, error) {
	eventData := EmailCreateEvent{
		Tenant:       tenant,
		RawEmail:     rawEmail,
		SourceFields: source,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		LinkWithType: linkWithType,
		LinkWithId:   linkWithId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate EmailCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, EmailCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for EmailCreateEvent")
	}
	return event, nil
}

type EmailUpdateEvent struct {
	RawEmail  string    `json:"rawEmail"`
	Tenant    string    `json:"tenant" validate:"required"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewEmailUpdateEvent(aggregate eventstore.Aggregate, tenant, rawEmail, source string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := EmailUpdateEvent{
		RawEmail:  rawEmail,
		Tenant:    tenant,
		Source:    source,
		UpdatedAt: updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate EmailUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, EmailUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for EmailUpdateEvent")
	}
	return event, nil
}

type EmailFailedValidationEvent struct {
	Tenant          string    `json:"tenant" validate:"required"`
	ValidationError string    `json:"validationError" validate:"required"`
	ValidatedAt     time.Time `json:"validatedAt" validate:"required"`
}

type EmailRequestValidationEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewEmailRequestValidationEvent(aggr eventstore.Aggregate) (eventstore.Event, error) {
	eventData := EmailRequestValidationEvent{
		Tenant: aggr.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate EmailRequestValidationEvent")
	}

	event := eventstore.NewBaseEvent(aggr, EmailValidateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for EmailRequestValidationEvent")
	}
	return event, nil
}
