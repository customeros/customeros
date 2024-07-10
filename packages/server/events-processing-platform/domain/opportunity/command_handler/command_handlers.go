package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	UpdateRenewalOpportunityNextCycleDate UpdateRenewalOpportunityNextCycleDateCommandHandler
	CloseWinOpportunity                   CloseWinOpportunityCommandHandler
	CloseLooseOpportunity                 CloseLooseOpportunityCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpdateRenewalOpportunityNextCycleDate: NewUpdateRenewalOpportunityNextCycleDateCommandHandler(log, es, cfg.Utils),
		CloseWinOpportunity:                   NewCloseWinOpportunityCommandHandler(log, es, cfg.Utils),
		CloseLooseOpportunity:                 NewCloseLooseOpportunityCommandHandler(log, es, cfg.Utils),
	}
}
