package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	UserCreateV1          = "V1_USER_CREATE"
	UserUpdateV1          = "V1_USER_UPDATE"
	UserPhoneNumberLinkV1 = "V1_USER_PHONE_NUMBER_LINK"
	// Deprecated
	UserEmailLinkV1 = "V1_USER_EMAIL_LINK"
	// Deprecated
	UserEmailUnlinkV1 = "V1_USER_EMAIL_UNLINK"
	UserJobRoleLinkV1 = "V1_USER_JOB_ROLE_LINK"
	UserAddRoleV1     = "V1_USER_ADD_ROLE"
	UserRemoveRoleV1  = "V1_USER_REMOVE_ROLE"
	//Deprecated
	UserAddPlayerV1 = "V1_USER_ADD_PLAYER"
)

type UserCreateEvent struct {
	Tenant          string                `json:"tenant" validate:"required"`
	Name            string                `json:"name"`
	FirstName       string                `json:"firstName"`
	LastName        string                `json:"lastName"`
	SourceFields    cmnmod.Source         `json:"sourceFields"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	Internal        bool                  `json:"internal"`
	Test            bool                  `json:"test"`
	Bot             bool                  `json:"bot"`
	ProfilePhotoUrl string                `json:"profilePhotoUrl"`
	Timezone        string                `json:"timezone"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewUserCreateEvent(aggregate eventstore.Aggregate, dataFields models.UserDataFields, sourceFields cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserCreateEvent{
		Tenant:          aggregate.GetTenant(),
		Name:            dataFields.Name,
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Internal:        dataFields.Internal,
		Bot:             dataFields.Bot,
		Test:            dataFields.Test,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		Timezone:        dataFields.Timezone,
		SourceFields:    sourceFields,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserCreateEvent")
	}
	return event, nil
}

type UserUpdateEvent struct {
	Tenant          string                `json:"tenant" validate:"required"`
	Source          string                `json:"source"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	Name            string                `json:"name"`
	FirstName       string                `json:"firstName"`
	LastName        string                `json:"lastName"`
	ProfilePhotoUrl string                `json:"profilePhotoUrl"`
	Timezone        string                `json:"timezone"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewUserUpdateEvent(aggregate eventstore.Aggregate, dataFields models.UserDataFields, source string, updatedAt time.Time, externalSystem cmnmod.ExternalSystem) (eventstore.Event, error) {
	eventData := UserUpdateEvent{
		Tenant:          aggregate.GetTenant(),
		Name:            dataFields.Name,
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		Timezone:        dataFields.Timezone,
		UpdatedAt:       updatedAt,
		Source:          source,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserUpdateEvent")
	}
	return event, nil
}

type UserLinkJobRoleEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
	JobRoleId string    `json:"jobRoleId" validate:"required"`
}

func NewUserLinkJobRoleEvent(aggregate eventstore.Aggregate, tenant, jobRoleId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserLinkJobRoleEvent{
		Tenant:    tenant,
		UpdatedAt: updatedAt,
		JobRoleId: jobRoleId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserLinkJobRoleEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserJobRoleLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserLinkJobRoleEvent")
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
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserLinkPhoneNumberEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserPhoneNumberLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserLinkPhoneNumberEvent")
	}
	return event, nil
}

type UserAddRoleEvent struct {
	Tenant string    `json:"tenant" validate:"required"`
	Role   string    `json:"role" validate:"required"`
	At     time.Time `json:"at"`
}

func NewUserAddRoleEvent(aggregate eventstore.Aggregate, role string, at time.Time) (eventstore.Event, error) {
	eventData := UserAddRoleEvent{
		Tenant: aggregate.GetTenant(),
		Role:   role,
		At:     at,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserAddRoleEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserAddRoleV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserAddRoleEvent")
	}
	return event, nil
}

type UserRemoveRoleEvent struct {
	Tenant string    `json:"tenant" validate:"required"`
	Role   string    `json:"role"`
	At     time.Time `json:"at"`
}

func NewUserRemoveRoleEvent(aggregate eventstore.Aggregate, role string, at time.Time) (eventstore.Event, error) {
	eventData := UserRemoveRoleEvent{
		Tenant: aggregate.GetTenant(),
		Role:   role,
		At:     at,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserRemoveRoleEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserRemoveRoleV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserRemoveRoleEvent")
	}
	return event, nil
}
