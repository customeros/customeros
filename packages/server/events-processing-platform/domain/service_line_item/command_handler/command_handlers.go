package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateServiceLineItem CreateServiceLineItemCommandHandler
	UpdateServiceLineItem UpdateServiceLineItemCommandHandler
}

func NewCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateServiceLineItem: NewCreateServiceLineItemCommandHandler(log, es),
		UpdateServiceLineItem: NewUpdateServiceLineItemCommandHandler(log, es),
	}
}
