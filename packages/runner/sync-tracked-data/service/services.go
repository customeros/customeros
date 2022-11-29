package service

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/repository"
	"gorm.io/gorm"
)

type Services struct {
	SyncService SyncService
}

func InitServices(driver *neo4j.Driver, gormDb *gorm.DB) *Services {
	repositories := repository.InitRepos(driver, gormDb)

	services := new(Services)

	services.SyncService = NewSyncService(repositories, services)

	return services
}
