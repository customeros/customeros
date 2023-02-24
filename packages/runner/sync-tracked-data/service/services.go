package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/repository"
	"gorm.io/gorm"
)

type Services struct {
	SyncService SyncService
	InitService InitService
}

func InitServices(driver *neo4j.DriverWithContext, gormDb *gorm.DB, gormTrackingDb *gorm.DB) *Services {
	repositories := repository.InitRepos(driver, gormDb, gormTrackingDb)

	services := new(Services)

	services.SyncService = NewSyncService(repositories, services)
	services.InitService = NewInitService(repositories)

	return services
}
