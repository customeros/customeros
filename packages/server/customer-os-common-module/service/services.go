package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Services struct {
	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jRepository.Repositories

	GrpcClients *grpc_client.Clients

	TenantService          TenantService
	ExternalSystemService  ExternalSystemService
	ContractService        ContractService
	ServiceLineItemService ServiceLineItemService
	InvoiceService         InvoiceService
	SlackChannelService    SlackChannelService
	CurrencyService        CurrencyService
	WorkflowService        WorkflowService
	SocialService          SocialService

	GoogleService GoogleService
	AzureService  AzureService

	ApiCacheService ApiCacheService
}

func InitServices(globalConfig *config.GlobalConfig, db *gorm.DB, driver *neo4j.DriverWithContext, neo4jDatabase string, grpcClients *grpc_client.Clients) *Services {
	services := &Services{
		GrpcClients:          grpcClients,
		PostgresRepositories: postgresRepository.InitRepositories(db),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, neo4jDatabase),
	}

	services.TenantService = NewTenantService(nil, services)
	services.ExternalSystemService = NewExternalSystemService(nil, services)
	services.ContractService = NewContractService(nil, services)
	services.ServiceLineItemService = NewServiceLineItemService(nil, services)
	services.InvoiceService = NewInvoiceService(nil, services)
	services.SlackChannelService = NewSlackChannelService(services.PostgresRepositories)
	services.CurrencyService = NewCurrencyService(services.PostgresRepositories)
	services.WorkflowService = NewWorkflowService(nil, services)
	services.SocialService = NewSocialService(nil, services)

	services.GoogleService = NewGoogleService(globalConfig.GoogleOAuthConfig, services.PostgresRepositories, services)
	services.AzureService = NewAzureService(services.PostgresRepositories, services)

	services.ApiCacheService = NewApiCacheService(services.Neo4jRepositories, services)

	return services
}
