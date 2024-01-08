package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateInvoicingCycle CreateInvoicingCycleCommandHandler
	UpdateInvoicingCycle UpdateInvoicingCycleCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateInvoicingCycle: NewCreateInvoicingCycleCommandHandler(log, es),
		UpdateInvoicingCycle: NewUpdateInvoicingCycleCommandHandler(log, es, cfg.Utils),
	}
}
