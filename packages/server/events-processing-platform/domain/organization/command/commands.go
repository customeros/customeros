package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

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

type LinkDomainCommand struct {
	eventstore.BaseCommand
	Domain    string
	AppSource string
}

func NewLinkDomainCommand(objectID, tenant, domain, loggedInUserId, appSource string) *LinkDomainCommand {
	return &LinkDomainCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, loggedInUserId),
		Domain:      domain,
		AppSource:   appSource,
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
	CustomFieldData model.CustomField
}

func NewUpsertCustomFieldCommand(organizationId, tenant, source, sourceOfTruth, appSource, userId string,
	createdAt, updatedAt *time.Time, customField model.CustomField) *UpsertCustomFieldCommand {
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
