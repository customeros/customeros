package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/config"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
)

type Services struct {
	CommonServices *commonService.Services

	AnthropicService AnthropicService
	OpenAiService    OpenAiService
}

func InitServices(cfg *config.Config, db *config.StorageDB) *Services {
	services := &Services{
		CommonServices: commonService.InitServices(db.GormDB, nil, "", nil),
	}

	services.OpenAiService = NewOpenAiService(cfg)
	services.AnthropicService = NewAnthropicService(cfg)

	return services
}
