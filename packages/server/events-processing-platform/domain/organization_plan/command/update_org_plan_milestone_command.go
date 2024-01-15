package command

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type OrganizationPlanMilestoneItems struct {
	Text      []string
	UpdatedAt *time.Time
	Status    string
}

type UpdateOrganizationPlanMilestoneCommand struct {
	eventstore.BaseCommand
	MilestoneId   string `validate:"required"`
	UpdatedAt     *time.Time
	AppSource     string
	Name          string
	Order         int64
	DurationHours int64
	Items         []OrganizationPlanMilestoneItems
	Optional      bool
	Retired       bool
	FieldsMask    []string
}

func NewUpdateOrganizationPlanMilestoneCommand(organizationPlanId, tenant, loggedInUserId, milestoneId, name, appSource string, order, durationHours int64, items []OrganizationPlanMilestoneItems, optional, retired bool, updatedAt *time.Time, fieldsMask []string) *UpdateOrganizationPlanMilestoneCommand {
	return &UpdateOrganizationPlanMilestoneCommand{
		BaseCommand:   eventstore.NewBaseCommand(organizationPlanId, tenant, loggedInUserId).WithAppSource(appSource),
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
