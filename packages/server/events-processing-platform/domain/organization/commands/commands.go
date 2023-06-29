package commands

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertOrganizationCommand struct {
	eventstore.BaseCommand
	CoreFields models.OrganizationCoreFields
	Source     common_models.Source
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func UpsertOrganizationCommandToOrganizationDto(command *UpsertOrganizationCommand) *models.OrganizationDto {
	return &models.OrganizationDto{
		ID:                     command.ObjectID,
		Tenant:                 command.Tenant,
		OrganizationCoreFields: command.CoreFields,
		Source:                 command.Source,
		CreatedAt:              command.CreatedAt,
		UpdatedAt:              command.UpdatedAt,
	}
}

func NewUpsertOrganizationCommand(organizationId, tenant, source, sourceOfTruth, appSource string, coreFields models.OrganizationCoreFields, createdAt, updatedAt *time.Time) *UpsertOrganizationCommand {
	return &UpsertOrganizationCommand{
		BaseCommand: eventstore.NewBaseCommand(organizationId, tenant),
		CoreFields:  coreFields,
		Source: common_models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type LinkPhoneNumberCommand struct {
	eventstore.BaseCommand
	PhoneNumberId string
	Primary       bool
	Label         string
}

func NewLinkPhoneNumberCommand(objectID, tenant, phoneNumberId, label string, primary bool) *LinkPhoneNumberCommand {
	return &LinkPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(objectID, tenant),
		PhoneNumberId: phoneNumberId,
		Primary:       primary,
		Label:         label,
	}
}

type LinkEmailCommand struct {
	eventstore.BaseCommand
	EmailId string
	Primary bool
	Label   string
}

func NewLinkEmailCommand(objectID, tenant, emailId, label string, primary bool) *LinkEmailCommand {
	return &LinkEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant),
		EmailId:     emailId,
		Primary:     primary,
		Label:       label,
	}
}
