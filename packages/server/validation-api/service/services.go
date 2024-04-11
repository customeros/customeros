package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
)

type Services struct {
	CommonServices *commonService.Services

	AddressValidationService     AddressValidationService
	PhoneNumberValidationService PhoneNumberValidationService
	EmailValidationService       EmailValidationService
}

func InitServices(config *config.Config, db *config.StorageDB, driver *neo4j.DriverWithContext) *Services {
	services := &Services{
		CommonServices: commonService.InitServices(db.GormDB, driver, config.Neo4j.Database),
	}

	services.AddressValidationService = NewAddressValidationService(config, services)
	services.PhoneNumberValidationService = NewPhoneNumberValidationService(services)
	services.EmailValidationService = NewEmailValidationService(config, services)

	return services
}
