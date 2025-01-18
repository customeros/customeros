package service

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	commonServices "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/customeros/packages/server/events/eventstore"

	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/logger"
)

type Services struct {
	CommonServices       *commonServices.CommonServices
	PostgresRepositories *postgres_repository.Repositories
	Neo4jRepositories    *neo4j_repository.Repositories
	Es                   eventstore.AggregateStore
}

func InitServices(
	log logger.Logger,
	neo4jRepositories *neo4j_repository.Repositories,
	postgresRepositories *postgres_repository.Repositories,
	config commonConfig.CommonConfig,
	grpcClients *grpc_client.Clients,
	es eventstore.AggregateStore,
) *Services {
	services := Services{
		Es:                   es,
		PostgresRepositories: postgresRepositories,
		Neo4jRepositories:    neo4jRepositories,
	}
	services.CommonServices = commonServices.InitCommonServices(
		log,
		neo4jRepositories,
		postgresRepositories,
		&config,
		grpcClients,
	)

	return &services
}
