package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonServices "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	neo4jRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
)

type Services struct {
	CommonServices       *commonServices.CommonServices
	PostgresRepositories *repository.Repositories
	Neo4jRepositories    *neo4jRepo.Repositories
	Es                   eventstore.AggregateStore
}

func InitServices(
	log logger.Logger,
	neo4jRepositories *neo4jRepo.Repositories,
	postgresRepositories *repository.Repositories,
	config *commonConfig.CommonConfig,
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
		config,
		grpcClients,
	)

	return &services
}
