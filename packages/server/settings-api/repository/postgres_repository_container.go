package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresRepositories struct {
	TenantSettingsRepository TenantSettingsRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositories {
	p := &PostgresRepositories{
		TenantSettingsRepository: NewTenantSettingsRepository(db),
	}

	err := db.AutoMigrate(&entity.TenantSettings{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
