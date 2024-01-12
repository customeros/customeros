package command

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type ReorderOrgPlanMilestonesCommand struct {
	eventstore.BaseCommand
	MilestoneIds []string `validate:"required"`
	UpdatedAt    *time.Time
	AppSource    string
}

func NewReorderOrgPlanMilestonesCommand(orgPlanId, tenant, loggedInUserId, appSource string, milestoneIds []string, updatedAt *time.Time) *ReorderOrgPlanMilestonesCommand {
	return &ReorderOrgPlanMilestonesCommand{
		BaseCommand:  eventstore.NewBaseCommand(orgPlanId, tenant, loggedInUserId).WithAppSource(appSource),
		MilestoneIds: milestoneIds,
		UpdatedAt:    updatedAt,
	}
}
