package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	ContactCreateV1          = "V1_CONTACT_CREATE"
	ContactUpdateV1          = "V1_CONTACT_UPDATE"
	ContactPhoneNumberLinkV1 = "V1_CONTACT_PHONE_NUMBER_LINK"
	ContactEmailLinkV1       = "V1_CONTACT_EMAIL_LINK"
)

type ContactCreateEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Name          string    `json:"name"`
	Prefix        string    `json:"prefix"`
	Description   string    `json:"description"`
	Timezone      string    `json:"timezone"`
	Source        string    `json:"source"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	AppSource     string    `json:"appSource"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewContactCreateEvent(aggregate eventstore.Aggregate, contactDto *models.ContactDto, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactCreateEvent{
		Tenant:        contactDto.Tenant,
		FirstName:     contactDto.FirstName,
		LastName:      contactDto.LastName,
		Name:          contactDto.Name,
		Prefix:        contactDto.Prefix,
		Description:   contactDto.Description,
		Timezone:      contactDto.Timezone,
		Source:        contactDto.Source.Source,
		SourceOfTruth: contactDto.Source.SourceOfTruth,
		AppSource:     contactDto.Source.AppSource,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, ContactCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactUpdateEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Name          string    `json:"name"`
	Prefix        string    `json:"prefix"`
	Description   string    `json:"description"`
	Timezone      string    `json:"timezone"`
}

func NewContactUpdateEvent(aggregate eventstore.Aggregate, contactDto *models.ContactDto, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactUpdateEvent{
		FirstName:     contactDto.FirstName,
		LastName:      contactDto.LastName,
		Prefix:        contactDto.Prefix,
		Description:   contactDto.Description,
		Timezone:      contactDto.Timezone,
		Name:          contactDto.Name,
		Tenant:        contactDto.Tenant,
		UpdatedAt:     updatedAt,
		SourceOfTruth: contactDto.Source.SourceOfTruth,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, ContactUpdateV1)
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

	event := eventstore.NewBaseEvent(aggregate, ContactPhoneNumberLinkV1)
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

	event := eventstore.NewBaseEvent(aggregate, ContactEmailLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
