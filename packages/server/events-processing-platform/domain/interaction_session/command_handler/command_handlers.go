package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type CommandHandlers struct {
	UpsertInteractionSession UpsertInteractionSessionCommandHandler
}

func NewCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertInteractionSession: NewUpsertInteractionSessionCommandHandler(log, es),
	}
}
