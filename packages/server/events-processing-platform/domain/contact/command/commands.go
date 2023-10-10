package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertContactCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      models.ContactDataFields
	Source          cmnmod.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertContactCommand(contactId, tenant, userId string, source cmnmod.Source, externalSystem cmnmod.ExternalSystem,
	dataFields models.ContactDataFields, createdAt, updatedAt *time.Time, create bool) *UpsertContactCommand {
	return &UpsertContactCommand{
		BaseCommand:     eventstore.NewBaseCommand(contactId, tenant, userId),
		DataFields:      dataFields,
		Source:          source,
		ExternalSystem:  externalSystem,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		IsCreateCommand: create,
	}
}

type LinkPhoneNumberCommand struct {
	eventstore.BaseCommand
	PhoneNumberId string `json:"phoneNumberId" validate:"required"`
	Primary       bool
	Label         string
}

func NewLinkPhoneNumberCommand(contactId, tenant, userId, phoneNumberId, label string, primary bool) *LinkPhoneNumberCommand {
	return &LinkPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(contactId, tenant, userId),
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

func NewLinkEmailCommand(contactId, tenant, userId, emailId, label string, primary bool) *LinkEmailCommand {
	return &LinkEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(contactId, tenant, userId),
		EmailId:     emailId,
		Primary:     primary,
		Label:       label,
	}
}

type LinkLocationCommand struct {
	eventstore.BaseCommand
	LocationId string `json:"locationId" validate:"required"`
}

func NewLinkLocationCommand(contactId, tenant, userId, locationId string) *LinkLocationCommand {
	return &LinkLocationCommand{
		BaseCommand: eventstore.NewBaseCommand(contactId, tenant, userId),
		LocationId:  locationId,
	}
}
