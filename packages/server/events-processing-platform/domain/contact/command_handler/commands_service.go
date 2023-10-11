package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type ContactCommands struct {
	UpsertContact           UpsertContactCommandHandler
	LinkPhoneNumberCommand  LinkPhoneNumberCommandHandler
	LinkEmailCommand        LinkEmailCommandHandler
	LinkLocationCommand     LinkLocationCommandHandler
	LinkOrganizationCommand LinkOrganizationCommandHandler
}

func NewContactCommands(log logger.Logger, es eventstore.AggregateStore) *ContactCommands {
	return &ContactCommands{
		UpsertContact:           NewUpsertContactCommandHandler(log, es),
		LinkPhoneNumberCommand:  NewLinkPhoneNumberCommandHandler(log, es),
		LinkEmailCommand:        NewLinkEmailCommandHandler(log, es),
		LinkLocationCommand:     NewLinkLocationCommandHandler(log, es),
		LinkOrganizationCommand: NewLinkOrganizationCommandHandler(log, es),
	}
}
