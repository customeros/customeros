package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"

	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
)

type Services struct {
	Cache       *caches.Cache
	GrpcClients *grpc_client.Clients

	CommonServices *commonService.Services
}

func InitServices(cfg *config.Config, driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB, grpcClients *grpc_client.Clients, cache *caches.Cache, appLogger logger.Logger) *Services {
	services := Services{
		Cache:       cache,
		GrpcClients: grpcClients,
	}

	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{
		GoogleOAuthConfig: &cfg.GoogleOAuth,
		RabbitMQConfig:    &cfg.RabbitMQ,
		ExternalServices: commonConfig.ExternalServices{
			OpenSRSConfig:  cfg.OpenSRS,
			PostmarkConfig: cfg.Postmark,
		},
		InternalServices: commonConfig.InternalServices{
			EnrichmentApiConfig: cfg.InternalServices.EnrichmentApi,
			ValidationApiConfig: cfg.InternalServices.ValidationApi,
		},
	}, postgresDB, driver, cfg.Neo4j.Database, grpcClients, appLogger)

	return &services
}
