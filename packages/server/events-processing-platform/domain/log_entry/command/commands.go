package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertLogEntryCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      models.LogEntryDataFields
	Source          cmnmod.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertLogEntryCommand(logEntryId, tenant, userId string, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, dataFields models.LogEntryDataFields, createdAt, updatedAt *time.Time) *UpsertLogEntryCommand {
	return &UpsertLogEntryCommand{
		BaseCommand:    eventstore.NewBaseCommand(logEntryId, tenant, userId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

type AddTagCommand struct {
	eventstore.BaseCommand
	TagId    string `json:"tagId" validate:"required"`
	TaggedAt *time.Time
}

func NewAddTagCommand(logEntryId, tenant, userId, tagId string, taggedAt *time.Time) *AddTagCommand {
	return &AddTagCommand{
		BaseCommand: eventstore.NewBaseCommand(logEntryId, tenant, userId),
		TagId:       tagId,
		TaggedAt:    taggedAt,
	}
}

type RemoveTagCommand struct {
	eventstore.BaseCommand
	TagId string
}

func NewRemoveTagCommand(logEntryId, tenant, userId, tagId string) *RemoveTagCommand {
	return &RemoveTagCommand{
		BaseCommand: eventstore.NewBaseCommand(logEntryId, tenant, userId),
		TagId:       tagId,
	}
}
