package container

import (
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
)

type Container struct {
	Cfg                           *config.Config
	Log                           logger.Logger
	Repositories                  *repository.Repositories
	CommonServices                *commonService.CommonServices
	EventProcessingServicesClient *grpc_client.Clients
}
