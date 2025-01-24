package command_handler

import (
	"github.com/customeros/customeros/packages/server/events/eventbuffer"
	"github.com/customeros/customeros/packages/server/events/eventstore"

	"github.com/customeros/customeros/packages/server/events-processing-platform/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	UpsertCustomFieldCommand UpsertCustomFieldCommandHandler
	RefreshArr               RefreshArrCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore, ebs *eventbuffer.EventBufferStoreService) *CommandHandlers {
	return &CommandHandlers{
		UpsertCustomFieldCommand: NewUpsertCustomFieldCommandHandler(log, es),
		RefreshArr:               NewRefreshArrCommandHandler(log, es, cfg.Utils),
	}
}
