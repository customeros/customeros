package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	ContactCreatedV1           = "V1_CONTACT_CREATED"
	ContactUpdatedV1           = "V1_CONTACT_UPDATED"
	ContactPhoneNumberLinkedV1 = "V1_CONTACT_PHONE_NUMBER_LINKED"
	ContactEmailLinkedV1       = "V1_CONTACT_EMAIL_LINKED"
)

type ContactCreatedEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Prefix        string    `json:"prefix"`
	Source        string    `json:"source"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	AppSource     string    `json:"appSource"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewContactCreatedEvent(aggregate eventstore.Aggregate, contactDto *models.ContactDto, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactCreatedEvent{
		Tenant:        contactDto.Tenant,
		FirstName:     contactDto.FirstName,
		LastName:      contactDto.LastName,
		Prefix:        contactDto.Prefix,
		Source:        contactDto.Source.Source,
		SourceOfTruth: contactDto.Source.SourceOfTruth,
		AppSource:     contactDto.Source.AppSource,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, ContactCreatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactUpdatedEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Prefix        string    `json:"prefix"`
}

func NewContactUpdatedEvent(aggregate eventstore.Aggregate, contactDto *models.ContactDto, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactUpdatedEvent{
		FirstName:     contactDto.FirstName,
		LastName:      contactDto.LastName,
		Prefix:        contactDto.Prefix,
		Tenant:        contactDto.Tenant,
		UpdatedAt:     updatedAt,
		SourceOfTruth: contactDto.Source.SourceOfTruth,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, ContactUpdatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactLinkPhoneNumberEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PhoneNumberId string    `json:"phoneNumberId" validate:"required"`
	Label         string    `json:"label"`
	Primary       bool      `json:"primary"`
}

func NewContactLinkPhoneNumberEvent(aggregate eventstore.Aggregate, tenant, phoneNumberId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkPhoneNumberEvent{
		Tenant:        tenant,
		UpdatedAt:     updatedAt,
		PhoneNumberId: phoneNumberId,
		Label:         label,
		Primary:       primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, ContactPhoneNumberLinkedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactLinkEmailEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
	EmailId   string    `json:"emailId" validate:"required"`
	Label     string    `json:"label"`
	Primary   bool      `json:"primary"`
}

func NewContactLinkEmailEvent(aggregate eventstore.Aggregate, tenant, emailId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkEmailEvent{
		Tenant:    tenant,
		UpdatedAt: updatedAt,
		EmailId:   emailId,
		Label:     label,
		Primary:   primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, ContactEmailLinkedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
