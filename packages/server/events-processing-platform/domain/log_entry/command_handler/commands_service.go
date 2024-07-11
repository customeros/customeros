package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type CommandHandlers struct {
	UpsertLogEntry UpsertLogEntryCommandHandler
	AddTag         AddTagCommandHandler
	RemoveTag      RemoveTagCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertLogEntry: NewUpsertLogEntryCommandHandler(log, es, cfg.Utils),
		AddTag:         NewAddTagCommandHandler(log, es),
		RemoveTag:      NewRemoveTagCommandHandler(log, es),
	}
}
