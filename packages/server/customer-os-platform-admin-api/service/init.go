package service

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	neo4jrepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-platform-admin-api/config"
)

type Services struct {
	cfg            *config.Config
	CommonServices *commonService.CommonServices
}

func InitServices(
	postgresRepositories *postgres_repository.Repositories,
	neo4jRepositories *neo4jrepo.Repositories,
	cfg *config.Config,
	appLogger logger.Logger,
) *Services {
	services := Services{
		cfg: cfg,
	}

	services.CommonServices = commonService.InitCommonServices(
		appLogger,
		neo4jRepositories,
		postgresRepositories,
		cfg.Common,
		&commonService.InitOptions{},
	)

	return &services
}
