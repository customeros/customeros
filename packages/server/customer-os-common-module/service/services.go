package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Services struct {
	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *repository.Repositories

	StateService        StateService
	SlackChannelService SlackChannelService
}

func InitServices(db *gorm.DB, driver *neo4j.DriverWithContext) *Services {
	services := &Services{
		PostgresRepositories: postgresRepository.InitRepositories(db),
		Neo4jRepositories:    repository.InitRepositories(driver),
	}

	services.StateService = NewStateService(services.Neo4jRepositories)
	services.SlackChannelService = NewSlackChannelService(services.PostgresRepositories)

	return services
}
