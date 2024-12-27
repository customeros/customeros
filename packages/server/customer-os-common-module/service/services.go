package service

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type Services struct {
	GlobalConfig *config.GlobalConfig
	Cache        *caches.Cache
	Logger       logger.Logger

	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jRepository.Repositories

	RabbitMQService RabbitMQService
	GrpcClients     *grpc_client.Clients

	AgentService               AgentService
	AttachmentService          AttachmentService
	AzureService               AzureService
	CloudflareService          CloudflareService
	ContactService             ContactService
	ContractService            ContractService
	CommonService              CommonService
	CommentService             CommentService
	CurrencyService            CurrencyService
	CustomFieldTemplateService CustomFieldTemplateService
	DomainService              DomainService
	EmailService               EmailService
	EmailingService            EmailingService
	EnrichmentService          EnrichmentService
	ExternalSystemService      ExternalSystemService
	FlowExecutionService       FlowExecutionService
	FlowService                FlowService
	GoogleService              GoogleService
	InteractionSessionService  InteractionSessionService
	InteractionEventService    InteractionEventService
	IssueService               IssueService
	InvoiceService             InvoiceService
	JobRoleService             JobRoleService
	LocationService            LocationService
	LogEntryService            LogEntryService
	MailboxService             MailboxService
	MailService                MailService
	MailstackService           MailstackService
	MarkdownEventService       MarkdownEventService
	NamecheapService           NamecheapService
	NovuService                NovuService
	OpenSrsService             OpenSrsService
	OpportunityService         OpportunityService
	OrganizationService        OrganizationService
	PhoneNumberService         PhoneNumberService
	PostmarkService            PostmarkService
	RegistrationService        RegistrationService
	ReminderService            ReminderService
	ServiceLineItemService     ServiceLineItemService
	SlackService               SlackService
	SocialService              SocialService
	TagService                 TagService
	TenantService              TenantService
	TenantSettingsService      TenantSettingsService
	UserService                UserService
	VerifyService              VerifyService
	WorkflowService            WorkflowService
	WorkspaceService           WorkspaceService
}

func InitServices(globalConfig *config.GlobalConfig, postgresDB *config.PostgresDB, driver *neo4j.DriverWithContext, neo4jDatabase string, grpcClients *grpc_client.Clients, log logger.Logger) *Services {
	services := &Services{
		GlobalConfig:         globalConfig,
		Cache:                caches.NewCommonCache(),
		Logger:               log,
		GrpcClients:          grpcClients,
		PostgresRepositories: postgresRepository.InitRepositories(postgresDB),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, neo4jDatabase),
	}

	if globalConfig.RabbitMQConfig != nil {
		services.RabbitMQService = NewRabbitMQService(globalConfig.RabbitMQConfig.Url, services)
	}

	services.AgentService = NewAgentService(services)
	services.AttachmentService = NewAttachmentService(services)
	services.AzureService = NewAzureService(globalConfig.AzureOAuthConfig, services.PostgresRepositories, services)
	services.CloudflareService = NewCloudflareService(log, services, globalConfig)
	services.CommonService = NewCommonService(services)
	services.ContactService = NewContactService(log, services)
	services.ContractService = NewContractService(log, services)
	services.CurrencyService = NewCurrencyService(services.PostgresRepositories)
	services.CustomFieldTemplateService = NewCustomFieldTemplateService(log, services)
	services.CommentService = NewCommentService(log, services)
	services.DomainService = NewDomainService(log, services)
	services.EmailService = NewEmailService(services)
	services.EmailingService = NewEmailingService(log, services)
	services.EnrichmentService = NewEnrichmentService(services, globalConfig)
	services.ExternalSystemService = NewExternalSystemService(log, services)
	services.FlowExecutionService = NewFlowExecutionService(services)
	services.FlowService = NewFlowService(services)
	services.GoogleService = NewGoogleService(globalConfig.GoogleOAuthConfig, services.PostgresRepositories, services)
	services.InteractionEventService = NewInteractionEventService(services)
	services.InteractionSessionService = NewInteractionSessionService(services)
	services.InvoiceService = NewInvoiceService(services)
	services.JobRoleService = NewJobRoleService(services)
	services.InteractionSessionService = NewInteractionSessionService(services)
	services.InteractionEventService = NewInteractionEventService(services)
	services.IssueService = NewIssueService(log, services)
	services.LocationService = NewLocationService(log, services)
	services.LogEntryService = NewLogEntryService(log, services)
	services.MailboxService = NewMailboxService(log, services)
	services.MailService = NewMailService(services)
	services.MailstackService = NewMailstackService(globalConfig, services)
	services.MarkdownEventService = NewMarkdownEventService(log, services)
	services.NamecheapService = NewNamecheapService(globalConfig, services)
	services.NovuService = NewNovuService(services)
	services.OpenSrsService = NewOpenSRSService(log, services)
	services.OpportunityService = NewOpportunityService(log, services)
	services.OrganizationService = NewOrganizationService(services)
	services.PhoneNumberService = NewPhoneNumberService(services)
	services.PostmarkService = NewPostmarkService(services)
	services.RegistrationService = NewRegistrationService(services)
	services.ReminderService = NewReminderService(services)
	services.ServiceLineItemService = NewServiceLineItemService(log, services)
	services.SlackService = NewSlackService(services, globalConfig)
	services.SocialService = NewSocialService(log, services)
	services.TagService = NewTagService(log, services)
	services.TenantService = NewTenantService(log, services)
	services.TenantSettingsService = NewTenantSettingsService(log, services)
	services.UserService = NewUserService(services)
	services.VerifyService = NewVerifyService(services)
	services.WorkflowService = NewWorkflowService(services)
	services.WorkspaceService = NewWorkspaceService(services)

	// init app cache
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
