package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type LogEntryCommands struct {
	UpsertLogEntry UpsertLogEntryCommandHandler
	AddTag         AddTagCommandHandler
	RemoveTag      RemoveTagCommandHandler
}

func NewLogEntryCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *LogEntryCommands {
	return &LogEntryCommands{
		UpsertLogEntry: NewUpsertLogEntryCommandHandler(log, cfg, es),
		AddTag:         NewAddTagCommandHandler(log, cfg, es),
		RemoveTag:      NewRemoveTagCommandHandler(log, cfg, es),
	}
}
