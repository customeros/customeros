package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
	"gorm.io/gorm"
)

type Services struct {
	TenantSettingsService TenantSettingsService
}

func InitServices(db *gorm.DB) *Services {
	repositories := repository.InitRepositories(db)

	return &Services{
		TenantSettingsService: NewTenantSettingsService(repositories),
	}
}
