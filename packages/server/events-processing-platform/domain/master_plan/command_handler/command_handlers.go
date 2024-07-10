package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateMasterPlan            CreateMasterPlanCommandHandler
	UpdateMasterPlan            UpdateMasterPlanCommandHandler
	CreateMasterPlanMilestone   CreateMasterPlanMilestoneCommandHandler
	UpdateMasterPlanMilestone   UpdateMasterPlanMilestoneCommandHandler
	ReorderMasterPlanMilestones ReorderMasterPlanMilestonesCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateMasterPlan:            NewCreateMasterPlanCommandHandler(log, es),
		UpdateMasterPlan:            NewUpdateMasterPlanCommandHandler(log, es, cfg.Utils),
		CreateMasterPlanMilestone:   NewCreateMasterPlanMilestoneCommandHandler(log, es, cfg.Utils),
		UpdateMasterPlanMilestone:   NewUpdateMasterPlanMilestoneCommandHandler(log, es, cfg.Utils),
		ReorderMasterPlanMilestones: NewReorderMasterPlanMilestonesCommandHandler(log, es, cfg.Utils),
	}
}
