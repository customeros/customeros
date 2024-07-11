package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type UpsertLogEntryCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      model.LogEntryDataFields
	Source          events.Source
	ExternalSystem  commonmodel.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertLogEntryCommand(logEntryId, tenant, userId string, source events.Source, externalSystem commonmodel.ExternalSystem, dataFields model.LogEntryDataFields, createdAt, updatedAt *time.Time) *UpsertLogEntryCommand {
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
