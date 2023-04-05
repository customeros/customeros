package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type ContactCommands struct {
	UpsertContact UpsertContactCommandHandler
	//CreateContact CreateContactCommandHandler
}

func NewContactCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *ContactCommands {
	return &ContactCommands{
		UpsertContact: NewUpsertContactHandler(log, cfg, es),
		//CreateContact: NewCreateContactHandler(log, cfg, es),
	}
}
