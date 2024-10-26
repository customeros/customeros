package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
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
	Source          common.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	CustomFieldData model.CustomField
}

func NewUpsertCustomFieldCommand(organizationId, tenant, source, sourceOfTruth, appSource, userId string,
	createdAt, updatedAt *time.Time, customField model.CustomField) *UpsertCustomFieldCommand {
	return &UpsertCustomFieldCommand{
		BaseCommand: eventstore.NewBaseCommand(organizationId, tenant, userId),
		Source: common.Source{
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
