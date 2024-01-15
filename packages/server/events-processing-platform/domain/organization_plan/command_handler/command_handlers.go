package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	CreateOrganizationPlan            CreateOrganizationPlanCommandHandler
	UpdateOrganizationPlan            UpdateOrganizationPlanCommandHandler
	CreateOrganizationPlanMilestone   CreateOrganizationPlanMilestoneCommandHandler
	UpdateOrganizationPlanMilestone   UpdateOrganizationPlanMilestoneCommandHandler
	ReorderOrganizationPlanMilestones ReorderOrganizationPlanMilestonesCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateOrganizationPlan:            NewCreateOrganizationPlanCommandHandler(log, es),
		UpdateOrganizationPlan:            NewUpdateOrganizationPlanCommandHandler(log, es, cfg.Utils),
		CreateOrganizationPlanMilestone:   NewCreateOrganizationPlanMilestoneCommandHandler(log, es, cfg.Utils),
		UpdateOrganizationPlanMilestone:   NewUpdateOrganizationPlanMilestoneCommandHandler(log, es, cfg.Utils),
		ReorderOrganizationPlanMilestones: NewReorderOrganizationPlanMilestonesCommandHandler(log, es, cfg.Utils),
	}
}
