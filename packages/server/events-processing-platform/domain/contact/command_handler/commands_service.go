package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type ContactCommandHandlers struct {
	Upsert           UpsertContactCommandHandler
	LinkPhoneNumber  LinkPhoneNumberCommandHandler
	LinkEmail        LinkEmailCommandHandler
	LinkLocation     LinkLocationCommandHandler
	LinkOrganization LinkOrganizationCommandHandler
}

func NewContactCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *ContactCommandHandlers {
	return &ContactCommandHandlers{
		Upsert:           NewUpsertContactCommandHandler(log, es),
		LinkPhoneNumber:  NewLinkPhoneNumberCommandHandler(log, es),
		LinkEmail:        NewLinkEmailCommandHandler(log, es),
		LinkLocation:     NewLinkLocationCommandHandler(log, es),
		LinkOrganization: NewLinkOrganizationCommandHandler(log, es),
	}
}
