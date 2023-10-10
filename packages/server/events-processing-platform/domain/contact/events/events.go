package events

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
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
	ContactLocationLinkV1    = "V1_CONTACT_LOCATION_LINK"
)

type ContactCreateEvent struct {
	Tenant          string                `json:"tenant" validate:"required"`
	FirstName       string                `json:"firstName"`
	LastName        string                `json:"lastName"`
	Name            string                `json:"name"`
	Prefix          string                `json:"prefix"`
	Description     string                `json:"description"`
	Timezone        string                `json:"timezone"`
	ProfilePhotoUrl string                `json:"profilePhotoUrl"`
	Source          string                `json:"source"`
	SourceOfTruth   string                `json:"sourceOfTruth"`
	AppSource       string                `json:"appSource"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewContactCreateEvent(aggregate eventstore.Aggregate, dataFields models.ContactDataFields, sourceFields cmnmod.Source,
	externalSystem cmnmod.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactCreateEvent{
		Tenant:          aggregate.GetTenant(),
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Name:            dataFields.Name,
		Prefix:          dataFields.Prefix,
		Description:     dataFields.Description,
		Timezone:        dataFields.Timezone,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		Source:          sourceFields.Source,
		SourceOfTruth:   sourceFields.SourceOfTruth,
		AppSource:       sourceFields.AppSource,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
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
	Tenant          string                `json:"tenant" validate:"required"`
	Source          string                `json:"source"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	FirstName       string                `json:"firstName"`
	LastName        string                `json:"lastName"`
	Name            string                `json:"name"`
	Prefix          string                `json:"prefix"`
	Description     string                `json:"description"`
	Timezone        string                `json:"timezone"`
	ProfilePhotoUrl string                `json:"profilePhotoUrl"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewContactUpdateEvent(aggregate eventstore.Aggregate, source string, dataFields models.ContactDataFields, externalSystem cmnmod.ExternalSystem, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactUpdateEvent{
		Tenant:          aggregate.GetTenant(),
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Prefix:          dataFields.Prefix,
		Description:     dataFields.Description,
		Timezone:        dataFields.Timezone,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		Name:            dataFields.Name,
		UpdatedAt:       updatedAt,
		Source:          source,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
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

func NewContactLinkPhoneNumberEvent(aggregate eventstore.Aggregate, phoneNumberId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkPhoneNumberEvent{
		Tenant:        aggregate.GetTenant(),
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

func NewContactLinkEmailEvent(aggregate eventstore.Aggregate, emailId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkEmailEvent{
		Tenant:    aggregate.GetTenant(),
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

type ContactLinkLocationEvent struct {
	Tenant     string    `json:"tenant" validate:"required"`
	UpdatedAt  time.Time `json:"updatedAt"`
	LocationId string    `json:"locationId" validate:"required"`
}

func NewContactLinkLocationEvent(aggregate eventstore.Aggregate, locationId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkLocationEvent{
		Tenant:     aggregate.GetTenant(),
		UpdatedAt:  updatedAt,
		LocationId: locationId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, ContactLocationLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
