package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	"gorm.io/gorm"
)

type Services struct {
	CommonAuthRepositories *repository.Repositories
	OAuthTokenService      OAuthTokenService
	GoogleService          GoogleService
}

func InitServices(cfg *config.Config, db *gorm.DB) *Services {
	repositories := repository.InitRepositories(db)

	services := &Services{
		CommonAuthRepositories: repositories,
		OAuthTokenService:      NewOAuthTokenService(repositories),
	}

	services.GoogleService = NewGoogleService(cfg, repositories, services)

	return services
}
