package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Services struct {
	GlobalConfig *config.GlobalConfig
	Cache        *caches.Cache

	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jRepository.Repositories

	GrpcClients *grpc_client.Clients

	AttachmentService         AttachmentService
	ContactService            ContactService
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
	OpportunityService        OpportunityService
	SlackChannelService       SlackChannelService
	ServiceLineItemService    ServiceLineItemService
	TenantService             TenantService
	UserService               UserService
	WorkflowService           WorkflowService
	WorkspaceService          WorkspaceService
	SocialService             SocialService
	DomainService             DomainService
	OrganizationService       OrganizationService
	TagService                TagService

	GoogleService  GoogleService
	AzureService   AzureService
	OpenSrsService OpenSrsService
	MailService    MailService

	ApiCacheService ApiCacheService
}

func InitServices(globalConfig *config.GlobalConfig, db *gorm.DB, driver *neo4j.DriverWithContext, neo4jDatabase string, grpcClients *grpc_client.Clients, log logger.Logger) *Services {
	services := &Services{
		GlobalConfig:         globalConfig,
		Cache:                caches.NewCommonCache(),
		GrpcClients:          grpcClients,
		PostgresRepositories: postgresRepository.InitRepositories(db),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, neo4jDatabase),
	}

	services.CommonService = NewCommonService(services)
	services.ApiCacheService = NewApiCacheService(services.Neo4jRepositories, services)

	services.AttachmentService = NewAttachmentService(services)
	services.AzureService = NewAzureService(globalConfig.AzureOAuthConfig, services.PostgresRepositories, services)
	services.ContactService = NewContactService(log, services)
	services.ContractService = NewContractService(log, services)
	services.CurrencyService = NewCurrencyService(services.PostgresRepositories)
	services.DomainService = NewDomainService(log, services)
	services.EmailService = NewEmailService(services)
	services.EmailingService = NewEmailingService(log, services)
	services.ExternalSystemService = NewExternalSystemService(log, services)
	services.FlowService = NewFlowService(services)
	services.FlowExecutionService = NewFlowExecutionService(services)
	services.GoogleService = NewGoogleService(globalConfig.GoogleOAuthConfig, services.PostgresRepositories, services)
	services.InvoiceService = NewInvoiceService(services)
	services.JobRoleService = NewJobRoleService(services)
	services.InteractionSessionService = NewInteractionSessionService(services)
	services.InteractionEventService = NewInteractionEventService(services)
	services.SocialService = NewSocialService(log, services)
	services.OpportunityService = NewOpportunityService(log, services)
	services.SocialService = NewSocialService(log, services)
	services.MailService = NewMailService(services)
	services.OrganizationService = NewOrganizationService(services)
	services.OpenSrsService = NewOpenSRSService(services)
	services.SlackChannelService = NewSlackChannelService(services.PostgresRepositories)
	services.ServiceLineItemService = NewServiceLineItemService(log, services)
	services.TagService = NewTagService(log, services)
	services.TenantService = NewTenantService(log, services)
	services.UserService = NewUserService(services)
	services.WorkflowService = NewWorkflowService(services)
	services.WorkspaceService = NewWorkspaceService(services)

	//init app cache
	personalEmailProviderEntities, err := services.PostgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
	if err != nil {
		log.Fatalf("Error getting personal email providers: %s", err.Error())
	}

	personalEmailProviders := make([]string, 0)
	for _, personalEmailProvider := range personalEmailProviderEntities {
		personalEmailProviders = append(personalEmailProviders, personalEmailProvider.ProviderDomain)
	}
	services.Cache.SetPersonalEmailProviders(personalEmailProviders)

	return services
}
