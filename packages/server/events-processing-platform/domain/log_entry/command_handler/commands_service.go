package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommandHandlers struct {
	UpsertLogEntry UpsertLogEntryCommandHandler
	AddTag         AddTagCommandHandler
	RemoveTag      RemoveTagCommandHandler
}

func NewCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertLogEntry: NewUpsertLogEntryCommandHandler(log, es),
		AddTag:         NewAddTagCommandHandler(log, es),
		RemoveTag:      NewRemoveTagCommandHandler(log, es),
	}
}
