package container

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/events_processing_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
)

type Container struct {
	Cfg                           *config.Config
	Log                           logger.Logger
	Repositories                  *repository.Repositories
	EventProcessingServicesClient *events_processing_client.Client
}
