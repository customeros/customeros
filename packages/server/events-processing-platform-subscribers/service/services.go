package service

import (
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/notifications"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type Services struct {
	Repositories *repository.Repositories

	CommonServices *commonService.Services

	Es eventstore.AggregateStore

	EventBufferStoreService *eventbuffer.EventBufferStoreService
	FileStoreApiService     fsc.FileStoreApiService
	PostmarkProvider        *notifications.PostmarkProvider
}

func InitServices(cfg *config.Config, es eventstore.AggregateStore, repositories *repository.Repositories, log logger.Logger, grpcClients *grpc_client.Clients) *Services {
	services := Services{}

	services.Repositories = repositories
	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{}, repositories.Drivers.GormDb, repositories.Drivers.Neo4jDriver, cfg.Neo4j.Database, grpcClients)

	services.Es = es

	services.EventBufferStoreService = eventbuffer.NewEventBufferStoreService(repositories.PostgresRepositories.EventBufferRepository, log)
	services.FileStoreApiService = fsc.NewFileStoreApiService(&cfg.Services.FileStoreApiConfig)
	services.PostmarkProvider = notifications.NewPostmarkProvider(log, repositories)

	return &services
}
