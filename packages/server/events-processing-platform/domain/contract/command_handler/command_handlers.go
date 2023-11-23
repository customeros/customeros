package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateContract                        CreateContractCommandHandler
	UpdateContract                        UpdateContractCommandHandler
	RefreshContractStatus                 RefreshContractStatusCommandHandler
	RolloutRenewalOpportunityOnExpiration RolloutRenewalOpportunityOnExpirationCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateContract:                        NewCreateContractCommandHandler(log, es),
		UpdateContract:                        NewUpdateContractCommandHandler(log, es, cfg.Utils),
		RefreshContractStatus:                 NewRefreshContractStatusCommandHandler(log, es, cfg.Utils),
		RolloutRenewalOpportunityOnExpiration: NewRolloutRenewalOpportunityOnExpirationCommandHandler(log, es, cfg.Utils),
	}
}
