package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/model"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type UpsertInteractionSessionCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      model.InteractionSessionDataFields
	Source          cmnmod.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertInteractionSessionCommand(interactionSessionId, tenant, loggedInUserId string, dataFields model.InteractionSessionDataFields, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt *time.Time) *UpsertInteractionSessionCommand {
	return &UpsertInteractionSessionCommand{
		BaseCommand:    eventstore.NewBaseCommand(interactionSessionId, tenant, loggedInUserId).WithAppSource(source.AppSource),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
