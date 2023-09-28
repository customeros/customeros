package repository

import (
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"gorm.io/gorm"
	"log"
)

type Repositories struct {
	OAuthTokenRepository repository.OAuthTokenRepository
	ApiKeyRepository     repository.ApiKeyRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		OAuthTokenRepository: repository.NewOAuthTokenRepository(db),
		ApiKeyRepository:     repository.NewApiKeyRepository(db),
	}

	var err error

	err = db.AutoMigrate(&entity.OAuthTokenEntity{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return repositories
}
