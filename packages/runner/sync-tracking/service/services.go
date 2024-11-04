package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/config"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"gorm.io/gorm"
)

type Services struct {
	Logger logger.Logger

	GrpcClient *grpc_client.Clients

	CommonServices *commonService.Services

	EnrichDetailsTrackingService TrackingService
}

func InitServices(cfg *config.Config, driver *neo4j.DriverWithContext, gormDb *gorm.DB, grpcClient *grpc_client.Clients, logger logger.Logger) *Services {
	services := new(Services)

	services.GrpcClient = grpcClient

	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{
		RabbitMQConfig: &cfg.RabbitMQConfig,
	}, gormDb, driver, cfg.Neo4j.Database, grpcClient, logger)

	services.EnrichDetailsTrackingService = NewTrackingService(cfg, services)

	services.Logger = logger

	return services
}
