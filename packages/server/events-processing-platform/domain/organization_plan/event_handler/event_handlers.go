package event_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

// EventHandlers acts as a container for all command handlers.
type EventHandlers struct {
	CreateOrganizationPlan            CreateOrganizationPlanHandler
	UpdateOrganizationPlan            UpdateOrganizationPlanHandler
	CreateOrganizationPlanMilestone   CreateOrganizationPlanMilestoneHandler
	UpdateOrganizationPlanMilestone   UpdateOrganizationPlanMilestoneCommandHandler
	ReorderOrganizationPlanMilestones ReorderOrganizationPlanMilestonesHandler
}

func NewEventHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *EventHandlers {
	return &EventHandlers{
		CreateOrganizationPlan:            NewCreateOrganizationPlanHandler(log, es),
		UpdateOrganizationPlan:            NewUpdateOrganizationPlanHandler(log, es, cfg.Utils),
		CreateOrganizationPlanMilestone:   NewCreateOrganizationPlanMilestoneHandler(log, es, cfg.Utils),
		UpdateOrganizationPlanMilestone:   NewUpdateOrganizationPlanMilestoneCommandHandler(log, es, cfg.Utils),
		ReorderOrganizationPlanMilestones: NewReorderOrganizationPlanMilestonesHandler(log, es, cfg.Utils),
	}
}
