package service

import (
	neo4jRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/action"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/agent"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/ai"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/attachment"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/azure"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/cloudflare"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/comment"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/currency"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/custom_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/domain"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/email"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/emailing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/enrichment"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/external_system"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/files"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/flow"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/flow_execution"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/google"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/industry"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/interaction_event"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/interaction_session"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/issue"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/jobrole"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/location"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/log_entry"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/mail"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/mailbox"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/mailstack"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/markdown_event"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/namecheap"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/notification"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/novu"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/opensrs"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/phone_number"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/postmark"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/registration"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/reminders"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/service_line_item"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/slack"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/social"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/tags"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/tenant_settings"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/user"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/verify"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/workflow"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/workspace"
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
	// initialize base services
	events, err := events.NewEventsService(cfg.Infrastructure.RabbitMQConfig.Url, log)
	if err != nil {
		log.Fatalf("Cannot start events service")
	}

	// Directly add services that don't require others services
	common := CommonServices{
		Neo4jRepositories:    neo4jRepositories,
		PostgresRepositories: postgresRepositories,

		Cache:                      caches.NewCommonCache(),
		AIService:                  ai.NewAIService(&cfg.External.AnthropicConfig),
		AttachmentService:          attachment.NewAttachmentService(neo4jRepositories),
		AzureService:               azure.NewAzureService(&cfg.Infrastructure.AzureOAuthConfig, postgresRepositories, neo4jRepositories),
		CloudflareService:          cloudflare.NewCloudflareService(log, &cfg.External.CloudflareConfig, postgresRepositories),
		CommentService:             comment.NewCommentService(log, neo4jRepositories, events),
		CurrencyService:            currency.NewCurrencyService(postgresRepositories),
		CustomFieldTemplateService: custom_fields.NewCustomFieldTemplateService(log, neo4jRepositories, events),
		EmailingService:            emailing.NewEmailingService(log, postgresRepositories),
		EnrichmentService:          enrichment.NewEnrichmentService(log, &cfg.External, postgresRepositories),
		Events:                     events,
		ExternalSystemService:      externalsystem.NewExternalSystemService(log, neo4jRepositories, events),
		GoogleService:              google.NewGoogleService(&cfg.Infrastructure.GoogleOAuthConfig, postgresRepositories, neo4jRepositories),
		IndustryService:            industry.NewIndustryService(log, neo4jRepositories),
		InteractionSessionService:  interaction_session.NewInteractionSessionService(neo4jRepositories),
		MarkdownEventService:       markdown_event.NewMarkdownEventService(log, neo4jRepositories, events),
		NotificationService:        notification.NewNotificationService(log, postgresRepositories),
		NamecheapService:           namecheap.NewNamecheapService(&cfg.External.NamecheapConfig, postgresRepositories),
		NovuService:                novu.NewNovuService(cfg.External.NovuCofig.ApiKey),
		OpenSRSService:             opensrs.NewOpenSRSService(log, &cfg.External.OpenSRSConfig, postgresRepositories),
		PhoneNumberService:         phone_number.NewPhoneNumberService(neo4jRepositories, events),
		PostmarkService:            postmark.NewPostmarkService(&cfg.External.PostmarkConfig, postgresRepositories),
		SlackService:               slack.NewSlackService(postgresRepositories),
		TagService:                 tags.NewTagService(log, neo4jRepositories, events),
		TenantService:              tenant.NewTenantService(log, neo4jRepositories, postgresRepositories),
		TenantSettingsService:      tenant_settings.NewTenantSettingsService(log, neo4jRepositories, events),
		UserService:                user.NewUserService(neo4jRepositories, postgresRepositories, events),
		WorkflowService:            workflow.NewWorkflowService(postgresRepositories),
		WorkspaceService:           workspace.NewWorkspaceService(neo4jRepositories),
	}

	// add services where service dependecies have been added above
	common.VerifyService = verify.NewVerifyService(
		log,
		postgresRepositories,
		cfg,
		common.EnrichmentService,
	)

	common.DomainService = domain.NewDomainService(
		log,
		common.Cache,
		postgresRepositories,
		neo4jRepositories,
		events,
	)

	common.FileService = files.NewFileService(
		log,
		&cfg.Internal.FileStoreConfig,
		neo4jRepositories,
		common.AttachmentService,
	)

	common.ReminderService = reminders.NewReminderService(
		neo4jRepositories,
		common.NovuService,
	)

	// For all service with other service dependencies, create services and defer initialization
	email := email.NewEmailService(
		neo4jRepositories,
		events,
		nil, // contact
		nil, // org
	)

	jobrole := jobrole.NewJobRoleService(
		neo4jRepositories,
		events,
	)

	issue := issue.NewIssueService(
		log,
		neo4jRepositories,
		events,
		nil, // org
	)

	contact := contact.NewContactService(
		log,
		neo4jRepositories,
		events,
		common.DomainService,
		nil, // email
		nil, // org
		nil, // jobrole
		nil, // social
		nil, // flow
	)

	social := social.NewSocialService(
		log,
		neo4jRepositories,
		events,
		nil, // contact
	)

	org := organization.NewOrganizationService(
		log,
		postgresRepositories,
		neo4jRepositories,
		events,
		common.DomainService,
		industry.NewIndustryService(log, neo4jRepositories),
		nil, // social
		user.NewUserService(neo4jRepositories, postgresRepositories, events),
	)

	contract := contract.NewContractService(
		log,
		neo4jRepositories,
		events,
		grpcClients,
		nil, // opportunity
		nil, // org
	)

	opportunity := opportunity.NewOpportunityService(
		log,
		grpcClients,
		neo4jRepositories,
		events,
		nil, // contract
		nil, // organization
		common.TenantSettingsService,
	)

	sliService := sli.NewServiceLineItemService(
		log,
		events,
		neo4jRepositories,
		nil, // contract
	)

	invoice := invoice.NewInvoiceService(
		log,
		grpcClients,
		neo4jRepositories,
		nil, // contract
		nil, // sli
		common.TenantSettingsService,
	)

	logEntry := logentry.NewLogEntryService(
		log,
		neo4jRepositories,
		events,
		nil, // org
	)

	email.SetContactService(contact)
	email.SetOrganizationService(org)
	email.SetDomainService(common.DomainService)
	if !email.IsInitialized() {
		log.Fatalf("Common Email Service not initialized")
	}
	common.EmailService = email

	mailbox := mailbox.NewMailboxService(
		log,
		postgresRepositories,
		neo4jRepositories,
		common.EmailService,
	)
	common.MailboxService = mailbox

	mailstack := mailstack.NewMailstackService(
		&cfg.External.StripeConfig,
		events,
		postgresRepositories,
		common.CloudflareService,
		common.NamecheapService,
		mailbox,
		common.OpenSRSService,
	)
	common.MailstackService = mailstack

	interactionEvent := interaction_event.NewInteractionEventService(
		neo4jRepositories,
		common.EmailService,
	)

	common.InteractionEventService = interactionEvent

	mail := mail.NewMailService(
		common.Cache,
		postgresRepositories,
		neo4jRepositories,
		common.AzureService,
		nil, // contact
		common.EmailService,
		common.GoogleService,
		interactionEvent,
		common.InteractionSessionService,
		common.OpenSRSService,
		nil, // org
	)

	flowExecution := flow_execution.NewFlowExecutionService(
		neo4jRepositories,
		postgresRepositories,
		events,
		nil, // email
		nil, // flow
		nil, // org
		nil, // social
	)

	flow := flow.NewFlowService(
		neo4jRepositories,
		events,
		nil, // flow execution
	)

	registration := registration.NewRegistrationService(
		events,
		postgresRepositories,
		neo4jRepositories,
		nil, // contact
		common.EmailService,
		nil, // flow
		mailbox,
		nil, // org
		common.PostmarkService,
		common.UserService,
	)

	action := action.NewActionService(
		log,
		neo4jRepositories,
		events,
		nil, // org
	)

	visitorIDAgent := agent.NewAgentVisitorIDService(
		postgresRepositories,
		nil, // action
		common.EnrichmentService,
		nil, // organization
		common.NotificationService,
		common.WorkspaceService,
	)

	location := location.NewLocationService(
		log,
		neo4jRepositories,
		postgresRepositories,
		events,
		&cfg.External.AnthropicPrompts,
		common.AIService,
		nil, // contact
		nil, // org
	)

	// resolve dependencies

	jobrole.SetOrganizationService(org)
	if !jobrole.IsInitialized() {
		log.Fatalf("Common JobRole Service not initialized")
	}
	common.JobRoleService = jobrole

	contact.SetEmailService(email)
	contact.SetOrganizationService(org)
	contact.SetJobRoleService(jobrole)
	contact.SetSocialService(social)
	contact.SetFlowService(flow)
	if !contact.IsInitialized() {
		log.Fatalf("Common Contact Service not initialized")
	}
	common.ContactService = contact

	social.SetContactService(contact)
	if !social.IsInitialized() {
		log.Fatalf("Common Social Service not initialized")
	}
	common.SocialService = social

	org.SetSocialService(social)
	if !org.IsInitialized() {
		log.Fatalf("Common Organization Service not initialized")
	}
	common.OrganizationService = org

	contract.SetOpportunityService(opportunity)
	contract.SetOrganizationService(org)
	if !contract.IsInitialized() {
		log.Fatalf("Common Contract Service not initialized")
	}
	common.ContractService = contract

	opportunity.SetContractService(contract)
	opportunity.SetOrganizationService(org)
	if !opportunity.IsInitialized() {
		log.Fatalf("Common Opportunity Service not initialized")
	}
	common.OpportunityService = opportunity

	sliService.SetContractService(contract)
	if !sliService.IsInitialized() {
		log.Fatalf("Common ServiceLineItem Service not initialized")
	}
	common.ServiceLineItemService = sliService

	invoice.SetContractService(contract)
	invoice.SetServiceLineItemService(sliService)
	if !invoice.IsInitialized() {
		log.Fatalf("Common Invoice Service not initialized")
	}
	common.InvoiceService = invoice

	mail.SetContactService(contact)
	mail.SetOrganizationService(org)
	if !mail.IsInitialized() {
		log.Fatalf("Mail Service not initialized")
	}
	common.MailService = mail

	flowExecution.SetEmailService(email)
	flowExecution.SetFlowService(flow)
	flowExecution.SetOrganizationService(org)
	flowExecution.SetSocialService(social)
	if !flowExecution.IsInitialized() {
		log.Fatalf("Flow Execution Service not initialized")
	}
	common.FlowExecutionService = flowExecution

	flow.SetFlowExecutionService(flowExecution)
	if !flow.IsInitialized() {
		log.Fatalf("Flow Service not initialized")
	}
	common.FlowService = flow

	registration.SetContactService(contact)
	registration.SetFlowService(flow)
	registration.SetOrganizationService(org)
	if !registration.IsInitialized() {
		log.Fatalf("Registration Service not initialized")
	}
	common.RegistrationService = registration

	action.SetOrganizationService(org)
	if !action.IsInitialized() {
		log.Fatalf("Action Service not initialized")
	}
	common.ActionService = action

	visitorIDAgent.SetActionService(action)
	visitorIDAgent.SetOrganizationService(org)
	if !visitorIDAgent.IsInitialized() {
		log.Fatalf("Visitor ID Agent Service not initialized")
	}
	common.AgentService = visitorIDAgent

	location.SetContactService(contact)
	location.SetOrganizationService(org)
	if !location.IsInitialized() {
		log.Fatalf("Location Service not initialized")
	}
	common.LocationService = location

	logEntry.SetOrganizationService(org)
	if !logEntry.IsInitialized() {
		log.Fatalf("Log Entry Service not initialized")
	}
	common.LogEntryService = logEntry

	issue.SetOrganizationService(org)
	if !issue.IsInitialized() {
		log.Fatalf("Issue Service not initialized")
	}
	common.IssueService = issue

	return &common
}
