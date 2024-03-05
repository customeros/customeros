package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateMasterPlanCommand struct {
	eventstore.BaseCommand
	SourceFields commonmodel.Source
	CreatedAt    *time.Time
	Name         string
}

func NewCreateMasterPlanCommand(masterPlanId, tenant, loggedInUserId, name string, sourceFields commonmodel.Source, createdAt *time.Time) *CreateMasterPlanCommand {
	return &CreateMasterPlanCommand{
		BaseCommand:  eventstore.NewBaseCommand(masterPlanId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
		SourceFields: sourceFields,
		CreatedAt:    createdAt,
		Name:         name,
	}
}
