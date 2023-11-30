package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateOpportunity                     CreateOpportunityCommandHandler
	UpdateOpportunity                     UpdateOpportunityCommandHandler
	CreateRenewalOpportunity              CreateRenewalOpportunityCommandHandler
	UpdateRenewalOpportunity              UpdateRenewalOpportunityCommandHandler
	UpdateRenewalOpportunityNextCycleDate UpdateRenewalOpportunityNextCycleDateCommandHandler
	CloseWinOpportunity                   CloseWinOpportunityCommandHandler
	CloseLooseOpportunity                 CloseLooseOpportunityCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateOpportunity:                     NewCreateOpportunityCommandHandler(log, es),
		UpdateOpportunity:                     NewUpdateOpportunityCommandHandler(log, es, cfg.Utils),
		CreateRenewalOpportunity:              NewCreateRenewalOpportunityCommandHandler(log, es),
		UpdateRenewalOpportunity:              NewUpdateRenewalOpportunityCommandHandler(log, es, cfg.Utils),
		UpdateRenewalOpportunityNextCycleDate: NewUpdateRenewalOpportunityNextCycleDateCommandHandler(log, es, cfg.Utils),
		CloseWinOpportunity:                   NewCloseWinOpportunityCommandHandler(log, es, cfg.Utils),
		CloseLooseOpportunity:                 NewCloseLooseOpportunityCommandHandler(log, es, cfg.Utils),
	}
}
