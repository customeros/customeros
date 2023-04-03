package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type CreateContactCommand struct {
	eventstore.BaseCommand
	UUID      string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
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
