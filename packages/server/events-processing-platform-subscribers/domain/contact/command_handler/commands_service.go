package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommandHandlers struct {
	Upsert           UpsertContactCommandHandler
	LinkPhoneNumber  LinkPhoneNumberCommandHandler
	LinkEmail        LinkEmailCommandHandler
	LinkLocation     LinkLocationCommandHandler
	LinkOrganization LinkOrganizationCommandHandler
}

func NewCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		Upsert:           NewUpsertContactCommandHandler(log, es),
		LinkPhoneNumber:  NewLinkPhoneNumberCommandHandler(log, es),
		LinkEmail:        NewLinkEmailCommandHandler(log, es),
		LinkLocation:     NewLinkLocationCommandHandler(log, es),
		LinkOrganization: NewLinkOrganizationCommandHandler(log, es),
	}
}
