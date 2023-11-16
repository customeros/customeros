package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonAuthService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/repository"
	"gorm.io/gorm"
)

type Services struct {
	cfg *config.Config

	GrpcClients *grpc_client.Clients

	Repositories *repository.Repositories

	CommonServices     *commonService.Services
	CommonAuthServices *commonAuthService.Services
}

func InitServices(
	driver *neo4j.DriverWithContext,
	gormDB *gorm.DB,
	cfg *config.Config,
	grpcClients *grpc_client.Clients) *Services {

	services := Services{
		cfg:                cfg,
		GrpcClients:        grpcClients,
		Repositories:       repository.InitRepos(driver, gormDB, cfg.Neo4j.Database),
		CommonServices:     commonservice.InitServices(gormDB, driver),
		CommonAuthServices: commonAuthService.InitServices(nil, gormDB),
	}

	return &services
}
