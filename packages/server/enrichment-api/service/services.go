package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"

	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/config"
)

type Services struct {
	Logger               logger.Logger
	BettercontactService BettercontactService
	BrandfetchService    BrandfetchService
	ScrapeInService      ScrapinService
	IPIdentityService    IPIdentityService
}

func InitServices(config *config.Config, postgresDB *commonConfig.PostgresDB, driver *neo4j.DriverWithContext, logger logger.Logger) *Services {
	services := &Services{
		CommonServices: commonservice.InitServices(&commonConfig.GlobalConfig{}, postgresDB, driver, config.Neo4j.Database, nil, logger),
	}
	services.Logger = logger
	services.BettercontactService = NewBettercontactService(config, services, logger)
	services.BrandfetchService = NewBrandfetchService(config, services, logger)
	services.ScrapeInService = NewScrapInService(config, services, logger)
	services.IPIdentityService = NewIPIdentityService(config, services)

	return services
}
