package commands

import (
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
)

type CreatePhoneNumberCommand struct {
	eventstore.BaseCommand
	Tenant      string `json:"tenant" validate:"required"`
	PhoneNumber string `json:"rawPhoneNumber" validate:"required"`
}

func NewCreatePhoneNumberCommand(aggregateID, tenant, rawPhoneNumber string) *CreatePhoneNumberCommand {
	return &CreatePhoneNumberCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		PhoneNumber: rawPhoneNumber,
	}
}
