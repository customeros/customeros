package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Services struct {
	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jRepository.Repositories

	GrpcClients *grpc_client.Clients

	TenantService       TenantService
	ContractService     ContractService
	InvoiceService      InvoiceService
	SlackChannelService SlackChannelService
}

func InitServices(db *gorm.DB, driver *neo4j.DriverWithContext, neo4jDatabase string, grpcClients *grpc_client.Clients) *Services {
	services := &Services{
		GrpcClients:          grpcClients,
		PostgresRepositories: postgresRepository.InitRepositories(db),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, neo4jDatabase),
	}

	services.SlackChannelService = NewSlackChannelService(services.PostgresRepositories)
	services.TenantService = NewTenantService(nil, services)
	services.ContractService = NewContractService(nil, services)
	services.InvoiceService = NewInvoiceService(nil, services)

	return services
}
