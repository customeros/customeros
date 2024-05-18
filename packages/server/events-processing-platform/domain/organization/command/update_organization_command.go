package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateOrganizationCommand struct {
	eventstore.BaseCommand
	DataFields   model.OrganizationDataFields
	Source       string
	UpdatedAt    *time.Time
	FieldsMask   []string
	EnrichDomain string
	EnrichSource string
}

func NewUpdateOrganizationCommand(organizationId, tenant, loggedInUser, appSource, source string, dataFields model.OrganizationDataFields,
	updatedAt *time.Time, enrichDomain, enrichSource string, fieldsMask []string) *UpdateOrganizationCommand {
	return &UpdateOrganizationCommand{
		BaseCommand:  eventstore.NewBaseCommand(organizationId, tenant, loggedInUser).WithAppSource(appSource),
		DataFields:   dataFields,
		Source:       source,
		UpdatedAt:    updatedAt,
		FieldsMask:   fieldsMask,
		EnrichDomain: enrichDomain,
		EnrichSource: enrichSource,
	}
}
