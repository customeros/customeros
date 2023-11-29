package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateServiceLineItem CreateServiceLineItemCommandHandler
	UpdateServiceLineItem UpdateServiceLineItemCommandHandler
	DeleteServiceLineItem DeleteServiceLineItemCommandHandler
	CloseServiceLineItem  CloseServiceLineItemCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateServiceLineItem: NewCreateServiceLineItemCommandHandler(log, es),
		UpdateServiceLineItem: NewUpdateServiceLineItemCommandHandler(log, es),
		DeleteServiceLineItem: NewDeleteServiceLineItemCommandHandler(log, es, cfg.Utils),
		CloseServiceLineItem:  NewCloseServiceLineItemCommandHandler(log, es, cfg.Utils),
	}
}
