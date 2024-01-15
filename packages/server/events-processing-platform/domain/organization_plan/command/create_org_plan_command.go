package command

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type CreateOrganizationPlanCommand struct {
	eventstore.BaseCommand
	SourceFields commonmodel.Source
	CreatedAt    *time.Time
	Name         string
}

func NewCreateOrganizationPlanCommand(organizationPlanId, tenant, loggedInUserId, name string, sourceFields commonmodel.Source, createdAt *time.Time) *CreateOrganizationPlanCommand {
	return &CreateOrganizationPlanCommand{
		BaseCommand:  eventstore.NewBaseCommand(organizationPlanId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
		SourceFields: sourceFields,
		CreatedAt:    createdAt,
		Name:         name,
	}
}
