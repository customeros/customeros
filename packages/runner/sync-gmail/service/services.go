package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	"gorm.io/gorm"
)

type Services struct {
	cfg *config.Config

	TenantService    TenantService
	UserService      UserService
	AnthropicService AnthropicService
	OpenAiService    OpenAiService
	EmailService     EmailService
}

func InitServices(driver *neo4j.DriverWithContext, gormDb *gorm.DB, cfg *config.Config) *Services {
	repositories := repository.InitRepos(driver, gormDb)

	services := new(Services)
	services.cfg = cfg

	services.TenantService = NewTenantService(repositories)
	services.UserService = NewUserService(repositories)
	services.AnthropicService = NewAnthropicService(cfg)
	services.OpenAiService = NewOpenAiService(cfg)
	services.EmailService = NewEmailService(cfg, repositories, services)

	return services
}
