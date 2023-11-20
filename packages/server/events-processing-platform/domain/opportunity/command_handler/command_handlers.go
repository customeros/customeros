package command_handler

import (
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
}

func NewCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateOpportunity:                     NewCreateOpportunityCommandHandler(log, es),
		UpdateOpportunity:                     NewUpdateOpportunityCommandHandler(log, es),
		CreateRenewalOpportunity:              NewCreateRenewalOpportunityCommandHandler(log, es),
		UpdateRenewalOpportunity:              NewUpdateRenewalOpportunityCommandHandler(log, es),
		UpdateRenewalOpportunityNextCycleDate: NewUpdateRenewalOpportunityNextCycleDateCommandHandler(log, es),
	}
}
