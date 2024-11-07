package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	service "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/client"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"gorm.io/gorm"
)

type Services struct {
	Cache       *caches.Cache
	GrpcClients *grpc_client.Clients

	CommonServices *commonService.Services

	CustomerOsClient    CustomerOsClient
	CustomerOSApiClient service.CustomerOSApiClient
	RegistrationService RegistrationService
}

func InitServices(cfg *config.Config, db *gorm.DB, driver *neo4j.DriverWithContext, grpcClients *grpc_client.Clients, cache *caches.Cache, appLogger logger.Logger) *Services {
	services := Services{
		Cache:               cache,
		GrpcClients:         grpcClients,
		CustomerOsClient:    NewCustomerOsClient(cfg, driver),
		CustomerOSApiClient: service.NewCustomerOsClient(cfg.CustomerOS.CustomerOsAPI, cfg.CustomerOS.CustomerOsAPIKey),
	}

	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{
		GoogleOAuthConfig: &cfg.GoogleOAuth,
		RabbitMQConfig:    &cfg.RabbitMQConfig,
		ExternalServices: commonConfig.ExternalServices{
			OpenSRSConfig: cfg.OpenSRSConfig,
		},
	}, db, driver, cfg.Neo4j.Database, grpcClients, appLogger)
	services.RegistrationService = NewRegistrationService(&services)

	return &services
}
