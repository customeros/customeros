package command

import (
	"time"

	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
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

func NewAddSocialCommand(organizationId, tenant, loggedInUserId, socialId, socialPlatform, socialUrl string, sourceFields cmnmod.Source, createdAt, updatedAt *time.Time) *AddSocialCommand {
	return &AddSocialCommand{
		BaseCommand:    eventstore.NewBaseCommand(organizationId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
		SocialId:       socialId,
		SocialPlatform: socialPlatform,
		SocialUrl:      socialUrl,
		Source:         sourceFields,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
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

type OrganizationOwnerUpdateCommand struct {
	eventstore.BaseCommand
	UpdatedAt      time.Time `json:"updatedAt"`
	OwnerUserId    string    `json:"userId" validate:"required"` // who became owner
	OrganizationId string    `json:"organizationId" validate:"required"`
	ActorUserId    string    `json:"actorUserId"` // who set the owner
}

func NewOrganizationOwnerUpdateEvent(organizationId, tenant, userId, actorUserId, appSource string, updatedAt time.Time) *OrganizationOwnerUpdateCommand {
	return &OrganizationOwnerUpdateCommand{
		BaseCommand:    eventstore.NewBaseCommand(organizationId, tenant, userId),
		UpdatedAt:      updatedAt,
		OwnerUserId:    userId,
		OrganizationId: organizationId,
		ActorUserId:    actorUserId,
	}
}
