package service

import (
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/openai-api/config"
)

type Services struct {
	CommonServices *commonService.Services

	OpenAiService OpenAiService
}

func InitServices(cfg *config.Config, db *config.StorageDB) *Services {
	services := &Services{
		CommonServices: commonService.InitServices(db.GormDB, nil),
	}

	services.OpenAiService = NewOpenAiService(cfg)

	return services
}
