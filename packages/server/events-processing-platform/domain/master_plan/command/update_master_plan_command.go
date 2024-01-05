package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateMasterPlanCommand struct {
	eventstore.BaseCommand
	UpdatedAt  *time.Time
	Name       string
	Retired    bool
	FieldsMask []string
}

func NewUpdateMasterPlanCommand(masterPlanId, tenant, loggedInUserId, appSource, name string, retired bool, updatedAt *time.Time, fieldsMask []string) *UpdateMasterPlanCommand {
	return &UpdateMasterPlanCommand{
		BaseCommand: eventstore.NewBaseCommand(masterPlanId, tenant, loggedInUserId).WithAppSource(appSource),
		UpdatedAt:   updatedAt,
		Name:        name,
		Retired:     retired,
		FieldsMask:  fieldsMask,
	}
}
