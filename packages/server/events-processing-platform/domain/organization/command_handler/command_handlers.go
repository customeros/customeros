package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	UpsertCustomFieldCommand UpsertCustomFieldCommandHandler
	RefreshArr               RefreshArrCommandHandler
	UpdateOrganizationOwner  UpdateOrganizationOwnerCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore, ebs *eventbuffer.EventBufferStoreService) *CommandHandlers {
	return &CommandHandlers{
		UpsertCustomFieldCommand: NewUpsertCustomFieldCommandHandler(log, es),
		RefreshArr:               NewRefreshArrCommandHandler(log, es, *cfg.Utils),
		UpdateOrganizationOwner:  NewUpdateOrganizationOwnerCommandHandler(log, es, *cfg.Utils, ebs),
	}
}
