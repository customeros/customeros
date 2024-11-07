package container

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
)

type Container struct {
	Cfg                           *config.Config
	Log                           logger.Logger
	Repositories                  *repository.Repositories
	CommonServices                *commonService.Services
	EventProcessingServicesClient *grpc_client.Clients
	EventBufferStoreService       *eventbuffer.EventBufferStoreService
}
