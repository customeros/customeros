package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type ContactCommands struct {
	UpsertContact          UpsertContactCommandHandler
	LinkPhoneNumberCommand LinkPhoneNumberCommandHandler
	LinkEmailCommand       LinkEmailCommandHandler
	LinkLocationCommand    LinkLocationCommandHandler
}

func NewContactCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *ContactCommands {
	return &ContactCommands{
		UpsertContact:          NewUpsertContactCommandHandler(log, es),
		LinkPhoneNumberCommand: NewLinkPhoneNumberCommandHandler(log, es),
		LinkEmailCommand:       NewLinkEmailCommandHandler(log, es),
		LinkLocationCommand:    NewLinkLocationCommandHandler(log, es),
	}
}
