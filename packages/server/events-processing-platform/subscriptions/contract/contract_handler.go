package contract

import (
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

type contractHandler struct {
	repositories        *repository.Repositories
	opportunityCommands *opportunitycmdhandler.CommandHandlers
	log                 logger.Logger
}
