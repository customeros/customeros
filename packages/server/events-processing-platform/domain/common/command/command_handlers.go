package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	locationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	organizationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type CommandHandlers struct {
	Organization *organizationcmdhandler.CommandHandlers
	Location     *locationcmdhandler.CommandHandlers
}

func NewCommandHandlers(log logger.Logger,
	cfg *config.Config,
	aggregateStore eventstore.AggregateStore,
	ebs *eventbuffer.EventBufferStoreService,
) *CommandHandlers {

	return &CommandHandlers{
		Organization: organizationcmdhandler.NewCommandHandlers(log, cfg, aggregateStore, ebs),
		Location:     locationcmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
	}
}
