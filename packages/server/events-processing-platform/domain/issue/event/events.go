package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	IssueAddUserAssigneeV1    = "V1_ISSUE_ADD_USER_ASSIGNEE"
	IssueRemoveUserAssigneeV1 = "V1_ISSUE_REMOVE_USER_ASSIGNEE"
	IssueAddUserFollowerV1    = "V1_ISSUE_ADD_USER_FOLLOWER"
	IssueRemoveUserFollowerV1 = "V1_ISSUE_REMOVE_USER_FOLLOWER"
)

type IssueAddUserAssigneeEvent struct {
	Tenant string    `json:"tenant" validate:"required"`
	At     time.Time `json:"at"`
	UserId string    `json:"userId" validate:"required"`
}

func NewIssueAddUserAssigneeEvent(aggregate eventstore.Aggregate, userId string, at time.Time) (eventstore.Event, error) {
	eventData := IssueAddUserAssigneeEvent{
		Tenant: aggregate.GetTenant(),
		At:     at,
		UserId: userId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate IssueAddUserAssigneeEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, IssueAddUserAssigneeV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for IssueAddUserAssigneeEvent")
	}
	return event, nil
}

type IssueRemoveUserAssigneeEvent struct {
	Tenant string    `json:"tenant" validate:"required"`
	At     time.Time `json:"at"`
	UserId string    `json:"userId" validate:"required"`
}

func NewIssueRemoveUserAssigneeEvent(aggregate eventstore.Aggregate, userId string, at time.Time) (eventstore.Event, error) {
	eventData := IssueRemoveUserAssigneeEvent{
		Tenant: aggregate.GetTenant(),
		At:     at,
		UserId: userId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate IssueRemoveUserAssigneeEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, IssueRemoveUserAssigneeV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for IssueRemoveUserAssigneeEvent")
	}
	return event, nil
}

type IssueAddUserFollowerEvent struct {
	Tenant string    `json:"tenant" validate:"required"`
	At     time.Time `json:"at"`
	UserId string    `json:"userId" validate:"required"`
}

func NewIssueAddUserFollowerEvent(aggregate eventstore.Aggregate, userId string, at time.Time) (eventstore.Event, error) {
	eventData := IssueAddUserFollowerEvent{
		Tenant: aggregate.GetTenant(),
		At:     at,
		UserId: userId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate IssueAddUserFollowerEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, IssueAddUserFollowerV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for IssueAddUserFollowerEvent")
	}
	return event, nil
}

type IssueRemoveUserFollowerEvent struct {
	Tenant string    `json:"tenant" validate:"required"`
	At     time.Time `json:"at"`
	UserId string    `json:"userId" validate:"required"`
}

func NewIssueRemoveUserFollowerEvent(aggregate eventstore.Aggregate, userId string, at time.Time) (eventstore.Event, error) {
	eventData := IssueRemoveUserFollowerEvent{
		Tenant: aggregate.GetTenant(),
		At:     at,
		UserId: userId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate IssueRemoveUserFollowerEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, IssueRemoveUserFollowerV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for IssueRemoveUserFollowerEvent")
	}
	return event, nil
}
