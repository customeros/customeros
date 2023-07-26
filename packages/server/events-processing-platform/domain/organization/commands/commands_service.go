package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type OrganizationCommands struct {
	UpsertOrganization     UpsertOrganizationCommandHandler
	UpdateOrganization     UpdateOrganizationCommandHandler
	LinkPhoneNumberCommand LinkPhoneNumberCommandHandler
	LinkEmailCommand       LinkEmailCommandHandler
	LinkDomainCommand      LinkDomainCommandHandler
	AddSocialCommand       AddSocialCommandHandler
}

func NewOrganizationCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *OrganizationCommands {
	return &OrganizationCommands{
		UpsertOrganization:     NewUpsertOrganizationCommandHandler(log, cfg, es),
		UpdateOrganization:     NewUpdateOrganizationCommandHandler(log, cfg, es),
		LinkPhoneNumberCommand: NewLinkPhoneNumberCommandHandler(log, cfg, es),
		LinkEmailCommand:       NewLinkEmailCommandHandler(log, cfg, es),
		LinkDomainCommand:      NewLinkDomainCommandHandler(log, cfg, es),
		AddSocialCommand:       NewAddSocialCommandHandler(log, cfg, es),
	}
}
