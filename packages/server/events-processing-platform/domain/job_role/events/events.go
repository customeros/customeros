package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

const (
	JobRoleCreateV1 = "V1_JOB_ROLE_CREATE"
)

type JobRoleCreateEvent struct {
	Tenant        string     `json:"tenant" validate:"required"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	EndedAt       *time.Time `json:"endedAt,omitempty"`
	JobTitle      string     `json:"jobTitle" validate:"required"`
	Description   *string    `json:"description,omitempty"`
	Primary       bool       `json:"primary"`
	Source        string     `json:"source"`
	SourceOfTruth string     `json:"sourceOfTruth"`
	AppSource     string     `json:"appSource"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type JobRoleUpdateEvent struct {
	Tenant        string     `json:"tenant" validate:"required"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	EndedAt       *time.Time `json:"endedAt,omitempty"`
	JobTitle      *string    `json:"jobTitle,omitempty"`
	Description   *string    `json:"description"`
	Primary       *bool      `json:"primary"`
	SourceOfTruth string     `json:"sourceOfTruth"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func NewJobRoleCreateEvent(aggregate eventstore.Aggregate, command *model.CreateJobRoleCommand) (eventstore.Event, error) {
	createdAt := utils.IfNotNilTimeWithDefault(command.CreatedAt, utils.Now())

	eventData := &JobRoleCreateEvent{
		Tenant:        command.Tenant,
		StartedAt:     command.StartedAt,
		EndedAt:       command.EndedAt,
		JobTitle:      command.JobTitle,
		Description:   command.Description,
		Primary:       command.Primary,
		Source:        command.Source.Source,
		SourceOfTruth: command.Source.SourceOfTruth,
		AppSource:     command.Source.AppSource,
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
	}
	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate JobRoleCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, JobRoleCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for JobRoleCreateEvent")
	}
	return event, nil
}
