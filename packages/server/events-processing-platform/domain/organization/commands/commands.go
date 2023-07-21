package commands

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
	CoreFields        models.OrganizationDataFields
	SourceOfTruth     string `json:"sourceOfTruth"`
	UpdatedAt         *time.Time
}

func NewUpdateOrganizationCommand(organizationId, tenant, sourceOfTruth string, dataFields models.OrganizationDataFields, updatedAt *time.Time, ignoreEmptyFields bool) *UpdateOrganizationCommand {
	return &UpdateOrganizationCommand{
		BaseCommand:       eventstore.NewBaseCommand(organizationId, tenant),
		IgnoreEmptyFields: ignoreEmptyFields,
		CoreFields:        dataFields,
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

type LinkDomainCommand struct {
	eventstore.BaseCommand
	Domain string
}

func NewLinkDomainCommand(objectID, tenant, domain string) *LinkDomainCommand {
	return &LinkDomainCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant),
		Domain:      domain,
	}
}

type AddSocialCommand struct {
	eventstore.BaseCommand
	SocialId       string
	SocialPlatform string
	SocialUrl      string
	Source         common_models.Source
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func NewAddSocialCommand(objectID, tenant, socialId, socialPlatform, socialUrl, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *AddSocialCommand {
	return &AddSocialCommand{
		BaseCommand:    eventstore.NewBaseCommand(objectID, tenant),
		SocialId:       socialId,
		SocialPlatform: socialPlatform,
		SocialUrl:      socialUrl,
		Source: common_models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
