package command

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertOrganizationCommand struct {
	eventstore.BaseCommand
	CoreFields models.OrganizationDataFields
	Source     common_models.Source
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func UpsertOrganizationCommandToOrganizationFields(command *UpsertOrganizationCommand) *models.OrganizationFields {
	return &models.OrganizationFields{
		ID:                     command.ObjectID,
		Tenant:                 command.Tenant,
		OrganizationDataFields: command.CoreFields,
		Source:                 command.Source,
		CreatedAt:              command.CreatedAt,
		UpdatedAt:              command.UpdatedAt,
	}
}

func NewUpsertOrganizationCommand(organizationId, tenant, source, sourceOfTruth, appSource string, coreFields models.OrganizationDataFields, createdAt, updatedAt *time.Time) *UpsertOrganizationCommand {
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

type UpdateOrganizationCommand struct {
	eventstore.BaseCommand
	IgnoreEmptyFields bool `json:"ignoreEmptyFields"`
	DataFields        models.OrganizationDataFields
	SourceOfTruth     string `json:"sourceOfTruth"`
	UpdatedAt         *time.Time
}

func NewUpdateOrganizationCommand(organizationId, tenant, sourceOfTruth string, dataFields models.OrganizationDataFields, updatedAt *time.Time, ignoreEmptyFields bool) *UpdateOrganizationCommand {
	return &UpdateOrganizationCommand{
		BaseCommand:       eventstore.NewBaseCommand(organizationId, tenant),
		IgnoreEmptyFields: ignoreEmptyFields,
		DataFields:        dataFields,
		SourceOfTruth:     sourceOfTruth,
		UpdatedAt:         updatedAt,
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

type UpdateRenewalLikelihoodCommand struct {
	eventstore.BaseCommand
	Fields models.RenewalLikelihoodFields
}

func NewUpdateRenewalLikelihoodCommand(tenant, orgId string, fields models.RenewalLikelihoodFields) *UpdateRenewalLikelihoodCommand {
	return &UpdateRenewalLikelihoodCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant),
		Fields:      fields,
	}
}

type RequestNextCycleDateCommand struct {
	eventstore.BaseCommand
}

func NewRequestNextCycleDateCommand(tenant, orgId string) *RequestNextCycleDateCommand {
	return &RequestNextCycleDateCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant),
	}
}

type RequestRenewalForecastCommand struct {
	eventstore.BaseCommand
}

func NewRequestRenewalForecastCommand(tenant, orgId string) *RequestRenewalForecastCommand {
	return &RequestRenewalForecastCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant),
	}
}

type UpdateRenewalForecastCommand struct {
	eventstore.BaseCommand
	Fields            models.RenewalForecastFields
	RenewalLikelihood models.RenewalLikelihoodProbability
}

func NewUpdateRenewalForecastCommand(tenant, orgId string, fields models.RenewalForecastFields, renewalLikelihood models.RenewalLikelihoodProbability) *UpdateRenewalForecastCommand {
	return &UpdateRenewalForecastCommand{
		BaseCommand:       eventstore.NewBaseCommand(orgId, tenant),
		Fields:            fields,
		RenewalLikelihood: renewalLikelihood,
	}
}

type UpdateBillingDetailsCommand struct {
	eventstore.BaseCommand
	Fields models.BillingDetailsFields
}

func NewUpdateBillingDetailsCommand(tenant, orgId string, fields models.BillingDetailsFields) *UpdateBillingDetailsCommand {
	return &UpdateBillingDetailsCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant),
		Fields:      fields,
	}
}
