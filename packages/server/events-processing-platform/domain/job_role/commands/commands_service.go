package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type CommandHandlers struct {
	CreateJobRoleCommand CreateJobRoleCommandHander
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		CreateJobRoleCommand: NewCreateJobRoleCommandHandler(log, cfg, es),
	}
}
