package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	UserCreateV1          = "V1_USER_CREATE"
	UserUpdateV1          = "V1_USER_UPDATE"
	UserAddPlayerV1       = "V1_USER_ADD_PLAYER"
	UserPhoneNumberLinkV1 = "V1_USER_PHONE_NUMBER_LINK"
	UserEmailLinkV1       = "V1_USER_EMAIL_LINK"
	UserJobRoleLinkV1     = "V1_USER_JOB_ROLE_LINK"
	UserAddRoleV1         = "V1_USER_ADD_ROLE"
	UserRemoveRoleV1      = "V1_USER_REMOVE_ROLE"
)

type UserCreateEvent struct {
	Tenant          string                `json:"tenant" validate:"required"`
	Name            string                `json:"name"`
	FirstName       string                `json:"firstName"`
	LastName        string                `json:"lastName"`
	SourceFields    events.Source         `json:"sourceFields"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	Internal        bool                  `json:"internal"`
	Bot             bool                  `json:"bot"`
	ProfilePhotoUrl string                `json:"profilePhotoUrl"`
	Timezone        string                `json:"timezone"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewUserCreateEvent(aggregate eventstore.Aggregate, dataFields models.UserDataFields, sourceFields events.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserCreateEvent{
		Tenant:          aggregate.GetTenant(),
		Name:            dataFields.Name,
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Internal:        dataFields.Internal,
		Bot:             dataFields.Bot,
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
	Internal        bool                  `json:"internal"`
	Bot             bool                  `json:"bot"`
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
		Internal:        dataFields.Internal,
		Bot:             dataFields.Bot,
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
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserLinkEmailEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserEmailLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserLinkEmailEvent")
	}
	return event, nil
}

type UserAddPlayerInfoEvent struct {
	Tenant       string        `json:"tenant" validate:"required"`
	Provider     string        `json:"provider" validate:"required"`
	AuthId       string        `json:"authId" validate:"required"`
	IdentityId   string        `json:"identityId"`
	CreatedAt    time.Time     `json:"createdAt"`
	SourceFields events.Source `json:"sourceFields"`
}

func NewUserAddPlayerInfoEvent(aggregate eventstore.Aggregate, dataFields models.PlayerInfo, source events.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := UserAddPlayerInfoEvent{
		Tenant:       aggregate.GetTenant(),
		Provider:     dataFields.Provider,
		AuthId:       dataFields.AuthId,
		IdentityId:   dataFields.IdentityId,
		CreatedAt:    createdAt,
		SourceFields: source,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserAddPlayerInfoEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserAddPlayerV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserAddPlayerInfoEvent")
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
