package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type UpsertContactCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      models.ContactDataFields
	Source          events.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	FieldsMask      []string
}

func NewUpsertContactCommand(contactId, tenant, userId string, source events.Source, externalSystem cmnmod.ExternalSystem,
	dataFields models.ContactDataFields, createdAt, updatedAt *time.Time, create bool, fieldsMask []string) *UpsertContactCommand {
	return &UpsertContactCommand{
		BaseCommand:     eventstore.NewBaseCommand(contactId, tenant, userId),
		DataFields:      dataFields,
		Source:          source,
		ExternalSystem:  externalSystem,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		IsCreateCommand: create,
		FieldsMask:      fieldsMask,
	}
}

type LinkPhoneNumberCommand struct {
	eventstore.BaseCommand
	PhoneNumberId string `json:"phoneNumberId" validate:"required"`
	Primary       bool
	Label         string
	AppSource     string
}

func NewLinkPhoneNumberCommand(contactId, tenant, userId, phoneNumberId, label, appSource string, primary bool) *LinkPhoneNumberCommand {
	return &LinkPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(contactId, tenant, userId),
		PhoneNumberId: phoneNumberId,
		Primary:       primary,
		Label:         label,
		AppSource:     appSource,
	}
}

type LinkEmailCommand struct {
	eventstore.BaseCommand
	EmailId   string `json:"emailId" validate:"required"`
	Primary   bool
	Label     string
	AppSource string
}

func NewLinkEmailCommand(contactId, tenant, userId, emailId, label, appSource string, primary bool) *LinkEmailCommand {
	return &LinkEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(contactId, tenant, userId),
		EmailId:     emailId,
		Primary:     primary,
		Label:       label,
		AppSource:   appSource,
	}
}

type LinkLocationCommand struct {
	eventstore.BaseCommand
	LocationId string `json:"locationId" validate:"required"`
	AppSource  string
}

func NewLinkLocationCommand(contactId, tenant, userId, locationId, appSource string) *LinkLocationCommand {
	return &LinkLocationCommand{
		BaseCommand: eventstore.NewBaseCommand(contactId, tenant, userId),
		LocationId:  locationId,
		AppSource:   appSource,
	}
}

type LinkOrganizationCommand struct {
	eventstore.BaseCommand
	OrganizationId string `json:"organizationId" validate:"required"`
	JobRoleFields  models.JobRole
	Source         events.Source
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func NewLinkOrganizationCommand(contactId, tenant, userId, organizationId string, source events.Source, jobRoleFields models.JobRole, createdAt, updatedAt *time.Time) *LinkOrganizationCommand {
	return &LinkOrganizationCommand{
		BaseCommand:    eventstore.NewBaseCommand(contactId, tenant, userId),
		OrganizationId: organizationId,
		JobRoleFields:  jobRoleFields,
		Source:         source,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
