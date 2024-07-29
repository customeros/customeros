package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type CreateMasterPlanCommand struct {
	eventstore.BaseCommand
	SourceFields common.Source
	CreatedAt    *time.Time
	Name         string
}

func NewCreateMasterPlanCommand(masterPlanId, tenant, loggedInUserId, name string, sourceFields common.Source, createdAt *time.Time) *CreateMasterPlanCommand {
	return &CreateMasterPlanCommand{
		BaseCommand:  eventstore.NewBaseCommand(masterPlanId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
		SourceFields: sourceFields,
		CreatedAt:    createdAt,
		Name:         name,
	}
}
