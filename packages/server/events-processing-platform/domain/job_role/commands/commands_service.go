package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type JobRoleCommands struct {
	CreateJobRole CreateJobRoleCommandHander
}

func NewJobRoleCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *JobRoleCommands {
	return &JobRoleCommands{
		CreateJobRole: NewCreateJobRoleCommandHandler(log, cfg, es),
	}
}
