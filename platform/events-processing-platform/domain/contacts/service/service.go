package service

import (
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/contacts/commands"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
)

type ContactService struct {
	Commands *commands.ContactCommands
}

func NewContactService(
	log logger.Logger,
	cfg *config.Config,
	es eventstore.AggregateStore,
) *ContactService {

	createContactHandler := commands.NewCreateContactHandler(log, cfg, es)
	updateContactCmdHandler := commands.NewUpdateContactCmdHandler(log, cfg, es)
	deleteContactCommandHandler := commands.NewDeleteContactCommandHandler(log, cfg, es)

	contactCommands := commands.NewContactCommands(
		createContactHandler,
		updateContactCmdHandler,
		deleteContactCommandHandler,
	)

	return &ContactService{Commands: contactCommands}
}
