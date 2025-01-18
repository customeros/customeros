package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/caches"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
)

type Services struct {
	cfg          *config.Config
	Repositories *repository.Repositories

	CommonServices *commonService.CommonServices

	grpcClients *grpc_client.Clients
	Cache       *caches.Cache

	SyncService    SyncService
	MeetingService MeetingService
}

func InitServices(cfg *config.Config, driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB, grpcClients *grpc_client.Clients, cache *caches.Cache, log logger.Logger) *Services {
	repositories := repository.InitRepos(cfg, driver, postgresDB)

	services := new(Services)
	services.cfg = cfg
	services.Cache = cache
	services.grpcClients = grpcClients
	services.Repositories = repositories

	services.SyncService = NewSyncService(cfg, repositories, services)
	services.MeetingService = NewMeetingService(cfg, repositories, services)

	neo4jRepositories := neo4jrepository.InitNeo4jRepositories(driver, cfg.Neo4jDb.Database)
	postgresRepositories := postgresRepository.InitRepositories(postgresDB)

	services.CommonServices = commonService.InitCommonServices(
		log,
		neo4jRepositories,
		postgresRepositories,
		&cfg.CommonConfig,
		grpc_client.InitClients(nil),
	)
	return services
}
