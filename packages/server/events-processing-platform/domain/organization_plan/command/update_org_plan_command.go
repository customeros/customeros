package command

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type UpdateOrgPlanCommand struct {
	eventstore.BaseCommand
	UpdatedAt  *time.Time
	Name       string
	Retired    bool
	FieldsMask []string
}

func NewUpdateOrgPlanCommand(orgPlanId, tenant, loggedInUserId, appSource, name string, retired bool, updatedAt *time.Time, fieldsMask []string) *UpdateOrgPlanCommand {
	return &UpdateOrgPlanCommand{
		BaseCommand: eventstore.NewBaseCommand(orgPlanId, tenant, loggedInUserId).WithAppSource(appSource),
		UpdatedAt:   updatedAt,
		Name:        name,
		Retired:     retired,
		FieldsMask:  fieldsMask,
	}
}
