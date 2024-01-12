package command

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type CreateOrgPlanCommand struct {
	eventstore.BaseCommand
	SourceFields commonmodel.Source
	CreatedAt    *time.Time
	Name         string
}

func NewCreateOrgPlanCommand(orgPlanId, tenant, loggedInUserId, name string, sourceFields commonmodel.Source, createdAt *time.Time) *CreateOrgPlanCommand {
	return &CreateOrgPlanCommand{
		BaseCommand:  eventstore.NewBaseCommand(orgPlanId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
		SourceFields: sourceFields,
		CreatedAt:    createdAt,
		Name:         name,
	}
}
