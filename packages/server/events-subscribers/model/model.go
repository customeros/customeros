package model

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type DependencyContainer struct {
	Logger               logger.Logger
	GRPCClients          *grpc_client.Clients
	CommonConfig         *config.CommonConfig
	PostgresRepositories *postgres_repository.Repositories
	Neo4jRepositories    *neo4j_repository.Repositories
	CommonServices       *service.CommonServices
}
