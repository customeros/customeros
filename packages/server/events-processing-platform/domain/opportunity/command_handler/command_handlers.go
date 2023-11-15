package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateOpportunity                     CreateOpportunityCommandHandler
	CreateRenewalOpportunity              CreateRenewalOpportunityCommandHandler
	UpdateRenewalOpportunityNextCycleDate UpdateRenewalOpportunityNextCycleDateCommandHandler
}

func NewCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateOpportunity:                     NewCreateOpportunityCommandHandler(log, es),
		CreateRenewalOpportunity:              NewCreateRenewalOpportunityCommandHandler(log, es),
		UpdateRenewalOpportunityNextCycleDate: NewUpdateRenewalOpportunityNextCycleDateCommandHandler(log, es),
	}
}
