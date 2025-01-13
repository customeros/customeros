package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"

	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
)

type Services struct {
	Cache    *caches.Cache
	Config   *config.Config
	Postgres *repository.Repositories
	Neo4j    *neoRepo.Repositories
	Log      logger.Logger

	AddressValidationService     AddressValidationService
	PhoneNumberValidationService PhoneNumberValidationService
	EmailValidationService       EmailValidationService
	IpIntelligenceService        IpIntelligenceService
}

func InitServices(cache *caches.Cache, config *config.Config, postgresDB *commonConfig.PostgresDB, driver *neo4j.DriverWithContext, log logger.Logger) *Services {
	services := &Services{}

	services.Cache = cache
	services.Log = log
	services.Config = config
	services.Postgres = repository.InitRepositories(postgresDB)
	services.Neo4j = neoRepo.InitNeo4jRepositories(driver, config.Neo4j.Database)

	services.AddressValidationService = NewAddressValidationService(config, services)
	services.PhoneNumberValidationService = NewPhoneNumberValidationService(services)
	services.EmailValidationService = NewEmailValidationService(config, services, log)
	services.IpIntelligenceService = NewIpIntelligenceService(config, services, log)

	return services
}
