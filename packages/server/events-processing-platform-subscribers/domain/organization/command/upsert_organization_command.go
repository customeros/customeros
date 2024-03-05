package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertOrganizationCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      model.OrganizationDataFields
	Source          cmnmod.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	FieldsMask      []string
}

func UpsertOrganizationCommandToOrganizationFieldsStruct(command *UpsertOrganizationCommand) *model.OrganizationFields {
	return &model.OrganizationFields{
		ID:                     command.ObjectID,
		Tenant:                 command.Tenant,
		OrganizationDataFields: command.DataFields,
		Source:                 command.Source,
		ExternalSystem:         command.ExternalSystem,
		CreatedAt:              command.CreatedAt,
		UpdatedAt:              command.UpdatedAt,
	}
}

func NewUpsertOrganizationCommand(organizationId, tenant, loggedInUserId string, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, coreFields model.OrganizationDataFields, createdAt, updatedAt *time.Time, maskFields []string) *UpsertOrganizationCommand {
	return &UpsertOrganizationCommand{
		BaseCommand:    eventstore.NewBaseCommand(organizationId, tenant, loggedInUserId).WithAppSource(source.AppSource),
		DataFields:     coreFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		FieldsMask:     maskFields,
	}
}
