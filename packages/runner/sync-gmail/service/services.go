package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	"gorm.io/gorm"
)

type Services struct {
	TenantService TenantService
	UserService   UserService
	EmailService  EmailService
}

func InitServices(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Services {
	repositories := repository.InitRepos(driver, gormDb)

	services := new(Services)

	services.TenantService = NewTenantService(repositories)
	services.UserService = NewUserService(repositories)
	services.EmailService = NewEmailService(repositories)

	return services
}
