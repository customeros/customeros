package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"gorm.io/gorm"
)

type Services struct {
	CommonServices *commonService.Services

	CommonAuthRepositories *repository.Repositories
	OAuthTokenService      OAuthTokenService
	GoogleService          GoogleService
}

func InitServices(cfg *config.Config, commonServices *commonService.Services, db *gorm.DB) *Services {
	repositories := repository.InitRepositories(db)

	services := &Services{
		CommonServices:         commonServices,
		CommonAuthRepositories: repositories,
		OAuthTokenService:      NewOAuthTokenService(repositories),
	}

	services.GoogleService = NewGoogleService(cfg, commonServices.CommonRepositories, repositories, services)

	return services
}
