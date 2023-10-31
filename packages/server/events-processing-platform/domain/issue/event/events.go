package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

const (
	IssueCreateV1             = "V1_ISSUE_CREATE"
	IssueUpdateV1             = "V1_ISSUE_UPDATE"
	IssueAddUserAssigneeV1    = "V1_ISSUE_ADD_USER_ASSIGNEE"
	IssueRemoveUserAssigneeV1 = "V1_ISSUE_REMOVE_USER_ASSIGNEE"
	IssueAddUserFollowerV1    = "V1_ISSUE_ADD_USER_FOLLOWER"
	IssueRemoveUserFollowerV1 = "V1_ISSUE_REMOVE_USER_FOLLOWER"
)

type IssueCreateEvent struct {
	Tenant                    string                `json:"tenant" validate:"required"`
	Subject                   string                `json:"subject" validate:"required_without=Description"`
	Description               string                `json:"description" validate:"required_without=Subject"`
	Status                    string                `json:"status"`
	Priority                  string                `json:"priority"`
	ReportedByOrganizationId  string                `json:"reportedByOrganizationId,omitempty"`
	SubmittedByOrganizationId string                `json:"submittedByOrganizationId,omitempty"`
	SubmittedByUserId         string                `json:"submittedByUserId,omitempty"`
	Source                    string                `json:"source"`
	AppSource                 string                `json:"appSource"`
	CreatedAt                 time.Time             `json:"createdAt"`
	UpdatedAt                 time.Time             `json:"updatedAt"`
	ExternalSystem            cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewIssueCreateEvent(aggregate eventstore.Aggregate, dataFields model.IssueDataFields, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := IssueCreateEvent{
		Tenant:                    aggregate.GetTenant(),
		Subject:                   dataFields.Subject,
		Description:               dataFields.Description,
		Status:                    dataFields.Status,
		Priority:                  dataFields.Priority,
		ReportedByOrganizationId:  utils.IfNotNilString(dataFields.ReportedByOrganizationId),
		SubmittedByOrganizationId: utils.IfNotNilString(dataFields.SubmittedByOrganizationId),
		SubmittedByUserId:         utils.IfNotNilString(dataFields.SubmittedByUserId),
		Source:                    source.Source,
		AppSource:                 source.AppSource,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate IssueCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, IssueCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for IssueCreateEvent")
	}
	return event, nil
}

type IssueUpdateEvent struct {
	Tenant         string                `json:"tenant" validate:"required"`
	Subject        string                `json:"subject" validate:"required_without=Description"`
	Description    string                `json:"description" validate:"required_without=Subject"`
	Status         string                `json:"status"`
	Priority       string                `json:"priority"`
	UpdatedAt      time.Time             `json:"updatedAt"`
	Source         string                `json:"source"`
	ExternalSystem cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewIssueUpdateEvent(aggregate eventstore.Aggregate, dataFields model.IssueDataFields, source string, externalSystem cmnmod.ExternalSystem, updatedAt time.Time) (eventstore.Event, error) {
	eventData := IssueUpdateEvent{
		Tenant:      aggregate.GetTenant(),
		Subject:     dataFields.Subject,
		Description: dataFields.Description,
		Status:      dataFields.Status,
		Priority:    dataFields.Priority,
		UpdatedAt:   updatedAt,
		Source:      source,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate IssueUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, IssueUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for IssueUpdateEvent")
	}
	return event, nil
}

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
