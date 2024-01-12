package command

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type OrgPlanMilestoneItems struct {
	Text      []string
	UpdatedAt *time.Time
	Status    string
}

type UpdateOrgPlanMilestoneCommand struct {
	eventstore.BaseCommand
	MilestoneId   string `validate:"required"`
	UpdatedAt     *time.Time
	AppSource     string
	Name          string
	Order         int64
	DurationHours int64
	Items         []OrgPlanMilestoneItems
	Optional      bool
	Retired       bool
	FieldsMask    []string
}

func NewUpdateOrgPlanMilestoneCommand(orgPlanId, tenant, loggedInUserId, milestoneId, name, appSource string, order, durationHours int64, items []OrgPlanMilestoneItems, optional, retired bool, updatedAt *time.Time, fieldsMask []string) *UpdateOrgPlanMilestoneCommand {
	return &UpdateOrgPlanMilestoneCommand{
		BaseCommand:   eventstore.NewBaseCommand(orgPlanId, tenant, loggedInUserId).WithAppSource(appSource),
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
