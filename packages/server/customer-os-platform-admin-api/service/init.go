package service

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
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
		nil,
		postgresRepositories,
		&cfg.Common,
		grpcClients,
	)

	return &services
}
