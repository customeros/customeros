package service

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/quickbooks"
	"log"
	"reflect"

	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/action"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/ai"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/attachment"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/azure"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/cloudflare"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/comment"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/contact"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/contract"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/currency"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/custom_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/domain"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/email"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/emailing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/enrichment"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	externalsystem "github.com/customeros/customeros/packages/server/customer-os-common-module/services/external_system"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/files"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/flow"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/flow_execution"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/google"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/industry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/interaction_event"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/interaction_session"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/invoice"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/issue"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/jobrole"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/location"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/log_entry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/mail"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/mailbox"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/mailstack"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/markdown_event"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/namecheap"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/notification"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/novu"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/opensrs"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/opportunity"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/organization"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/phone_number"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/postmark"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/registration"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/reminders"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/service_line_item"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/slack"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/social"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/tags"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/tenant"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/tenant_settings"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/user"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/verify"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/workflow"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/workspace"
)

type CommonServices struct {
	// Core infrastructure
	Cache                *caches.Cache
	Events               *events.EventsService
	Neo4jRepositories    *neo4j_repository.Repositories
	PostgresRepositories *postgres_repository.Repositories

	// Services
	ActionService              interfaces.ActionService
	AgentService               interfaces.AgentService
	AgentCapabilityService     interfaces.AgentCapabilityService
	AgentVisitorIDService      *agent.AgentVisitorIDService
	AIService                  interfaces.AIService
	AttachmentService          interfaces.AttachmentService
	AzureService               interfaces.AzureService
	CloudflareService          interfaces.CloudflareService
	CommentService             interfaces.CommentService
	ContactService             interfaces.ContactService
	ContractService            interfaces.ContractService
	CurrencyService            interfaces.CurrencyService
	CustomFieldTemplateService interfaces.CustomFieldTemplateService
	DomainService              interfaces.DomainService
	EmailingService            interfaces.EmailingService
	EmailService               interfaces.EmailService
	EnrichmentService          interfaces.EnrichmentService
	ExternalSystemService      interfaces.ExternalSystemService
	FileService                interfaces.FileService
	FlowExecutionService       interfaces.FlowExecutionService
	FlowService                interfaces.FlowService
	GoogleService              interfaces.GoogleService
	IndustryService            interfaces.IndustryService
	InteractionEventService    interfaces.InteractionEventService
	InteractionSessionService  interfaces.InteractionSessionService
	InvoiceService             interfaces.InvoiceService
	IssueService               interfaces.IssueService
	JobRoleService             interfaces.JobRoleService
	LocationService            interfaces.LocationService
	LogEntryService            interfaces.LogEntryService
	MailboxService             interfaces.MailboxService
	MailService                interfaces.MailService
	MailstackService           interfaces.MailstackService
	MarkdownEventService       interfaces.MarkdownEventService
	NamecheapService           interfaces.NamecheapService
	NotificationService        interfaces.NotificationService
	NovuService                interfaces.NovuService
	OpenSRSService             interfaces.OpenSrsService
	OpportunityService         interfaces.OpportunityService
	OrganizationService        interfaces.OrganizationService
	PhoneNumberService         interfaces.PhoneNumberService
	PostmarkService            interfaces.PostmarkService
	RegistrationService        interfaces.RegistrationService
	ReminderService            interfaces.ReminderService
	ServiceLineItemService     interfaces.ServiceLineItemService
	SlackService               interfaces.SlackService
	SocialService              interfaces.SocialService
	TagService                 interfaces.TagService
	TenantService              interfaces.TenantService
	TenantSettingsService      interfaces.TenantSettingsService
	UserService                interfaces.UserService
	VerifyService              interfaces.VerifyService
	QuickbooksService          interfaces.QuickbooksService
	WorkflowService            interfaces.WorkflowService
	WorkspaceService           interfaces.WorkspaceService
}

type InitOptions struct {
	LoadPersonalEmailProviders bool
	LoadEmailExclusionList     bool
}

func InitCommonServices(
	log logger.Logger,
	neo4jRepositories *neo4j_repository.Repositories,
	postgresRepositories *postgres_repository.Repositories,
	cfg *config.CommonConfig,
	grpcClients *grpc_client.Clients,
	options *InitOptions,
) *CommonServices {
	var err error

	// Base services
	cacheImpl := caches.NewCommonCache()
	eventsImpl := &events.EventsService{}
	if cfg.Infrastructure.RabbitMQConfig.Url != "" {
		eventsImpl, err = events.NewEventsService(cfg.Infrastructure.RabbitMQConfig.Url, log)
		if err != nil {
			log.Fatalf("Cannot start events service")
		}
	}

	// Simple - Services that depend only on base services
	agentImpl := agent.NewAgentService(postgresRepositories)
	aiImpl := ai.NewAIService(log, &cfg.External.AnthropicConfig)
	attachmentImpl := attachment.NewAttachmentService(neo4jRepositories)
	azureImpl := azure.NewAzureService(&cfg.Infrastructure.AzureOAuthConfig, postgresRepositories, neo4jRepositories)
	cloudfareImpl := cloudflare.NewCloudflareService(log, &cfg.External.CloudflareConfig, postgresRepositories)
	commentImpl := comment.NewCommentService(log, neo4jRepositories, eventsImpl)
	currencyImpl := currency.NewCurrencyService(postgresRepositories)
	customFieldTemplateImpl := custom_fields.NewCustomFieldTemplateService(log, neo4jRepositories, eventsImpl)
	domainImpl := domain.NewDomainService(log, cacheImpl, postgresRepositories, neo4jRepositories, eventsImpl)
	emailingImpl := emailing.NewEmailingService(log, postgresRepositories)
	enrichmentImpl := enrichment.NewEnrichmentService(log, &cfg.External, postgresRepositories)
	externalSystemImpl := externalsystem.NewExternalSystemService(log, neo4jRepositories, eventsImpl)
	googleImpl := google.NewGoogleService(&cfg.Infrastructure.GoogleOAuthConfig, postgresRepositories, neo4jRepositories)
	industryImpl := industry.NewIndustryService(log, neo4jRepositories)
	interactionSessionImpl := interaction_session.NewInteractionSessionService(neo4jRepositories)
	markdownEventImpl := markdown_event.NewMarkdownEventService(log, neo4jRepositories, eventsImpl)
	namecheapImpl := namecheap.NewNamecheapService(&cfg.External.NamecheapConfig, postgresRepositories)
	novuImpl := novu.NewNovuService(cfg.External.NovuCofig.ApiKey)
	openSRSImpl := opensrs.NewOpenSRSService(log, &cfg.External.OpenSRSConfig, postgresRepositories)
	phoneNumberImpl := phone_number.NewPhoneNumberService(neo4jRepositories, eventsImpl)
	postmarkImpl := postmark.NewPostmarkService(&cfg.External.PostmarkConfig, postgresRepositories)
	slackImpl := slack.NewSlackService(log, postgresRepositories)
	tagImpl := tags.NewTagService(log, neo4jRepositories, eventsImpl)
	tenantImpl := tenant.NewTenantService(log, neo4jRepositories, postgresRepositories)
	tenantSettingsImpl := tenant_settings.NewTenantSettingsService(log, neo4jRepositories, eventsImpl)
	userImpl := user.NewUserService(neo4jRepositories, postgresRepositories, eventsImpl)
	quickbooksImpl := quickbooks.NewQuickbooksService(&cfg.External.QuickbooksConfig, postgresRepositories)
	workflowImpl := workflow.NewWorkflowService(postgresRepositories)
	workspaceImpl := workspace.NewWorkspaceService(neo4jRepositories)

	// Services that only depend on Simple
	fileImpl := files.NewFileService(log, &cfg.Internal.FileStoreConfig, neo4jRepositories, attachmentImpl)
	notificationImpl := notification.NewNotificationService(log, postgresRepositories, slackImpl)
	reminderImpl := reminders.NewReminderService(neo4jRepositories, novuImpl)
	verifyImpl := verify.NewVerifyService(log, postgresRepositories, cfg, enrichmentImpl)

	// Complex dependencies (ordered by dependency chain)
	emailImpl := email.NewEmailService(neo4jRepositories, eventsImpl, nil, nil, nil)
	jobroleImpl := jobrole.NewJobRoleService(neo4jRepositories, eventsImpl, nil)
	issueImpl := issue.NewIssueService(log, neo4jRepositories, eventsImpl, nil)
	contactImpl := contact.NewContactService(log, neo4jRepositories, eventsImpl, domainImpl, emailImpl, nil, jobroleImpl, nil, nil)
	agentCapabilityImpl, err := agent_capability.NewAgentCapabilityService(postgresRepositories, enrichmentImpl, nil, nil, notificationImpl)
	if err != nil {
		log.Fatalf("Cannot start agent capability service")
	}
	socialImpl := social.NewSocialService(log, neo4jRepositories, eventsImpl, contactImpl)
	orgImpl := organization.NewOrganizationService(log, postgresRepositories, neo4jRepositories, eventsImpl, domainImpl, industryImpl, socialImpl, userImpl)
	agentVisitorIdImpl := agent.NewAgentVisitorIDService(postgresRepositories, agentImpl, agentCapabilityImpl, workspaceImpl)
	contractImpl := contract.NewContractService(log, neo4jRepositories, eventsImpl, grpcClients, nil, orgImpl)
	opportunityImpl := opportunity.NewOpportunityService(log, grpcClients, neo4jRepositories, eventsImpl, contractImpl, orgImpl, tenantSettingsImpl)
	sliImpl := sli.NewServiceLineItemService(log, eventsImpl, neo4jRepositories, contractImpl)
	invoiceImpl := invoice.NewInvoiceService(log, grpcClients, neo4jRepositories, contractImpl, sliImpl, tenantSettingsImpl)
	logEntry := logentry.NewLogEntryService(log, neo4jRepositories, eventsImpl, orgImpl)
	mailboxImpl := mailbox.NewMailboxService(log, postgresRepositories, neo4jRepositories, emailImpl)
	interactionEventImpl := interaction_event.NewInteractionEventService(neo4jRepositories, emailImpl)
	mailstackImpl := mailstack.NewMailstackService(&cfg.External.StripeConfig, eventsImpl, postgresRepositories, cloudfareImpl, namecheapImpl, mailboxImpl, openSRSImpl)
	mailImpl := mail.NewMailService(cacheImpl, postgresRepositories, neo4jRepositories, azureImpl, contactImpl, emailImpl, googleImpl, interactionEventImpl, interactionSessionImpl, openSRSImpl, orgImpl)
	flowExecutionImpl := flow_execution.NewFlowExecutionService(neo4jRepositories, postgresRepositories, eventsImpl, emailImpl, nil, orgImpl, socialImpl)
	flowImpl := flow.NewFlowService(neo4jRepositories, postgresRepositories, eventsImpl, flowExecutionImpl)
	locationImpl := location.NewLocationService(log, neo4jRepositories, postgresRepositories, eventsImpl, &cfg.External.AnthropicConfig.Prompts, aiImpl, contactImpl, orgImpl)
	actionImpl := action.NewActionService(log, neo4jRepositories, eventsImpl, orgImpl)
	registrationImpl := registration.NewRegistrationService(eventsImpl, postgresRepositories, neo4jRepositories, contactImpl, emailImpl, flowImpl, mailboxImpl, orgImpl, postmarkImpl, userImpl)

	// Resolve circular dependencies
	emailImpl.SetContactService(contactImpl)
	emailImpl.SetOrganizationService(orgImpl)
	emailImpl.SetDomainService(domainImpl)
	issueImpl.SetOrganizationService(orgImpl)
	contactImpl.SetOrganizationService(orgImpl)
	contactImpl.SetSocialService(socialImpl)
	contactImpl.SetFlowService(flowImpl)
	contractImpl.SetOpportunityService(opportunityImpl)
	flowExecutionImpl.SetFlowService(flowImpl)
	jobroleImpl.SetOrganizationService(orgImpl)
	agentCapabilityImpl.SetOrganizationService(orgImpl)
	agentCapabilityImpl.SetActionService(actionImpl)

	// Initialize CommonServices struct
	common := CommonServices{
		// Core components
		Cache:                cacheImpl,
		Events:               eventsImpl,
		Neo4jRepositories:    neo4jRepositories,
		PostgresRepositories: postgresRepositories,

		// All other services (alphabetically)
		ActionService:              actionImpl,
		AgentService:               agentImpl,
		AgentCapabilityService:     agentCapabilityImpl,
		AgentVisitorIDService:      agentVisitorIdImpl,
		AIService:                  aiImpl,
		AttachmentService:          attachmentImpl,
		AzureService:               azureImpl,
		CloudflareService:          cloudfareImpl,
		CommentService:             commentImpl,
		ContactService:             contactImpl,
		ContractService:            contractImpl,
		CurrencyService:            currencyImpl,
		CustomFieldTemplateService: customFieldTemplateImpl,
		DomainService:              domainImpl,
		EmailService:               emailImpl,
		EmailingService:            emailingImpl,
		EnrichmentService:          enrichmentImpl,
		ExternalSystemService:      externalSystemImpl,
		FileService:                fileImpl,
		FlowService:                flowImpl,
		FlowExecutionService:       flowExecutionImpl,
		GoogleService:              googleImpl,
		IndustryService:            industryImpl,
		InteractionEventService:    interactionEventImpl,
		InteractionSessionService:  interactionSessionImpl,
		InvoiceService:             invoiceImpl,
		IssueService:               issueImpl,
		JobRoleService:             jobroleImpl,
		LocationService:            locationImpl,
		LogEntryService:            logEntry,
		MailService:                mailImpl,
		MailboxService:             mailboxImpl,
		MailstackService:           mailstackImpl,
		MarkdownEventService:       markdownEventImpl,
		NamecheapService:           namecheapImpl,
		NotificationService:        notificationImpl,
		NovuService:                novuImpl,
		OpenSRSService:             openSRSImpl,
		OpportunityService:         opportunityImpl,
		OrganizationService:        orgImpl,
		PhoneNumberService:         phoneNumberImpl,
		PostmarkService:            postmarkImpl,
		RegistrationService:        registrationImpl,
		ReminderService:            reminderImpl,
		ServiceLineItemService:     sliImpl,
		SlackService:               slackImpl,
		SocialService:              socialImpl,
		TagService:                 tagImpl,
		TenantService:              tenantImpl,
		TenantSettingsService:      tenantSettingsImpl,
		UserService:                userImpl,
		VerifyService:              verifyImpl,
		QuickbooksService:          quickbooksImpl,
		WorkflowService:            workflowImpl,
		WorkspaceService:           workspaceImpl,
	}

	// Check that all services are initialized
	CheckIsInitialized(&common)

	// Process options
	if options != nil {
		if options.LoadPersonalEmailProviders {
			// init app cache
			personalEmailProviderEntities, err := postgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
			if err != nil {
				log.Fatalf("Error getting personal email providers: %s", err.Error())
			}
			personalEmailProviders := make([]string, 0)
			for _, personalEmailProvider := range personalEmailProviderEntities {
				personalEmailProviders = append(personalEmailProviders, personalEmailProvider.ProviderDomain)
			}
			common.Cache.SetPersonalEmailProviders(personalEmailProviders)
		}
		if options.LoadEmailExclusionList {
			// init app cache
			exclusionList, err := postgresRepositories.TenantSettingsEmailExclusionRepository.GetExclusionList(context.Background())
			if err != nil {
				log.Fatalf("Error getting exclusion list: %s", err.Error())
			}
			common.Cache.SetEmailExclusion(exclusionList)
		}
	}

	return &common
}

// CheckIsInitialized iterates over all fields of a struct and calls IsInitialized if the field implements it.
func CheckIsInitialized(common *CommonServices) {
	v := reflect.ValueOf(common).Elem() // struct value
	t := v.Type()                       // struct type

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)

		// We only care if the field is non-nil (pointer or interface)
		if (field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface) && !field.IsNil() {
			fieldValue := field.Interface()
			fieldType := reflect.TypeOf(fieldValue)

			// Check if the underlying type has IsInitialized method
			_, exists := fieldType.MethodByName("IsInitialized")
			if exists {
				// Invoke IsInitialized
				results := reflect.ValueOf(fieldValue).MethodByName("IsInitialized").Call(nil)
				if len(results) == 1 && results[0].Kind() == reflect.Bool {
					isInitialized := results[0].Bool()
					if !isInitialized {
						log.Fatalf("Service %s is not initialized", t.Field(i).Name)
					}
				}
			}
		}
	}
}
