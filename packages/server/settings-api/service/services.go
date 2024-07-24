package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/config"
	"gorm.io/gorm"
)

type Services struct {
	CommonServices *commonService.Services

	TenantSettingsService       TenantSettingsService
	PersonalIntegrationsService PersonalIntegrationsService
	OAuthUserSettingsService    OAuthUserSettingsService
	SlackSettingsService        SlackSettingsService
}

func InitServices(cfg *config.Config, db *gorm.DB, driver *neo4j.DriverWithContext, logger logger.Logger) *Services {
	services := &Services{
		CommonServices: commonService.InitServices(&commconf.GlobalConfig{}, db, driver, cfg.Neo4j.Database, nil),
	}

	services.TenantSettingsService = NewTenantSettingsService(services, logger, cfg)
	services.OAuthUserSettingsService = NewUserSettingsService(services, logger)
	services.SlackSettingsService = NewSlackSettingsService(services, logger)
	services.PersonalIntegrationsService = NewPersonalIntegrationsService(services, logger)

	return services
}
