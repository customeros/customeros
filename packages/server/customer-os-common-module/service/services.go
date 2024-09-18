package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
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

	AttachmentService         AttachmentService
	ContractService           ContractService
	CommonService             CommonService
	CurrencyService           CurrencyService
	EmailService              EmailService
	EmailingService           EmailingService
	ExternalSystemService     ExternalSystemService
	FlowService               FlowService
	FlowExecutionService      FlowExecutionService
	JobRoleService            JobRoleService
	InvoiceService            InvoiceService
	InteractionSessionService InteractionSessionService
	InteractionEventService   InteractionEventService
	SlackChannelService       SlackChannelService
	ServiceLineItemService    ServiceLineItemService
	TenantService             TenantService
	UserService               UserService
	WorkflowService           WorkflowService
	WorkspaceService          WorkspaceService
	SocialService             SocialService
	DomainService             DomainService

	GoogleService  GoogleService
	AzureService   AzureService
	OpenSrsService OpenSrsService
	MailService    MailService

	ApiCacheService ApiCacheService
}

func InitServices(globalConfig *config.GlobalConfig, db *gorm.DB, driver *neo4j.DriverWithContext, neo4jDatabase string, grpcClients *grpc_client.Clients) *Services {
	services := &Services{
		GrpcClients:          grpcClients,
		PostgresRepositories: postgresRepository.InitRepositories(db),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, neo4jDatabase),
	}

	cache := caches.NewCommonCache()

	services.CommonService = NewCommonService(services)
	services.ApiCacheService = NewApiCacheService(services.Neo4jRepositories, services)

	services.AttachmentService = NewAttachmentService(services)
	services.AzureService = NewAzureService(globalConfig.AzureOAuthConfig, services.PostgresRepositories, services)
	services.ContractService = NewContractService(nil, services)
	services.CurrencyService = NewCurrencyService(services.PostgresRepositories)
	services.DomainService = NewDomainService(nil, services, cache)
	services.EmailService = NewEmailService(services)
	services.EmailingService = NewEmailingService(nil, services)
	services.ExternalSystemService = NewExternalSystemService(nil, services)
	services.FlowService = NewFlowService(services)
	services.FlowExecutionService = NewFlowExecutionService(services)
	services.GoogleService = NewGoogleService(globalConfig.GoogleOAuthConfig, services.PostgresRepositories, services)
	services.InvoiceService = NewInvoiceService(services)
	services.JobRoleService = NewJobRoleService(services)
	services.InteractionSessionService = NewInteractionSessionService(services)
	services.InteractionEventService = NewInteractionEventService(services)
	services.SocialService = NewSocialService(nil, services)
	services.MailService = NewMailService(services)
	services.OpenSrsService = NewOpenSRSService(services)
	services.SlackChannelService = NewSlackChannelService(services.PostgresRepositories)
	services.ServiceLineItemService = NewServiceLineItemService(nil, services)
	services.TenantService = NewTenantService(nil, services)
	services.UserService = NewUserService(services)
	services.WorkflowService = NewWorkflowService(services)
	services.WorkspaceService = NewWorkspaceService(services)

	return services
}
