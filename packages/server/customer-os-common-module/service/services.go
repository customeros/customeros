package service

import (
	"context"
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
	Logger       logger.Logger

	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jRepository.Repositories

	RabbitMQService RabbitMQService
	GrpcClients     *grpc_client.Clients

	AttachmentService          AttachmentService
	ContactService             ContactService
	ContractService            ContractService
	CommonService              CommonService
	CurrencyService            CurrencyService
	CustomFieldTemplateService CustomFieldTemplateService
	EmailService               EmailService
	EmailingService            EmailingService
	ExternalSystemService      ExternalSystemService
	FlowService                FlowService
	FlowExecutionService       FlowExecutionService
	JobRoleService             JobRoleService
	InvoiceService             InvoiceService
	InteractionSessionService  InteractionSessionService
	InteractionEventService    InteractionEventService
	LocationService            LocationService
	OpportunityService         OpportunityService
	SlackChannelService        SlackChannelService
	ServiceLineItemService     ServiceLineItemService
	TenantService              TenantService
	UserService                UserService
	WorkflowService            WorkflowService
	WorkspaceService           WorkspaceService
	SocialService              SocialService
	DomainService              DomainService
	OrganizationService        OrganizationService
	TagService                 TagService
	RegistrationService        RegistrationService

	GoogleService   GoogleService
	AzureService    AzureService
	OpenSrsService  OpenSrsService
	MailService     MailService
	PostmarkService PostmarkService
	NovuService     NovuService
}

func InitServices(globalConfig *config.GlobalConfig, db *gorm.DB, driver *neo4j.DriverWithContext, neo4jDatabase string, grpcClients *grpc_client.Clients, log logger.Logger) *Services {
	services := &Services{
		GlobalConfig:         globalConfig,
		Cache:                caches.NewCommonCache(),
		Logger:               log,
		GrpcClients:          grpcClients,
		PostgresRepositories: postgresRepository.InitRepositories(db),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, neo4jDatabase),
	}

	if globalConfig.RabbitMQConfig != nil {
		services.RabbitMQService = NewRabbitMQService(globalConfig.RabbitMQConfig.Url, services)
	}

	services.CommonService = NewCommonService(services)

	services.AttachmentService = NewAttachmentService(services)
	services.AzureService = NewAzureService(globalConfig.AzureOAuthConfig, services.PostgresRepositories, services)
	services.ContactService = NewContactService(log, services)
	services.ContractService = NewContractService(log, services)
	services.CurrencyService = NewCurrencyService(services.PostgresRepositories)
	services.CustomFieldTemplateService = NewCustomFieldTemplateService(log, services)
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
	services.LocationService = NewLocationService(log, services)
	services.SocialService = NewSocialService(log, services)
	services.OpportunityService = NewOpportunityService(log, services)
	services.SocialService = NewSocialService(log, services)
	services.MailService = NewMailService(services)
	services.PostmarkService = NewPostmarkService(services)
	services.NovuService = NewNovuService(services)
	services.OrganizationService = NewOrganizationService(services)
	services.OpenSrsService = NewOpenSRSService(services)
	services.SlackChannelService = NewSlackChannelService(services.PostgresRepositories)
	services.ServiceLineItemService = NewServiceLineItemService(log, services)
	services.TagService = NewTagService(log, services)
	services.TenantService = NewTenantService(log, services)
	services.UserService = NewUserService(services)
	services.WorkflowService = NewWorkflowService(services)
	services.WorkspaceService = NewWorkspaceService(services)
	services.RegistrationService = NewRegistrationService(services)

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

	emailExclusionEntities, err := services.PostgresRepositories.TenantSettingsEmailExclusionRepository.GetExclusionList(context.Background())
	if err != nil {
		log.Fatalf("Error getting email exclusion list: %s", err.Error())
	}
	services.Cache.SetEmailExclusion(emailExclusionEntities)

	return services
}
