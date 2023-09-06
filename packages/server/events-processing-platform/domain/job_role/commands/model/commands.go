package model

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateJobRoleCommand struct {
	eventstore.BaseCommand
	StartedAt   *time.Time
	EndedAt     *time.Time
	JobTitle    string
	Description *string
	Primary     bool
	Source      common_models.Source
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
		Source: common_models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
	}
}
