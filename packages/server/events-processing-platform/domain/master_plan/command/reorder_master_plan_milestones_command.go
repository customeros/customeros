package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type ReorderMasterPlanMilestonesCommand struct {
	eventstore.BaseCommand
	MilestoneIds []string `validate:"required"`
	UpdatedAt    *time.Time
	AppSource    string
}

func NewReorderMasterPlanMilestonesCommand(masterPlanId, tenant, loggedInUserId, appSource string, milestoneIds []string, updatedAt *time.Time) *ReorderMasterPlanMilestonesCommand {
	return &ReorderMasterPlanMilestonesCommand{
		BaseCommand:  eventstore.NewBaseCommand(masterPlanId, tenant, loggedInUserId).WithAppSource(appSource),
		MilestoneIds: milestoneIds,
		UpdatedAt:    updatedAt,
	}
}
