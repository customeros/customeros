package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"gorm.io/gorm"
)

type Services struct {
	CommonServices *commonService.Services

	AddressValidationService     AddressValidationService
	PhoneNumberValidationService PhoneNumberValidationService
	EmailValidationService       EmailValidationService
	IpIntelligenceService        IpIntelligenceService
}

func InitServices(config *config.Config, gormDb *gorm.DB, driver *neo4j.DriverWithContext, log logger.Logger) *Services {
	services := &Services{
		CommonServices: commonService.InitServices(&commonConfig.GlobalConfig{}, gormDb, driver, config.Neo4j.Database, nil),
	}

	services.AddressValidationService = NewAddressValidationService(config, services)
	services.PhoneNumberValidationService = NewPhoneNumberValidationService(services)
	services.EmailValidationService = NewEmailValidationService(config, services, log)
	services.IpIntelligenceService = NewIpIntelligenceService(config, services, log)

	return services
}
