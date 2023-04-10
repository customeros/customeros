package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type OrganizationCommands struct {
	UpsertOrganization     UpsertOrganizationCommandHandler
	LinkPhoneNumberCommand LinkPhoneNumberCommandHandler
	LinkEmailCommand       LinkEmailCommandHandler
}

func NewOrganizationCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *OrganizationCommands {
	return &OrganizationCommands{
		UpsertOrganization:     NewUpsertOrganizationCommandHandler(log, cfg, es),
		LinkPhoneNumberCommand: NewLinkPhoneNumberCommandHandler(log, cfg, es),
		LinkEmailCommand:       NewLinkEmailCommandHandler(log, cfg, es),
	}
}
