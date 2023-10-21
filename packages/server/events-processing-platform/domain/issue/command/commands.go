package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertIssueCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      model.IssueDataFields
	Source          cmnmod.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertIssueCommand(issueId, tenant, loggedInUserId string, dataFields model.IssueDataFields, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt *time.Time) *UpsertIssueCommand {
	return &UpsertIssueCommand{
		BaseCommand:    eventstore.NewBaseCommand(issueId, tenant, loggedInUserId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
