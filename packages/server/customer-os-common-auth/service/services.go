package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	"gorm.io/gorm"
)

type Services struct {
	CommonAuthRepositories *repository.Repositories
	OAuthTokenService      OAuthTokenService
}

func InitServices(db *gorm.DB) *Services {
	repositories := repository.InitRepositories(db)

	services := &Services{
		CommonAuthRepositories: repositories,
		OAuthTokenService:      NewOAuthTokenService(repositories),
	}

	return services
}
