package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"gorm.io/gorm"
)

type Services struct {
	TenantSettingsService    TenantSettingsService
	OAuthUserSettingsService OAuthUserSettingsService
}

func InitServices(db *gorm.DB, logger logger.Logger) *Services {
	repositories := repository.InitRepositories(db)

	err := db.AutoMigrate(entity.TenantSettings{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(entity.TenantAPIKey{})
	if err != nil {
		panic(err)
	}

	return &Services{
		TenantSettingsService:    NewTenantSettingsService(repositories, logger),
		OAuthUserSettingsService: NewUserSettingsService(repositories, logger),
	}
}
