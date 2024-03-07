package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"gorm.io/gorm"
)

type Services struct {
	CommonRepositories *repository.Repositories

	StateService        StateService
	SlackChannelService SlackChannelService
}

func InitServices(db *gorm.DB, driver *neo4j.DriverWithContext) *Services {
	repositories := repository.InitRepositories(db, driver)

	services := &Services{
		CommonRepositories:  repositories,
		StateService:        NewStateService(repositories),
		SlackChannelService: NewSlackChannelService(repositories),
	}

	return services
}
