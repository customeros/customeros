package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type LogEntryCommandHandlers struct {
	UpsertLogEntry UpsertLogEntryCommandHandler
	AddTag         AddTagCommandHandler
	RemoveTag      RemoveTagCommandHandler
}

func NewLogEntryCommands(log logger.Logger, es eventstore.AggregateStore) *LogEntryCommandHandlers {
	return &LogEntryCommandHandlers{
		UpsertLogEntry: NewUpsertLogEntryCommandHandler(log, es),
		AddTag:         NewAddTagCommandHandler(log, es),
		RemoveTag:      NewRemoveTagCommandHandler(log, es),
	}
}
