package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type JobRoleCommands struct {
	CreateJobRoleCommand CreateJobRoleCommandHander
}

func NewJobRoleCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *JobRoleCommands {
	return &JobRoleCommands{
		CreateJobRoleCommand: NewCreateJobRoleCommandHandler(log, cfg, es),
	}
}
