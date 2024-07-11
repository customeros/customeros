package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type CommandHandlers struct {
	UpsertUser         UpsertUserCommandHandler
	AddPlayerInfo      AddPlayerInfoCommandHandler
	AddRole            AddRoleCommandHandler
	RemoveRole         RemoveRoleCommandHandler
	LinkJobRoleCommand LinkJobRoleCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertUser:         NewUpsertUserCommandHandler(log, cfg, es),
		AddPlayerInfo:      NewAddPlayerInfoCommandHandler(log, cfg, es),
		LinkJobRoleCommand: NewLinkJobRoleCommandHandler(log, cfg, es),
		AddRole:            NewAddRoleCommandHandler(log, cfg, es),
		RemoveRole:         NewRemoveRoleCommandHandler(log, cfg, es),
	}
}
