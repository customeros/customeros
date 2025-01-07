package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/config"
)

type Services struct {
	cfg *config.Config

	GrpcClients *grpc_client.Clients

	CommonServices *commonService.Services
}

func InitServices(
	driver *neo4j.DriverWithContext,
	postgresDB *commonConfig.PostgresDB,
	cfg *config.Config,
	grpcClients *grpc_client.Clients,
	appLogger logger.Logger) *Services {

	services := Services{
		cfg:         cfg,
		GrpcClients: grpcClients,
	}

	services.CommonServices = commonservice.InitServices(&commonConfig.GlobalConfig{}, postgresDB, driver, cfg.Neo4j.Database, grpcClients, appLogger)

	return &services
}
