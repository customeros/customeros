package command

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type ReorderOrganizationPlanMilestonesCommand struct {
	eventstore.BaseCommand
	MilestoneIds []string `validate:"required"`
	UpdatedAt    *time.Time
	AppSource    string
}

func NewReorderOrganizationPlanMilestonesCommand(organizationPlanId, tenant, loggedInUserId, appSource string, milestoneIds []string, updatedAt *time.Time) *ReorderOrganizationPlanMilestonesCommand {
	return &ReorderOrganizationPlanMilestonesCommand{
		BaseCommand:  eventstore.NewBaseCommand(organizationPlanId, tenant, loggedInUserId).WithAppSource(appSource),
		MilestoneIds: milestoneIds,
		UpdatedAt:    updatedAt,
	}
}
