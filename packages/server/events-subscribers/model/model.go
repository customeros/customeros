package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	service "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	neo4jRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
)

type DependencyContainer struct {
	Logger               logger.Logger
	GRPCClients          *grpc_client.Clients
	CommonConfig         *config.CommonConfig
	PostgresRepositories *repository.Repositories
	Neo4jRepositories    *neo4jRepo.Repositories
	CommonServices       *service.CommonServices
}
