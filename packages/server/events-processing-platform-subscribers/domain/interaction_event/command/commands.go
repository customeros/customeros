package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertInteractionEventCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      model.InteractionEventDataFields
	Source          cmnmod.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertInteractionEventCommand(interactionEventId, tenant, loggedInUserId string, dataFields model.InteractionEventDataFields, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt *time.Time) *UpsertInteractionEventCommand {
	return &UpsertInteractionEventCommand{
		BaseCommand:    eventstore.NewBaseCommand(interactionEventId, tenant, loggedInUserId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

type ReplaceActionItemsCommand struct {
	eventstore.BaseCommand
	ActionItems []string
	UpdatedAt   *time.Time
}

func NewReplaceActionItemsCommand(tenant, interactionEventId string, actionItems []string, updatedAt *time.Time) *ReplaceActionItemsCommand {
	return &ReplaceActionItemsCommand{
		BaseCommand: eventstore.NewBaseCommand(interactionEventId, tenant, ""),
		ActionItems: actionItems,
		UpdatedAt:   updatedAt,
	}
}

type ReplaceSummaryCommand struct {
	eventstore.BaseCommand
	Summary     string
	ContentType string
	UpdatedAt   *time.Time
}

func NewReplaceSummaryCommand(tenant, interactionEventId, summary, contentType string, updatedAt *time.Time) *ReplaceSummaryCommand {
	return &ReplaceSummaryCommand{
		BaseCommand: eventstore.NewBaseCommand(interactionEventId, tenant, ""),
		Summary:     summary,
		ContentType: contentType,
		UpdatedAt:   updatedAt,
	}
}

type RequestActionItemsCommand struct {
	eventstore.BaseCommand
}

func NewRequestActionItemsCommand(tenant, interactionEventId string) *RequestActionItemsCommand {
	return &RequestActionItemsCommand{
		BaseCommand: eventstore.NewBaseCommand(interactionEventId, tenant, ""),
	}
}

type RequestSummaryCommand struct {
	eventstore.BaseCommand
}

func NewRequestSummaryCommand(tenant, interactionEventId string) *RequestSummaryCommand {
	return &RequestSummaryCommand{
		BaseCommand: eventstore.NewBaseCommand(interactionEventId, tenant, ""),
	}
}
