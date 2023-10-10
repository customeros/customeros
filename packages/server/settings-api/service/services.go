package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"gorm.io/gorm"
)

type Services struct {
	Repositories *repository.PostgresRepositories

	TenantSettingsService       TenantSettingsService
	PersonalIntegrationsService PersonalIntegrationsService
	OAuthUserSettingsService    OAuthUserSettingsService
	SlackSettingsService        SlackSettingsService
}

func InitServices(db *gorm.DB, driver *neo4j.DriverWithContext, logger logger.Logger) *Services {
	repositories := repository.InitRepositories(db, driver)

	err := db.AutoMigrate(entity.TenantSettings{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(entity.TenantAPIKey{})
	if err != nil {
		panic(err)
	}

	return &Services{
		Repositories:                repositories,
		TenantSettingsService:       NewTenantSettingsService(repositories, logger),
		OAuthUserSettingsService:    NewUserSettingsService(repositories, logger),
		SlackSettingsService:        NewSlackSettingsService(repositories, logger),
		PersonalIntegrationsService: NewPersonalIntegrationsService(repositories, logger),
	}
}
