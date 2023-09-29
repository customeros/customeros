package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type UserCommands struct {
	UpsertUser             UpsertUserCommandHandler
	LinkJobRoleCommand     LinkJobRoleCommandHandler
	LinkPhoneNumberCommand LinkPhoneNumberCommandHandler
	LinkEmailCommand       LinkEmailCommandHandler
}

func NewUserCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *UserCommands {
	return &UserCommands{
		UpsertUser:             NewUpsertUserCommandHandler(log, cfg, es),
		LinkPhoneNumberCommand: NewLinkPhoneNumberCommandHandler(log, cfg, es),
		LinkEmailCommand:       NewLinkEmailCommandHandler(log, cfg, es),
		LinkJobRoleCommand:     NewLinkJobRoleCommandHandler(log, cfg, es),
	}
}
