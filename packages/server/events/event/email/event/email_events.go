package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	EmailUpsertV1    = "V1_EMAIL_UPSERT"
	EmailValidateV1  = "V1_EMAIL_VALIDATE"
	EmailValidatedV2 = "V2_EMAIL_VALIDATED"
	EmailDeleteV1    = "V1_EMAIL_DELETE"
	// Deprecated
	EmailCreateV1 = "V1_EMAIL_CREATE"
	// Deprecated
	EmailUpdateV1 = "V1_EMAIL_UPDATE"
	//Deprecated
	EmailValidatedV1 = "V1_EMAIL_VALIDATED"
	//Deprecated
	EmailValidationFailedV1 = "V1_EMAIL_VALIDATION_FAILED"
)

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
