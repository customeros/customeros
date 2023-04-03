package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type ContactCommandsService struct {
	Commands *commands.ContactCommands
}

func NewContactCommandsService(
	log logger.Logger,
	cfg *config.Config,
	es eventstore.AggregateStore,
) *ContactCommandsService {

	createContactHandler := commands.NewCreateContactHandler(log, cfg, es)
	updateContactCmdHandler := commands.NewUpdateContactCmdHandler(log, cfg, es)
	deleteContactCommandHandler := commands.NewDeleteContactCommandHandler(log, cfg, es)

	contactCommands := commands.NewContactCommands(
		createContactHandler,
		updateContactCmdHandler,
		deleteContactCommandHandler,
	)

	return &ContactCommandsService{Commands: contactCommands}
}
