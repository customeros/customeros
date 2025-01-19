package service

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	neo4jrepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-platform-admin-api/config"
)

type Services struct {
	cfg *config.Config

	GrpcClients *grpc_client.Clients

	CommonServices *commonService.CommonServices
}

func InitServices(
	postgresRepositories *postgres_repository.Repositories,
	neo4jRepositories *neo4jrepo.Repositories,
	cfg *config.Config,
	grpcClients *grpc_client.Clients,
	appLogger logger.Logger,
) *Services {
	services := Services{
		cfg:         cfg,
		GrpcClients: grpcClients,
	}

	services.CommonServices = commonService.InitCommonServices(
		appLogger,
		neo4jRepositories,
		postgresRepositories,
		cfg.Common,
		grpcClients,
		&commonService.InitOptions{},
	)

	return &services
}
