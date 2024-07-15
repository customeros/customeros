package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type CreateJobRoleCommand struct {
	eventstore.BaseCommand
	StartedAt   *time.Time
	EndedAt     *time.Time
	JobTitle    string
	Description *string
	Primary     bool
	Source      common.Source
	CreatedAt   *time.Time
}

func NewCreateJobRoleCommand(objectID, tenant, jobTitle string, description *string, primary bool, source, sourceOfTruth, appSource string, startedAt, endedAt, createdAt *time.Time) *CreateJobRoleCommand {
	return &CreateJobRoleCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, ""),
		StartedAt:   startedAt,
		EndedAt:     endedAt,
		JobTitle:    jobTitle,
		Description: description,
		Primary:     primary,
		Source: common.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
	}
}
