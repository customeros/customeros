package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateMasterPlanMilestoneCommand struct {
	eventstore.BaseCommand
	MilestoneId   string `validate:"required"`
	UpdatedAt     *time.Time
	AppSource     string
	Name          string
	Order         int64
	DurationHours int64
	Items         []string
	Optional      bool
	Retired       bool
	FieldsMask    []string
}

func NewUpdateMasterPlanMilestoneCommand(masterPlanId, tenant, loggedInUserId, milestoneId, name, appSource string, order, durationHours int64, items []string, optional, retired bool, updatedAt *time.Time, fieldsMask []string) *UpdateMasterPlanMilestoneCommand {
	return &UpdateMasterPlanMilestoneCommand{
		BaseCommand:   eventstore.NewBaseCommand(masterPlanId, tenant, loggedInUserId).WithAppSource(appSource),
		MilestoneId:   milestoneId,
		UpdatedAt:     updatedAt,
		Name:          name,
		Order:         order,
		DurationHours: durationHours,
		Items:         items,
		Optional:      optional,
		Retired:       retired,
		FieldsMask:    fieldsMask,
	}
}
