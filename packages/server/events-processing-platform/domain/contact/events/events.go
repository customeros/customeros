package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

const (
	ContactCreated           = "CONTACT_CREATED"
	ContactUpdated           = "CONTACT_UPDATED"
	ContactPhoneNumberLinked = "CONTACT_PHONE_NUMBER_LINKED"
)

type ContactCreatedEvent struct {
	Tenant        string    `json:"tenant"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Name          string    `json:"name"`
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
		Name:          contactDto.Name,
		Prefix:        contactDto.Prefix,
		Source:        contactDto.Source.Source,
		SourceOfTruth: contactDto.Source.SourceOfTruth,
		AppSource:     contactDto.Source.AppSource,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
	event := eventstore.NewBaseEvent(aggregate, ContactCreated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactUpdatedEvent struct {
	Tenant        string    `json:"tenant"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Name          string    `json:"name"`
	Prefix        string    `json:"prefix"`
}

func NewContactUpdatedEvent(aggregate eventstore.Aggregate, contactDto *models.ContactDto, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactUpdatedEvent{
		FirstName:     contactDto.FirstName,
		LastName:      contactDto.LastName,
		Name:          contactDto.Name,
		Prefix:        contactDto.Prefix,
		Tenant:        contactDto.Tenant,
		UpdatedAt:     updatedAt,
		SourceOfTruth: contactDto.Source.SourceOfTruth,
	}
	event := eventstore.NewBaseEvent(aggregate, ContactUpdated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactLinkPhoneNumberEvent struct {
	Tenant        string    `json:"tenant"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PhoneNumberId string    `json:"phoneNumberId"`
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
	event := eventstore.NewBaseEvent(aggregate, ContactPhoneNumberLinked)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

// FIXME alexb implement
//type ContactDeletedEvent struct {
//	Uuid string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
//}
//
//func NewContactDeletedEvent(aggregate eventstore.Aggregate, uuid string) (eventstore.Event, error) {
//	eventData := ContactDeletedEvent{
//		Uuid: uuid,
//	}
//	event := eventstore.NewBaseEvent(aggregate, ContactDeleted)
//	err := event.SetJsonData(&eventData)
//	if err != nil {
//		return eventstore.Event{}, err
//	}
//	return event, nil
//}
