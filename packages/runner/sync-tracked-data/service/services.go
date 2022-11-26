package service

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/gen"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/repository"
)

type Services struct {
	SyncService    SyncService
	VisitorService VisitorService
}

func InitServices(driver *neo4j.Driver, client *gen.Client) *Services {
	repositories := repository.InitRepos(driver, client)

	services := new(Services)

	services.SyncService = NewSyncService(repositories, services)
	services.VisitorService = NewVisitorService(&repositories.TrackedVisitorRepository)

	return services
}
