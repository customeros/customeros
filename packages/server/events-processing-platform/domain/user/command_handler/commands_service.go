package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommandHandlers struct {
	UpsertUser             UpsertUserCommandHandler
	AddPlayerInfo          AddPlayerInfoCommandHandler
	AddRole                AddRoleCommandHandler
	RemoveRole             RemoveRoleCommandHandler
	LinkJobRoleCommand     LinkJobRoleCommandHandler
	LinkPhoneNumberCommand LinkPhoneNumberCommandHandler
	LinkEmailCommand       LinkEmailCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertUser:             NewUpsertUserCommandHandler(log, cfg, es),
		AddPlayerInfo:          NewAddPlayerInfoCommandHandler(log, cfg, es),
		LinkPhoneNumberCommand: NewLinkPhoneNumberCommandHandler(log, cfg, es),
		LinkEmailCommand:       NewLinkEmailCommandHandler(log, cfg, es),
		LinkJobRoleCommand:     NewLinkJobRoleCommandHandler(log, cfg, es),
		AddRole:                NewAddRoleCommandHandler(log, cfg, es),
		RemoveRole:             NewRemoveRoleCommandHandler(log, cfg, es),
	}
}
