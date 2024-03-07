package repository

import (
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres"
	authEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	OAuthTokenRepository    repository.OAuthTokenRepository
	SlackSettingsRepository repository.SlackSettingsRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		OAuthTokenRepository:    repository.NewOAuthTokenRepository(db),
		SlackSettingsRepository: repository.NewSlackSettingsRepository(db),
	}

	return repositories
}

func Migration(db *gorm.DB) {

	var err error

	err = db.AutoMigrate(&authEntity.OAuthTokenEntity{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&authEntity.SlackSettingsEntity{})
	if err != nil {
		panic(err)
	}
}
