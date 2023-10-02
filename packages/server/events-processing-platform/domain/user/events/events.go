package events

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
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
	Tenant          string                       `json:"tenant" validate:"required"`
	Name            string                       `json:"name"`
	FirstName       string                       `json:"firstName"`
	LastName        string                       `json:"lastName"`
	SourceFields    common_models.Source         `json:"sourceFields"`
	CreatedAt       time.Time                    `json:"createdAt"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
	Internal        bool                         `json:"internal"`
	ProfilePhotoUrl string                       `json:"profilePhotoUrl"`
	Timezone        string                       `json:"timezone"`
	ExternalSystem  common_models.ExternalSystem `json:"externalSystem"`
}

func NewUserCreateEvent(aggregate eventstore.Aggregate, dataFields models.UserDataFields, source common_models.Source, externalSystem common_models.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserCreateEvent{
		Tenant:          aggregate.GetTenant(),
		Name:            dataFields.Name,
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Internal:        dataFields.Internal,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		Timezone:        dataFields.Timezone,
		SourceFields:    source,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type UserUpdateEvent struct {
	Tenant          string                       `json:"tenant" validate:"required"`
	Source          string                       `json:"source"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
	Name            string                       `json:"name"`
	FirstName       string                       `json:"firstName"`
	LastName        string                       `json:"lastName"`
	Internal        bool                         `json:"internal"`
	ProfilePhotoUrl string                       `json:"profilePhotoUrl"`
	Timezone        string                       `json:"timezone"`
	ExternalSystem  common_models.ExternalSystem `json:"externalSystem"`
}

func NewUserUpdateEvent(aggregate eventstore.Aggregate, dataFields models.UserDataFields, source string, updatedAt time.Time, externalSystem common_models.ExternalSystem) (eventstore.Event, error) {
	eventData := UserUpdateEvent{
		Tenant:          aggregate.GetTenant(),
		Name:            dataFields.Name,
		FirstName:       dataFields.FirstName,
		LastName:        dataFields.LastName,
		Internal:        dataFields.Internal,
		ProfilePhotoUrl: dataFields.ProfilePhotoUrl,
		Timezone:        dataFields.Timezone,
		UpdatedAt:       updatedAt,
		Source:          source,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
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
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserJobRoleLinkV1)
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

	event := eventstore.NewBaseEvent(aggregate, UserPhoneNumberLinkV1)
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

	event := eventstore.NewBaseEvent(aggregate, UserEmailLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type UserAddPlayerInfoEvent struct {
	Tenant       string               `json:"tenant" validate:"required"`
	Provider     string               `json:"provider" validate:"required"`
	AuthId       string               `json:"authId" validate:"required"`
	IdentityId   string               `json:"identityId"`
	CreatedAt    time.Time            `json:"createdAt"`
	SourceFields common_models.Source `json:"sourceFields"`
}

func NewUserAddPlayerInfoEvent(aggregate eventstore.Aggregate, dataFields models.PlayerInfo, source common_models.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := UserAddPlayerInfoEvent{
		Tenant:       aggregate.GetTenant(),
		Provider:     dataFields.Provider,
		AuthId:       dataFields.AuthId,
		IdentityId:   dataFields.IdentityId,
		CreatedAt:    createdAt,
		SourceFields: source,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserAddPlayerV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
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
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserAddRoleV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
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
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, UserRemoveRoleV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
