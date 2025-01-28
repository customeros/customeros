package command_handler

import (
	"github.com/customeros/customeros/packages/server/events/eventstore"

	"github.com/customeros/customeros/packages/server/events-processing-platform/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	RefreshArr RefreshArrCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		RefreshArr: NewRefreshArrCommandHandler(log, es, cfg.Utils),
	}
}
