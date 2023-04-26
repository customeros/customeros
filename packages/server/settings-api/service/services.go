package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"gorm.io/gorm"
)

type Services struct {
	TenantSettingsService TenantSettingsService
}

func InitServices(db *gorm.DB) *Services {
	repositories := repository.InitRepositories(db)

	err := db.AutoMigrate(entity.TenantSettings{})
	if err != nil {
		panic(err)
	}

	return &Services{
		TenantSettingsService: NewTenantSettingsService(repositories),
	}
}
