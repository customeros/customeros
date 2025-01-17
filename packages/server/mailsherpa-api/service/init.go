package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/enrichment"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/verify"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
)

type Services struct {
	Cache                *caches.Cache
	PostgresRepositories *repository.Repositories
	EnrichmentService    interfaces.EnrichmentService
	VerifyService        interfaces.VerifyService
}

func InitServices(
	log logger.Logger,
	postgres *repository.Repositories,
	config *config.CommonConfig,
) *Services {

	cache := caches.NewCommonCache()

	services := Services{
		Cache:                cache,
		PostgresRepositories: postgres,
		EnrichmentService: enrichment.NewEnrichmentService(
			log,
			config.ExternalServices,
			postgres,
		),
	}

	services.VerifyService = verify.NewVerifyService(
		log,
		postgres,
		config,
		services.EnrichmentService,
	)

	return &services
}
