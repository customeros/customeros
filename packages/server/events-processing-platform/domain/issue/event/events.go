package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	IssueCreateV1 = "V1_ISSUE_CREATE"
	IssueUpdateV1 = "V1_ISSUE_UPDATE"
)

type IssueCreateEvent struct {
	Tenant                   string                `json:"tenant" validate:"required"`
	Subject                  string                `json:"subject" validate:"required_without=Description"`
	Description              string                `json:"description" validate:"required_without=Subject"`
	Status                   string                `json:"status"`
	Priority                 string                `json:"priority"`
	ReportedByOrganizationId string                `json:"loggedOrganizationId,omitempty"`
	Source                   string                `json:"source"`
	AppSource                string                `json:"appSource"`
	CreatedAt                time.Time             `json:"createdAt"`
	UpdatedAt                time.Time             `json:"updatedAt"`
	ExternalSystem           cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewIssueCreateEvent(aggregate eventstore.Aggregate, dataFields model.IssueDataFields, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := IssueCreateEvent{
		Tenant:                   aggregate.GetTenant(),
		Subject:                  dataFields.Subject,
		Description:              dataFields.Description,
		Status:                   dataFields.Status,
		Priority:                 dataFields.Priority,
		ReportedByOrganizationId: utils.IfNotNilString(dataFields.ReportedByOrganizationId),
		Source:                   source.Source,
		AppSource:                source.AppSource,
		CreatedAt:                createdAt,
		UpdatedAt:                updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, IssueCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
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
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, IssueUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
