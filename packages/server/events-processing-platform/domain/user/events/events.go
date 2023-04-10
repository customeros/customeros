package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	UserCreatedV1           = "V1_USER_CREATED"
	UserUpdatedV1           = "V1_USER_UPDATED"
	UserPhoneNumberLinkedV1 = "V1_USER_PHONE_NUMBER_LINKED"
	UserEmailLinkedV1       = "V1_USER_EMAIL_LINKED"
)

type UserCreatedEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	Name          string    `json:"name"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Source        string    `json:"source"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	AppSource     string    `json:"appSource"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewUserCreatedEvent(aggregate eventstore.Aggregate, userDto *models.UserDto, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserCreatedEvent{
		Tenant:        userDto.Tenant,
		Name:          userDto.Name,
		FirstName:     userDto.FirstName,
		LastName:      userDto.LastName,
		Source:        userDto.Source.Source,
		SourceOfTruth: userDto.Source.SourceOfTruth,
		AppSource:     userDto.Source.AppSource,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserCreatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type UserUpdatedEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Name          string    `json:"name"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
}

func NewUserUpdatedEvent(aggregate eventstore.Aggregate, userDto *models.UserDto, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserUpdatedEvent{
		Name:          userDto.Name,
		FirstName:     userDto.FirstName,
		LastName:      userDto.LastName,
		Tenant:        userDto.Tenant,
		UpdatedAt:     updatedAt,
		SourceOfTruth: userDto.Source.SourceOfTruth,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserUpdatedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type UserLinkPhoneNumberEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PhoneNumberId string    `json:"phoneNumberId" validate:"required"`
	Label         string    `json:"label"`
	Primary       bool      `json:"primary"`
}

func NewUserLinkPhoneNumberEvent(aggregate eventstore.Aggregate, tenant, phoneNumberId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserLinkPhoneNumberEvent{
		Tenant:        tenant,
		UpdatedAt:     updatedAt,
		PhoneNumberId: phoneNumberId,
		Label:         label,
		Primary:       primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserPhoneNumberLinkedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type UserLinkEmailEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
	EmailId   string    `json:"emailId" validate:"required"`
	Label     string    `json:"label"`
	Primary   bool      `json:"primary"`
}

func NewUserLinkEmailEvent(aggregate eventstore.Aggregate, tenant, emailId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserLinkEmailEvent{
		Tenant:    tenant,
		UpdatedAt: updatedAt,
		EmailId:   emailId,
		Label:     label,
		Primary:   primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserEmailLinkedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
