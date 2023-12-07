package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateOrganizationCommand struct {
	eventstore.BaseCommand
	DataFields model.OrganizationDataFields
	Source     string
	UpdatedAt  *time.Time
	FieldsMask []string
}

func NewUpdateOrganizationCommand(organizationId, tenant, source string, dataFields model.OrganizationDataFields, updatedAt *time.Time, fieldsMask []string) *UpdateOrganizationCommand {
	return &UpdateOrganizationCommand{
		BaseCommand: eventstore.NewBaseCommand(organizationId, tenant, ""),
		DataFields:  dataFields,
		Source:      source,
		UpdatedAt:   updatedAt,
		FieldsMask:  fieldsMask,
	}
}
