package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateOrgPlan            CreateOrgPlanCommandHandler
	UpdateOrgPlan            UpdateOrgPlanCommandHandler
	CreateOrgPlanMilestone   CreateOrgPlanMilestoneCommandHandler
	UpdateOrgPlanMilestone   UpdateOrgPlanMilestoneCommandHandler
	ReorderOrgPlanMilestones ReorderOrgPlanMilestonesCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateOrgPlan:            NewCreateOrgPlanCommandHandler(log, es),
		UpdateOrgPlan:            NewUpdateOrgPlanCommandHandler(log, es, cfg.Utils),
		CreateOrgPlanMilestone:   NewCreateOrgPlanMilestoneCommandHandler(log, es, cfg.Utils),
		UpdateOrgPlanMilestone:   NewUpdateOrgPlanMilestoneCommandHandler(log, es, cfg.Utils),
		ReorderOrgPlanMilestones: NewReorderOrgPlanMilestonesCommandHandler(log, es, cfg.Utils),
	}
}
