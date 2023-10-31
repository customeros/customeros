package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertOrganizationCommand struct {
	eventstore.BaseCommand
	IsCreateCommand   bool
	IgnoreEmptyFields bool
	DataFields        models.OrganizationDataFields
	Source            cmnmod.Source
	ExternalSystem    cmnmod.ExternalSystem
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}

func UpsertOrganizationCommandToOrganizationFieldsStruct(command *UpsertOrganizationCommand) *models.OrganizationFields {
	return &models.OrganizationFields{
		ID:                     command.ObjectID,
		Tenant:                 command.Tenant,
		OrganizationDataFields: command.DataFields,
		Source:                 command.Source,
		ExternalSystem:         command.ExternalSystem,
		CreatedAt:              command.CreatedAt,
		UpdatedAt:              command.UpdatedAt,
		IgnoreEmptyFields:      command.IgnoreEmptyFields,
	}
}

func NewUpsertOrganizationCommand(organizationId, tenant, userId string, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, coreFields models.OrganizationDataFields, createdAt, updatedAt *time.Time, ignoreEmptyFields bool) *UpsertOrganizationCommand {
	return &UpsertOrganizationCommand{
		BaseCommand:       eventstore.NewBaseCommand(organizationId, tenant, userId),
		IgnoreEmptyFields: ignoreEmptyFields,
		DataFields:        coreFields,
		Source:            source,
		ExternalSystem:    externalSystem,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

type UpdateOrganizationCommand struct {
	eventstore.BaseCommand
	IgnoreEmptyFields bool
	DataFields        models.OrganizationDataFields
	Source            string
	UpdatedAt         *time.Time
}

func NewUpdateOrganizationCommand(organizationId, tenant, source string, dataFields models.OrganizationDataFields, updatedAt *time.Time, ignoreEmptyFields bool) *UpdateOrganizationCommand {
	return &UpdateOrganizationCommand{
		BaseCommand:       eventstore.NewBaseCommand(organizationId, tenant, ""),
		IgnoreEmptyFields: ignoreEmptyFields,
		DataFields:        dataFields,
		Source:            source,
		UpdatedAt:         updatedAt,
	}
}

type LinkPhoneNumberCommand struct {
	eventstore.BaseCommand
	PhoneNumberId string `json:"phoneNumberId" validate:"required"`
	Primary       bool
	Label         string
}

func NewLinkPhoneNumberCommand(objectID, tenant, userId, phoneNumberId, label string, primary bool) *LinkPhoneNumberCommand {
	return &LinkPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(objectID, tenant, userId),
		PhoneNumberId: phoneNumberId,
		Primary:       primary,
		Label:         label,
	}
}

type LinkEmailCommand struct {
	eventstore.BaseCommand
	EmailId string `json:"emailId" validate:"required"`
	Primary bool
	Label   string
}

func NewLinkEmailCommand(objectID, tenant, userId, emailId, label string, primary bool) *LinkEmailCommand {
	return &LinkEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, userId),
		EmailId:     emailId,
		Primary:     primary,
		Label:       label,
	}
}

type LinkLocationCommand struct {
	eventstore.BaseCommand
	LocationId string `json:"locationId" validate:"required"`
}

func NewLinkLocationCommand(organizationId, tenant, userId, locationId string) *LinkLocationCommand {
	return &LinkLocationCommand{
		BaseCommand: eventstore.NewBaseCommand(organizationId, tenant, userId),
		LocationId:  locationId,
	}
}

type UpdateRenewalLikelihoodCommand struct {
	eventstore.BaseCommand
	Fields models.RenewalLikelihoodFields
}

func NewUpdateRenewalLikelihoodCommand(tenant, orgId, userId string, fields models.RenewalLikelihoodFields) *UpdateRenewalLikelihoodCommand {
	return &UpdateRenewalLikelihoodCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId),
		Fields:      fields,
	}
}

type RequestNextCycleDateCommand struct {
	eventstore.BaseCommand
}

func NewRequestNextCycleDateCommand(tenant, orgId, userId string) *RequestNextCycleDateCommand {
	return &RequestNextCycleDateCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId),
	}
}

type RequestRenewalForecastCommand struct {
	eventstore.BaseCommand
}

func NewRequestRenewalForecastCommand(tenant, orgId, userId string) *RequestRenewalForecastCommand {
	return &RequestRenewalForecastCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId),
	}
}

type UpdateRenewalForecastCommand struct {
	eventstore.BaseCommand
	Fields            models.RenewalForecastFields
	RenewalLikelihood models.RenewalLikelihoodProbability
}

func NewUpdateRenewalForecastCommand(tenant, orgId, userId string, fields models.RenewalForecastFields, renewalLikelihood models.RenewalLikelihoodProbability) *UpdateRenewalForecastCommand {
	return &UpdateRenewalForecastCommand{
		BaseCommand:       eventstore.NewBaseCommand(orgId, tenant, userId),
		Fields:            fields,
		RenewalLikelihood: renewalLikelihood,
	}
}

type UpdateBillingDetailsCommand struct {
	eventstore.BaseCommand
	Fields models.BillingDetailsFields
}

func NewUpdateBillingDetailsCommand(tenant, orgId, userId string, fields models.BillingDetailsFields) *UpdateBillingDetailsCommand {
	return &UpdateBillingDetailsCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId),
		Fields:      fields,
	}
}

type LinkDomainCommand struct {
	eventstore.BaseCommand
	Domain string
}

func NewLinkDomainCommand(objectID, tenant, domain, userId string) *LinkDomainCommand {
	return &LinkDomainCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, userId),
		Domain:      domain,
	}
}

type AddSocialCommand struct {
	eventstore.BaseCommand
	SocialId       string
	SocialPlatform string
	SocialUrl      string `json:"socialUrl" validate:"required"`
	Source         cmnmod.Source
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func NewAddSocialCommand(objectID, tenant, socialId, socialPlatform, socialUrl, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *AddSocialCommand {
	return &AddSocialCommand{
		BaseCommand:    eventstore.NewBaseCommand(objectID, tenant, ""),
		SocialId:       socialId,
		SocialPlatform: socialPlatform,
		SocialUrl:      socialUrl,
		Source: cmnmod.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type HideOrganizationCommand struct {
	eventstore.BaseCommand
}

func NewHideOrganizationCommand(tenant, orgId, userId string) *HideOrganizationCommand {
	return &HideOrganizationCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId),
	}
}

type ShowOrganizationCommand struct {
	eventstore.BaseCommand
}

func NewShowOrganizationCommand(tenant, orgId, userId string) *ShowOrganizationCommand {
	return &ShowOrganizationCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId),
	}
}

type RefreshLastTouchpointCommand struct {
	eventstore.BaseCommand
	AppSource string
}

func NewRefreshLastTouchpointCommand(tenant, orgId, userId, appSource string) *RefreshLastTouchpointCommand {
	return &RefreshLastTouchpointCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId),
		AppSource:   appSource,
	}
}

type UpsertCustomFieldCommand struct {
	eventstore.BaseCommand
	Source          cmnmod.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	CustomFieldData models.CustomField
}

func NewUpsertCustomFieldCommand(organizationId, tenant, source, sourceOfTruth, appSource, userId string,
	createdAt, updatedAt *time.Time, customField models.CustomField) *UpsertCustomFieldCommand {
	return &UpsertCustomFieldCommand{
		BaseCommand: eventstore.NewBaseCommand(organizationId, tenant, userId),
		Source: cmnmod.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		CustomFieldData: customField,
	}
}

type AddParentCommand struct {
	eventstore.BaseCommand
	ParentOrganizationId string `json:"parentOrganizationId" validate:"required,nefield=ObjectID"`
	Type                 string
	AppSource            string
}

func NewAddParentCommand(organizationId, tenant, userId, parentOrganizationId, relType, appSource string) *AddParentCommand {
	return &AddParentCommand{
		BaseCommand:          eventstore.NewBaseCommand(organizationId, tenant, userId),
		ParentOrganizationId: parentOrganizationId,
		Type:                 relType,
		AppSource:            appSource,
	}
}

type RemoveParentCommand struct {
	eventstore.BaseCommand
	ParentOrganizationId string `json:"parentOrganizationId" validate:"required,nefield=ObjectID"`
	AppSource            string
}

func NewRemoveParentCommand(organizationId, tenant, userId, parentOrganizationId, appSource string) *RemoveParentCommand {
	return &RemoveParentCommand{
		BaseCommand:          eventstore.NewBaseCommand(organizationId, tenant, userId),
		ParentOrganizationId: parentOrganizationId,
		AppSource:            appSource,
	}
}
