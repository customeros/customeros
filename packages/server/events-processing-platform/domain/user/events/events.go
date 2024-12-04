package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	UserPhoneNumberLinkV1 = "V1_USER_PHONE_NUMBER_LINK"
	UserJobRoleLinkV1     = "V1_USER_JOB_ROLE_LINK"
	UserAddRoleV1         = "V1_USER_ADD_ROLE"
	UserRemoveRoleV1      = "V1_USER_REMOVE_ROLE"
)

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
