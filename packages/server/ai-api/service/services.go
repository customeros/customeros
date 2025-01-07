package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/config"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
)

type Services struct {
	CommonServices *commonService.Services
	AIService      AIService
}

func InitServices(cfg *config.Config, postgresDB *commonConfig.PostgresDB, appLogger logger.Logger) *Services {
	services := &Services{
		CommonServices: commonService.InitServices(&commonConfig.GlobalConfig{}, postgresDB, nil, "", nil, appLogger),
	}

	services.AIService = NewAIService(cfg, services)

	return services
}
