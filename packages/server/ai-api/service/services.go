package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/config"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
)

type Services struct {
	CommonServices *commonService.Services

	AnthropicService AnthropicService
	OpenAiService    OpenAiService
}

func InitServices(cfg *config.Config, postgresDB *commonConfig.PostgresDB, appLogger logger.Logger) *Services {
	services := &Services{
		CommonServices: commonService.InitServices(&commonConfig.GlobalConfig{}, postgresDB, nil, "", nil, appLogger),
	}

	services.OpenAiService = NewOpenAiService(cfg)
	services.AnthropicService = NewAnthropicService(cfg)

	return services
}
