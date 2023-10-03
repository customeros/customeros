package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	EmailCreateV1           = "V1_EMAIL_CREATE"
	EmailUpdateV1           = "V1_EMAIL_UPDATE"
	EmailValidationFailedV1 = "V1_EMAIL_VALIDATION_FAILED"
	EmailValidatedV1        = "V1_EMAIL_VALIDATED"
)

type EmailCreateEvent struct {
	Tenant        string               `json:"tenant" validate:"required"`
	RawEmail      string               `json:"rawEmail" validate:"required"`
	Source        string               `json:"source"`        //Deprecated
	SourceOfTruth string               `json:"sourceOfTruth"` //Deprecated
	AppSource     string               `json:"appSource"`     //Deprecated
	SourceFields  common_models.Source `json:"sourceFields"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}

func NewEmailCreateEvent(aggregate eventstore.Aggregate, tenant, rawEmail string, source common_models.Source, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := EmailCreateEvent{
		Tenant:       tenant,
		RawEmail:     rawEmail,
		SourceFields: source,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, EmailCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type EmailUpdateEvent struct {
	RawEmail  string    `json:"rawEmail" validate:"required"`
	Tenant    string    `json:"tenant" validate:"required"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewEmailUpdateEvent(aggregate eventstore.Aggregate, rawEmail, tenant, source string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := EmailUpdateEvent{
		RawEmail:  rawEmail,
		Tenant:    tenant,
		Source:    source,
		UpdatedAt: updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, EmailUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type EmailFailedValidationEvent struct {
	Tenant          string    `json:"tenant" validate:"required"`
	ValidationError string    `json:"validationError" validate:"required"`
	ValidatedAt     time.Time `json:"validatedAt" validate:"required"`
}

func NewEmailFailedValidationEvent(aggregate eventstore.Aggregate, tenant, validationError string) (eventstore.Event, error) {
	eventData := EmailFailedValidationEvent{
		Tenant:          tenant,
		ValidationError: validationError,
		ValidatedAt:     utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, EmailValidationFailedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type EmailValidatedEvent struct {
	Tenant          string    `json:"tenant" validate:"required"`
	RawEmail        string    `json:"rawEmail" validate:"required"`
	IsReachable     string    `json:"isReachable" validate:"required"`
	ValidatedAt     time.Time `json:"validatedAt" validate:"required"`
	ValidationError string    `json:"validationError"`
	AcceptsMail     bool      `json:"acceptsMail"`
	CanConnectSmtp  bool      `json:"canConnectSmtp"`
	HasFullInbox    bool      `json:"hasFullInbox"`
	IsCatchAll      bool      `json:"isCatchAll"`
	IsDeliverable   bool      `json:"isDeliverable"`
	IsDisabled      bool      `json:"isDisabled"`
	Domain          string    `json:"domain"`
	IsValidSyntax   bool      `json:"isValidSyntax"`
	Username        string    `json:"username"`
	EmailAddress    string    `json:"email"`
}

func NewEmailValidatedEvent(aggregate eventstore.Aggregate, tenant, rawEmail, isReachable, validationError, domain, username, emailAddress string,
	acceptsMail, canConnectSmtp, hasFullInbox, isCatchAll, IsDeliverable, isDisabled, isValidSyntax bool) (eventstore.Event, error) {
	eventData := EmailValidatedEvent{
		Tenant:          tenant,
		RawEmail:        rawEmail,
		IsReachable:     isReachable,
		ValidationError: validationError,
		AcceptsMail:     acceptsMail,
		CanConnectSmtp:  canConnectSmtp,
		HasFullInbox:    hasFullInbox,
		IsCatchAll:      isCatchAll,
		IsDeliverable:   IsDeliverable,
		IsDisabled:      isDisabled,
		Domain:          domain,
		IsValidSyntax:   isValidSyntax,
		Username:        username,
		EmailAddress:    emailAddress,
		ValidatedAt:     utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, EmailValidatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
