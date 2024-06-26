package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
)

type Services struct {
	CommonServices *commonService.Services

	AddressValidationService     AddressValidationService
	PhoneNumberValidationService PhoneNumberValidationService
	EmailValidationService       EmailValidationService
	EmailFinderService           EmailFinderService
}

func InitServices(config *config.Config, db *config.StorageDB, driver *neo4j.DriverWithContext, log logger.Logger) *Services {
	services := &Services{
		CommonServices: commonService.InitServices(&commonConfig.GlobalConfig{}, db.GormDB, driver, config.Neo4j.Database, nil),
	}

	services.AddressValidationService = NewAddressValidationService(config, services)
	services.PhoneNumberValidationService = NewPhoneNumberValidationService(services)
	services.EmailValidationService = NewEmailValidationService(config, services, log)
	services.EmailFinderService = NewEmailFinderService(config, services, log)

	return services
}
