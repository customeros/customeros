package commands

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"

	"time"
)

type ContactCoreFields struct {
	FirstName   string
	LastName    string
	Prefix      string
	Description string
}

func NewContactCoreFields(firstName, lastName, prefix, description string) ContactCoreFields {
	return ContactCoreFields{
		FirstName:   firstName,
		LastName:    lastName,
		Prefix:      prefix,
		Description: description,
	}
}

type UpsertContactCommand struct {
	eventstore.BaseCommand
	CoreFields ContactCoreFields
	Source     common_models.Source
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func UpsertContactCommandToContactDto(command *UpsertContactCommand) *models.ContactDto {
	return &models.ContactDto{
		ID:        command.ObjectID,
		Tenant:    command.Tenant,
		FirstName: command.CoreFields.FirstName,
		LastName:  command.CoreFields.LastName,
		Prefix:    command.CoreFields.Prefix,
		Source:    command.Source,
		CreatedAt: command.CreatedAt,
		UpdatedAt: command.UpdatedAt,
	}
}

func NewUpsertContactCommand(objectID, tenant, source, sourceOfTruth, appSource string, coreFields ContactCoreFields, createdAt, updatedAt *time.Time) *UpsertContactCommand {
	return &UpsertContactCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant),
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

type CreateContactCommand struct {
	eventstore.BaseCommand
	ContactCoreFields
	common_models.Source
	CreatedAt *time.Time
}

// FIXME alexb re-implement all below

type UpdateContactCommand struct {
	eventstore.BaseCommand
	UUID      string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
}

func NewContactCreateCommand(objectID, tenant, firstName, lastName, prefix, description, source, sourceOfTruth, appSource string, createdAt *time.Time) *CreateContactCommand {
	return &CreateContactCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant),
		ContactCoreFields: ContactCoreFields{
			FirstName:   firstName,
			LastName:    lastName,
			Prefix:      prefix,
			Description: description,
		},
		Source: common_models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
	}
}
