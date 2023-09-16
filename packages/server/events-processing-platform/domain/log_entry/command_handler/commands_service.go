package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type LogEntryCommands struct {
	UpsertLogEntry UpsertLogEntryCommandHandler
}

func NewLogEntryCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *LogEntryCommands {
	return &LogEntryCommands{
		UpsertLogEntry: NewUpsertLogEntryCommandHandler(log, cfg, es),
	}
}
