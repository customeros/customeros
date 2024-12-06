package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
)

type Services struct {
	cfg *config.Config

	Repositories *repository.Repositories

	CommonServices *commonService.Services

	EmailService   EmailService
	MeetingService MeetingService
}

func InitServices(driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB, cfg *config.Config, logger *logger.ExtendedLogger) *Services {
	services := new(Services)
	services.cfg = cfg
	services.Repositories = repository.InitRepos(driver, postgresDB)

	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{GoogleOAuthConfig: &cfg.GoogleOAuthConfig, AzureOAuthConfig: &cfg.AzureOAuthConfig}, postgresDB, driver, cfg.Neo4jDb.Database, nil, logger)

	services.EmailService = NewEmailService(cfg, services.Repositories, services)
	services.MeetingService = NewMeetingService(cfg, services.Repositories, services)

	return services
}
