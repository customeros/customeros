package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type Services struct {
	CommonServices *commonService.Services

	Es eventstore.AggregateStore

	EventBufferStoreService *eventbuffer.EventBufferStoreService
	FileStoreApiService     fsc.FileStoreApiService
}

func InitServices(cfg *config.Config, es eventstore.AggregateStore, log logger.Logger, grpcClients *grpc_client.Clients, postgresDB *commonConfig.PostgresDB, driver *neo4j.DriverWithContext) *Services {
	services := Services{}

	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{
		RabbitMQConfig: &cfg.RabbitMQConfig,
	}, postgresDB, driver, cfg.Neo4j.Database, grpcClients, log)

	services.Es = es

	services.EventBufferStoreService = eventbuffer.NewEventBufferStoreService(services.CommonServices.PostgresRepositories.EventBufferRepository, log)
	services.FileStoreApiService = fsc.NewFileStoreApiService(&cfg.Services.FileStoreApiConfig)

	return &services
}
