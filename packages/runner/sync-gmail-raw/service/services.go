package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"gorm.io/gorm"
)

type Services struct {
	cfg *config.Config

	Repositories *repository.Repositories

	CommonServices *commonService.Services

	UserService    UserService
	TenantService  TenantService
	EmailService   EmailService
	MeetingService MeetingService
}

func InitServices(driver *neo4j.DriverWithContext, gormDb *gorm.DB, cfg *config.Config) *Services {
	services := new(Services)
	services.cfg = cfg
	services.Repositories = repository.InitRepos(driver, gormDb)

	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{GoogleOAuthConfig: &cfg.AuthConfig}, gormDb, driver, cfg.Neo4jDb.Database, nil)

	services.TenantService = NewTenantService(services.Repositories)
	services.UserService = NewUserService(services.Repositories)
	services.EmailService = NewEmailService(cfg, services.Repositories, services)
	services.MeetingService = NewMeetingService(cfg, services.Repositories, services)

	return services
}
