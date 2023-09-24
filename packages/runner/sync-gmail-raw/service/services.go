package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	"gorm.io/gorm"
)

type Services struct {
	cfg *config.Config

	Repositories *repository.Repositories

	UserService   UserService
	TenantService TenantService

	GmailService GmailService

	EmailService   EmailService
	MeetingService MeetingService
}

func InitServices(driver *neo4j.DriverWithContext, gormDb *gorm.DB, cfg *config.Config) *Services {
	services := new(Services)
	services.cfg = cfg

	services.Repositories = repository.InitRepos(driver, gormDb)
	services.TenantService = NewTenantService(services.Repositories)
	services.UserService = NewUserService(services.Repositories)
	services.GmailService = NewGmailService(cfg, services.Repositories, services)
	services.EmailService = NewEmailService(cfg, services.Repositories, services)
	services.MeetingService = NewMeetingService(cfg, services.Repositories, services)

	return services
}
