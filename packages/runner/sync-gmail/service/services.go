package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/caches"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	config2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"gorm.io/gorm"
)

type Services struct {
	cfg          *config.Config
	Repositories *repository.Repositories

	CommonServices *commonService.Services

	grpcClients *grpc_client.Clients
	Cache       *caches.Cache

	SyncService    SyncService
	EmailService   EmailService
	MeetingService MeetingService
}

func InitServices(cfg *config.Config, driver *neo4j.DriverWithContext, gormDb *gorm.DB, grpcClients *grpc_client.Clients, cache *caches.Cache, appLogger logger.Logger) *Services {
	repositories := repository.InitRepos(cfg, driver, gormDb)

	services := new(Services)
	services.cfg = cfg
	services.Cache = cache
	services.grpcClients = grpcClients
	services.Repositories = repositories

	services.SyncService = NewSyncService(cfg, repositories, services)
	services.EmailService = NewEmailService(cfg, repositories, services)
	services.MeetingService = NewMeetingService(cfg, repositories, services)

	services.CommonServices = commonService.InitServices(&config2.GlobalConfig{}, repositories.PostgresDriver, repositories.Neo4jDriver, "neo4j", grpcClients, appLogger)

	return services
}
