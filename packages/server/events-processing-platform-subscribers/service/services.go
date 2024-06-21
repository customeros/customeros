package service

import (
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/notifications"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
)

type Services struct {
	FileStoreApiService fsc.FileStoreApiService
	CommonServices      *commonService.Services
	PostmarkProvider    *notifications.PostmarkProvider
}

func InitServices(cfg *config.Config, repositories *repository.Repositories, log logger.Logger, grpcClients *grpc_client.Clients) *Services {
	services := Services{}

	services.FileStoreApiService = fsc.NewFileStoreApiService(&cfg.Services.FileStoreApiConfig)
	services.CommonServices = commonService.InitServices(repositories.Drivers.GormDb, repositories.Drivers.Neo4jDriver, cfg.Neo4j.Database, grpcClients)
	services.PostmarkProvider = notifications.NewPostmarkProvider(log, repositories)

	return &services
}
