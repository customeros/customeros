package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/ai"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/attachment"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/azure"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/cloudflare"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/comment"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/currency"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/custom_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/domain"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/emailing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/enrichment"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/google"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/interaction_session"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/markdown_event"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/namecheap"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/novu"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/opensrs"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/phone_number"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/postmark"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/reminders"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/slack"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/tags"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/tenant_settings"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/verify"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/workflow"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/workspace"
	neo4jRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"log"
	"reflect"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/action"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/email"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/files"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/flow"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/flow_execution"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/industry"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/interaction_event"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/issue"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/jobrole"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/location"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/log_entry"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/mail"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/mailbox"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/mailstack"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/registration"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/service_line_item"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/social"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/user"
)

type CommonServices struct {
	Neo4jRepositories    *neo4jRepo.Repositories
	PostgresRepositories *repository.Repositories

	Cache                      *caches.Cache
	Events                     *events.EventsService
	ActionService              interfaces.ActionService
	AIService                  interfaces.AIService
	AgentService               interfaces.AgentService
	AttachmentService          interfaces.AttachmentService
	AzureService               interfaces.AzureService
	CloudflareService          interfaces.CloudflareService
	CommentService             interfaces.CommentService
	ContactService             interfaces.ContactService
	ContractService            interfaces.ContractService
	CurrencyService            interfaces.CurrencyService
	CustomFieldTemplateService interfaces.CustomFieldTemplateService
	DomainService              interfaces.DomainService
	EmailService               interfaces.EmailService
	EmailingService            interfaces.EmailingService
	EnrichmentService          interfaces.EnrichmentService
	ExternalSystemService      interfaces.ExternalSystemService
	FileService                interfaces.FileService
	FlowService                interfaces.FlowService
	FlowExecutionService       interfaces.FlowExecutionService
	GoogleService              interfaces.GoogleService
	IndustryService            interfaces.IndustryService
	InteractionEventService    interfaces.InteractionEventService
	InteractionSessionService  interfaces.InteractionSessionService
	InvoiceService             interfaces.InvoiceService
	IssueService               interfaces.IssueService
	JobRoleService             interfaces.JobRoleService
	LocationService            interfaces.LocationService
	LogEntryService            interfaces.LogEntryService
	MailService                interfaces.MailService
	MailboxService             interfaces.MailboxService
	MailstackService           interfaces.MailstackService
	MarkdownEventService       interfaces.MarkdownEventService
	NamecheapService           interfaces.NamecheapService
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
	WorkflowService            interfaces.WorkflowService
	WorkspaceService           interfaces.WorkspaceService
}

func InitCommonServices(
	log logger.Logger,
	neo4jRepositories *neo4jRepo.Repositories,
	postgresRepositories *repository.Repositories,
	cfg *config.CommonConfig,
	grpcClients *grpc_client.Clients,
) *CommonServices {
	var err error

	eventsImpl := &events.EventsService{}
	if cfg.Infrastructure.RabbitMQConfig.Url != "" {
		eventsImpl, err = events.NewEventsService(cfg.Infrastructure.RabbitMQConfig.Url, log)
		if err != nil {
			log.Fatalf("Cannot start events service")
		}
	}

	cacheImpl := caches.NewCommonCache()
	emailImpl := email.NewEmailService(neo4jRepositories, eventsImpl, nil, nil, nil)
	aiImpl := ai.NewAIService(&cfg.External.AnthropicConfig)
	attachmentImpl := attachment.NewAttachmentService(neo4jRepositories)
	azureImpl := azure.NewAzureService(&cfg.Infrastructure.AzureOAuthConfig, postgresRepositories, neo4jRepositories)
	cloudfareImpl := cloudflare.NewCloudflareService(log, &cfg.External.CloudflareConfig, postgresRepositories)
	currencyImpl := currency.NewCurrencyService(postgresRepositories)
	emailingImpl := emailing.NewEmailingService(log, postgresRepositories)
	enrichmentImpl := enrichment.NewEnrichmentService(log, &cfg.External, postgresRepositories)
	googleImpl := google.NewGoogleService(&cfg.Infrastructure.GoogleOAuthConfig, postgresRepositories, neo4jRepositories)
	industryImpl := industry.NewIndustryService(log, neo4jRepositories)
	interactionSessionImpl := interaction_session.NewInteractionSessionService(neo4jRepositories)
	namecheapImpl := namecheap.NewNamecheapService(&cfg.External.NamecheapConfig, postgresRepositories)
	novuImpl := novu.NewNovuService(cfg.External.NovuCofig.ApiKey)
	openSRSImpl := opensrs.NewOpenSRSService(log, &cfg.External.OpenSRSConfig, postgresRepositories)
	postmarkImpl := postmark.NewPostmarkService(&cfg.External.PostmarkConfig, postgresRepositories)
	slackImpl := slack.NewSlackService(log, postgresRepositories)
	tenantImpl := tenant.NewTenantService(log, neo4jRepositories, postgresRepositories)
	workflowImpl := workflow.NewWorkflowService(postgresRepositories)
	workspaceImpl := workspace.NewWorkspaceService(neo4jRepositories)
	commentImpl := comment.NewCommentService(log, neo4jRepositories, eventsImpl)
	customFieldTemplateImpl := custom_fields.NewCustomFieldTemplateService(log, neo4jRepositories, eventsImpl)
	markdownEventImpl := markdown_event.NewMarkdownEventService(log, neo4jRepositories, eventsImpl)
	phoneNumberImpl := phone_number.NewPhoneNumberService(neo4jRepositories, eventsImpl)
	tagImpl := tags.NewTagService(log, neo4jRepositories, eventsImpl)
	tenantSettingsImpl := tenant_settings.NewTenantSettingsService(log, neo4jRepositories, eventsImpl)
	userImpl := user.NewUserService(neo4jRepositories, postgresRepositories, eventsImpl)
	verifyImpl := verify.NewVerifyService(log, postgresRepositories, cfg, enrichmentImpl)
	domainImpl := domain.NewDomainService(log, cacheImpl, postgresRepositories, neo4jRepositories, eventsImpl)
	fileImpl := files.NewFileService(log, &cfg.Internal.FileStoreConfig, neo4jRepositories, attachmentImpl)
	reminderImpl := reminders.NewReminderService(neo4jRepositories, novuImpl)
	jobroleImpl := jobrole.NewJobRoleService(neo4jRepositories, eventsImpl)
	issueImpl := issue.NewIssueService(log, neo4jRepositories, eventsImpl, nil)
	contactImpl := contact.NewContactService(log, neo4jRepositories, eventsImpl, domainImpl, emailImpl, nil, jobroleImpl, nil, nil)
	socialImpl := social.NewSocialService(log, neo4jRepositories, eventsImpl, contactImpl)
	orgImpl := organization.NewOrganizationService(log, postgresRepositories, neo4jRepositories, eventsImpl, domainImpl, industryImpl, socialImpl, userImpl)
	contractImpl := contract.NewContractService(log, neo4jRepositories, eventsImpl, grpcClients, nil, orgImpl)
	opportunityImpl := opportunity.NewOpportunityService(log, grpcClients, neo4jRepositories, eventsImpl, contractImpl, orgImpl, tenantSettingsImpl)
	sliImpl := sli.NewServiceLineItemService(log, eventsImpl, neo4jRepositories, contractImpl)
	invoiceImpl := invoice.NewInvoiceService(log, grpcClients, neo4jRepositories, contractImpl, sliImpl, tenantSettingsImpl)
	logEntry := logentry.NewLogEntryService(log, neo4jRepositories, eventsImpl, orgImpl)
	mailboxImpl := mailbox.NewMailboxService(log, postgresRepositories, neo4jRepositories, emailImpl)
	mailstackImpl := mailstack.NewMailstackService(&cfg.External.StripeConfig, eventsImpl, postgresRepositories, cloudfareImpl, namecheapImpl, mailboxImpl, openSRSImpl)
	interactionEventImpl := interaction_event.NewInteractionEventService(neo4jRepositories, emailImpl)
	mailImpl := mail.NewMailService(cacheImpl, postgresRepositories, neo4jRepositories, azureImpl, contactImpl, emailImpl, googleImpl, interactionEventImpl, interactionSessionImpl, openSRSImpl, orgImpl)
	flowExecutionImpl := flow_execution.NewFlowExecutionService(neo4jRepositories, postgresRepositories, eventsImpl, emailImpl, nil, orgImpl, socialImpl)
	flowImpl := flow.NewFlowService(neo4jRepositories, eventsImpl, flowExecutionImpl)
	registrationImpl := registration.NewRegistrationService(eventsImpl, postgresRepositories, neo4jRepositories, contactImpl, emailImpl, flowImpl, mailboxImpl, orgImpl, postmarkImpl, userImpl)
	actionImpl := action.NewActionService(log, neo4jRepositories, eventsImpl, orgImpl)
	//visitorIDAgentImpl := agent.NewAgentVisitorIDService(postgresRepositories, actionImpl, enrichmentImpl, orgImpl, slackImpl, workspaceImpl)
	locationImpl := location.NewLocationService(log, neo4jRepositories, postgresRepositories, eventsImpl, &cfg.External.AnthropicPrompts, aiImpl, contactImpl, orgImpl)

	// resolve dependencies
	emailImpl.SetContactService(contactImpl)
	emailImpl.SetOrganizationService(orgImpl)
	emailImpl.SetDomainService(domainImpl)
	issueImpl.SetOrganizationService(orgImpl)
	contactImpl.SetOrganizationService(orgImpl)
	contactImpl.SetSocialService(socialImpl)
	contactImpl.SetFlowService(flowImpl)
	contractImpl.SetOpportunityService(opportunityImpl)
	flowExecutionImpl.SetFlowService(flowImpl)

	// initialize base services
	common := CommonServices{
		Neo4jRepositories:          neo4jRepositories,
		PostgresRepositories:       postgresRepositories,
		Cache:                      cacheImpl,
		Events:                     eventsImpl,
		EmailService:               emailImpl,
		AIService:                  aiImpl,
		AttachmentService:          attachmentImpl,
		AzureService:               azureImpl,
		CloudflareService:          cloudfareImpl,
		CurrencyService:            currencyImpl,
		EmailingService:            emailingImpl,
		EnrichmentService:          enrichmentImpl,
		GoogleService:              googleImpl,
		IndustryService:            industryImpl,
		InteractionSessionService:  interactionSessionImpl,
		NamecheapService:           namecheapImpl,
		NovuService:                novuImpl,
		OpenSRSService:             openSRSImpl,
		PostmarkService:            postmarkImpl,
		SlackService:               slackImpl,
		TenantService:              tenantImpl,
		WorkflowService:            workflowImpl,
		WorkspaceService:           workspaceImpl,
		CommentService:             commentImpl,
		CustomFieldTemplateService: customFieldTemplateImpl,
		MarkdownEventService:       markdownEventImpl,
		PhoneNumberService:         phoneNumberImpl,
		TagService:                 tagImpl,
		TenantSettingsService:      tenantSettingsImpl,
		UserService:                userImpl,
		VerifyService:              verifyImpl,
		DomainService:              domainImpl,
		FileService:                fileImpl,
		ReminderService:            reminderImpl,
		JobRoleService:             jobroleImpl,
		IssueService:               issueImpl,
		ContactService:             contactImpl,
		SocialService:              socialImpl,
		OrganizationService:        orgImpl,
		ContractService:            contractImpl,
		OpportunityService:         opportunityImpl,
		ServiceLineItemService:     sliImpl,
		InvoiceService:             invoiceImpl,
		LogEntryService:            logEntry,
		MailboxService:             mailboxImpl,
		MailstackService:           mailstackImpl,
		InteractionEventService:    interactionEventImpl,
		MailService:                mailImpl,
		FlowExecutionService:       flowExecutionImpl,
		FlowService:                flowImpl,
		RegistrationService:        registrationImpl,
		ActionService:              actionImpl,
		LocationService:            locationImpl,
		//VisitorIDAgentService:      visitorIDAgentImpl,
	}

	// Check that all services are initialized
	CheckIsInitialized(&common)

	return &common
}

// CheckIsInitialized iterates over all fields of a struct and calls IsInitialized if the field implements it.
func CheckIsInitialized(common *CommonServices) {
	v := reflect.ValueOf(common).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			method := field.MethodByName("IsInitialized")
			if method.IsValid() {
				// Call IsInitialized method
				results := method.Call(nil)
				if len(results) == 1 && results[0].Kind() == reflect.Bool {
					isInitialized := results[0].Bool()
					if !isInitialized {
						log.Fatalf("Service %s not initialized", t.Field(i).Name)
					}
				}
			}
		}
	}
}
