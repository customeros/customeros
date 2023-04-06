package commands

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type ContactCoreFields struct {
	FirstName string
	LastName  string
	Name      string
	Prefix    string
}

type UpsertContactCommand struct {
	eventstore.BaseCommand
	Tenant     string
	CoreFields ContactCoreFields
	Source     commonModels.Source
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func UpsertContactCommandToContactDto(command *UpsertContactCommand) *models.ContactDto {
	return &models.ContactDto{
		ID:        command.AggregateID,
		Tenant:    command.Tenant,
		FirstName: command.CoreFields.FirstName,
		LastName:  command.CoreFields.LastName,
		Name:      command.CoreFields.Name,
		Prefix:    command.CoreFields.Prefix,
		Source:    command.Source,
		CreatedAt: command.CreatedAt,
		UpdatedAt: command.UpdatedAt,
	}
}

func NewUpsertContactCommand(aggregateID, tenant, source, sourceOfTruth, appSource string, coreFields ContactCoreFields, createdAt, updatedAt *time.Time) *UpsertContactCommand {
	return &UpsertContactCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		CoreFields:  coreFields,
		Source: commonModels.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type AddPhoneNumberCommand struct {
	eventstore.BaseCommand
	Tenant        string
	PhoneNumberId string
	Primary       bool
	Label         string
}

func NewAddPhoneNumberCommand(aggregateID, tenant, phoneNumberId, label string, primary bool) *AddPhoneNumberCommand {
	return &AddPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(aggregateID),
		Tenant:        tenant,
		PhoneNumberId: phoneNumberId,
		Primary:       primary,
		Label:         label,
	}
}

// FIXME alexb re-implement all below
type CreateContactCommand struct {
	eventstore.BaseCommand
	UUID      string `json:"uuid"  validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func NewCreateContactCommand(aggregateID string, uuid string, firstName string, lastName string) *CreateContactCommand {
	return &CreateContactCommand{BaseCommand: eventstore.NewBaseCommand(aggregateID), UUID: uuid, FirstName: firstName, LastName: lastName}
}

type UpdateContactCommand struct {
	eventstore.BaseCommand
	UUID      string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
}

func NewUpdateContactCommand(aggregateID string, uuid string, firstName string, lastName string) *UpdateContactCommand {
	return &UpdateContactCommand{BaseCommand: eventstore.NewBaseCommand(aggregateID), UUID: uuid, FirstName: firstName, LastName: lastName}
}

type DeleteContactCommand struct {
	eventstore.BaseCommand
	UUID string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
}

func NewDeleteContactCommand(aggregateID string, uuid string) *DeleteContactCommand {
	return &DeleteContactCommand{BaseCommand: eventstore.NewBaseCommand(aggregateID), UUID: uuid}
}
