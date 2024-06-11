package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	service "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/client"
	authServices "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"gorm.io/gorm"
)

type Services struct {
	Cache       *caches.Cache
	GrpcClients *grpc_client.Clients

	CommonServices *commonService.Services
	AuthServices   *authServices.Services

	CustomerOsClient    CustomerOsClient
	CustomerOSApiClient service.CustomerOSApiClient
	RegistrationService RegistrationService
	TenantDataInjector  TenantDataInjector
}

func InitServices(cfg *config.Config, db *gorm.DB, driver *neo4j.DriverWithContext, grpcClients *grpc_client.Clients, cache *caches.Cache) *Services {
	services := Services{
		Cache:               cache,
		GrpcClients:         grpcClients,
		CustomerOsClient:    NewCustomerOsClient(cfg, driver),
		CustomerOSApiClient: service.NewCustomerOsClient(cfg.CustomerOS.CustomerOsAPI, cfg.CustomerOS.CustomerOsAPIKey),
	}

	services.CommonServices = commonService.InitServices(db, driver, cfg.Neo4j.Database, grpcClients)
	services.AuthServices = authServices.InitServices(nil, services.CommonServices, db)
	services.RegistrationService = NewRegistrationService(&services)
	services.TenantDataInjector = NewTenantDataInjector(&services)

	return &services
}
