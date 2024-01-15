package command

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type UpdateOrganizationPlanCommand struct {
	eventstore.BaseCommand
	UpdatedAt  *time.Time
	Name       string
	Retired    bool
	FieldsMask []string
}

func NewUpdateOrganizationPlanCommand(organizationPlanId, tenant, loggedInUserId, appSource, name string, retired bool, updatedAt *time.Time, fieldsMask []string) *UpdateOrganizationPlanCommand {
	return &UpdateOrganizationPlanCommand{
		BaseCommand: eventstore.NewBaseCommand(organizationPlanId, tenant, loggedInUserId).WithAppSource(appSource),
		UpdatedAt:   updatedAt,
		Name:        name,
		Retired:     retired,
		FieldsMask:  fieldsMask,
	}
}
