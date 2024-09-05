package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonconfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/config"
	"gorm.io/gorm"
)

type Services struct {
	Logger               logger.Logger
	CommonServices       *commonservice.Services
	ScrapeInService      ScrapinService
	BettercontactService BettercontactService
	BrandfetchService    BrandfetchService
}

func InitServices(config *config.Config, gormDb *gorm.DB, driver *neo4j.DriverWithContext, logger logger.Logger) *Services {
	services := &Services{
		CommonServices: commonservice.InitServices(&commonconfig.GlobalConfig{}, gormDb, driver, config.Neo4j.Database, nil),
	}
	services.Logger = logger
	services.ScrapeInService = NewScrapeInService(config, services, logger)
	services.BettercontactService = NewBettercontactService(config, services, logger)
	services.BrandfetchService = NewBrandfetchService(config, services, logger)

	return services
}
