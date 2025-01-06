package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"

	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/logger"
)

type Services struct {
	Logger         logger.Logger
	CommonServices *service.Services
	AIModelService AIModelService
}

func InitServices(config *config.Config, postgresDB *commonConfig.PostgresDB, driver *neo4j.DriverWithContext, logger logger.Logger) *Services {
	services := &Services{
		CommonServices: service.InitServices(&commonConfig.GlobalConfig{}, postgresDB, driver, config.Neo4j.Database, nil, logger),
		Logger:         logger,
	}

	services.AIModelService = NewAIModelService(config, services)

	return services
}
