package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type UserCommandHandlers struct {
	UpsertUser             UpsertUserCommandHandler
	AddPlayerInfo          AddPlayerInfoCommandHandler
	AddRole                AddRoleCommandHandler
	RemoveRole             RemoveRoleCommandHandler
	LinkJobRoleCommand     LinkJobRoleCommandHandler
	LinkPhoneNumberCommand LinkPhoneNumberCommandHandler
	LinkEmailCommand       LinkEmailCommandHandler
}

func NewUserCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *UserCommandHandlers {
	return &UserCommandHandlers{
		UpsertUser:             NewUpsertUserCommandHandler(log, cfg, es),
		AddPlayerInfo:          NewAddPlayerInfoCommandHandler(log, cfg, es),
		LinkPhoneNumberCommand: NewLinkPhoneNumberCommandHandler(log, cfg, es),
		LinkEmailCommand:       NewLinkEmailCommandHandler(log, cfg, es),
		LinkJobRoleCommand:     NewLinkJobRoleCommandHandler(log, cfg, es),
		AddRole:                NewAddRoleCommandHandler(log, cfg, es),
		RemoveRole:             NewRemoveRoleCommandHandler(log, cfg, es),
	}
}
