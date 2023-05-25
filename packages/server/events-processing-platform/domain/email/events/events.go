package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	EmailCreatedV1          = "V1_EMAIL_CREATED"
	EmailUpdatedV1          = "V1_EMAIL_UPDATED"
	EmailValidationFailedV1 = "V1_EMAIL_VALIDATION_FAILED"
	EmailValidatedV1        = "V1_EMAIL_VALIDATED"
)

// TODO handle case when any event arrives before EmailCreatedV1 event

type EmailCreatedEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	RawEmail      string    `json:"rawEmail" validate:"required"`
	Source        string    `json:"source"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	AppSource     string    `json:"appSource"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewEmailCreatedEvent(aggregate eventstore.Aggregate, tenant, rawEmail, source, sourceOfTruth, appSource string, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := EmailCreatedEvent{
		Tenant:        tenant,
		RawEmail:      rawEmail,
		Source:        source,
		SourceOfTruth: sourceOfTruth,
		AppSource:     appSource,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, EmailCreatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type EmailUpdatedEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewEmailUpdatedEvent(aggregate eventstore.Aggregate, tenant, sourceOfTruth string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := EmailUpdatedEvent{
		Tenant:        tenant,
		SourceOfTruth: sourceOfTruth,
		UpdatedAt:     updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, EmailUpdatedV1)
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
