package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertContactCommand struct {
	eventstore.BaseCommand
	Tenant    string
	FirstName string
	LastName  string
	Name      string
	Prefix    string
	Source    models.Source
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func NewUpsertContactCommand(aggregateID, tenant, firstName, lastName, name, prefix, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *UpsertContactCommand {
	return &UpsertContactCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		FirstName:   firstName,
		LastName:    lastName,
		Name:        name,
		Prefix:      prefix,
		Source: models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
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
