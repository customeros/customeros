package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	UserCreateV1          = "V1_USER_CREATE"
	UserUpdateV1          = "V1_USER_UPDATE"
	UserPhoneNumberLinkV1 = "V1_USER_PHONE_NUMBER_LINK"
	UserEmailLinkV1       = "V1_USER_EMAIL_LINK"
	UserJobRoleLinkV1     = "V1_USER_JOB_ROLE_LINK"
)

type UserCreateEvent struct {
	Tenant          string    `json:"tenant" validate:"required"`
	Name            string    `json:"name"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Source          string    `json:"source"`
	SourceOfTruth   string    `json:"sourceOfTruth"`
	AppSource       string    `json:"appSource"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Internal        bool      `json:"internal"`
	ProfilePhotoUrl string    `json:"profilePhotoUrl"`
}

func NewUserCreateEvent(aggregate eventstore.Aggregate, userDto *models.UserDto, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserCreateEvent{
		Tenant:          userDto.Tenant,
		Name:            userDto.UserCoreFields.Name,
		FirstName:       userDto.UserCoreFields.FirstName,
		LastName:        userDto.UserCoreFields.LastName,
		Internal:        userDto.UserCoreFields.Internal,
		ProfilePhotoUrl: userDto.UserCoreFields.ProfilePhotoUrl,
		Source:          userDto.Source.Source,
		SourceOfTruth:   userDto.Source.SourceOfTruth,
		AppSource:       userDto.Source.AppSource,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
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
	Tenant          string    `json:"tenant" validate:"required"`
	SourceOfTruth   string    `json:"sourceOfTruth"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Name            string    `json:"name"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Internal        bool      `json:"internal"`
	ProfilePhotoUrl string    `json:"profilePhotoUrl"`
}

func NewUserUpdateEvent(aggregate eventstore.Aggregate, userDto *models.UserDto, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UserUpdateEvent{
		Name:            userDto.UserCoreFields.Name,
		FirstName:       userDto.UserCoreFields.FirstName,
		LastName:        userDto.UserCoreFields.LastName,
		Internal:        userDto.UserCoreFields.Internal,
		ProfilePhotoUrl: userDto.UserCoreFields.ProfilePhotoUrl,
		Tenant:          userDto.Tenant,
		UpdatedAt:       updatedAt,
		SourceOfTruth:   userDto.Source.SourceOfTruth,
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
