package cosapi_services

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	cosapi_interfaces "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	api_action_item "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/action_item"
	api_billable "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/billable"
	api_billing_profile "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/billing_profile"
	api_cache "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/cache"
	api_calendar "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/calendar"
	api_comment "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/comment"
	api_contact "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/contact"
	api_contract "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/contract"
	api_country "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/country"
	api_customfields "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/custom_fields"
	api_custom_fields_template "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/custom_fields_template"
	api_dashboard "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/dashboard"
	api_email "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/email"
	api_external_system "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/external_system"
	api_invoice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/invoice"
	api_jwt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/jwt"
	api_location "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/location"
	api_log_entry "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/log_entry"
	api_meeting "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/meeting"
	api_note "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/note"
	api_oauthuser "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/oauth_user"
	api_opportunity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/opportunity"
	api_organization "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/organization"
	api_personal_integrations "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/personal_integrations"
	api_search "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/search"
	api_sli "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/service_line_item"
	api_slack "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/slack"
	api_tenant_settings "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/tenant_settings"
	api_timeline_event "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/timeline_event"
	api_user "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/user"
	api_webhook "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services/webhook"
)

type Services struct {
	CommonServices *commonService.CommonServices
	Cfg            *config.Config
	Log            logger.Logger
	Repositories   *repository.Repositories
	Cache          *caches.Cache

	ActionItemService           cosapi_interfaces.ActionItemService
	BankAccountService          cosapi_interfaces.BankAccountService
	BillableService             cosapi_interfaces.BillableService
	BillingProfileService       cosapi_interfaces.BillingProfileService
	CacheService                cosapi_interfaces.CacheService
	CalendarService             cosapi_interfaces.CalendarService
	CommentService              cosapi_interfaces.CommentService
	ContactService              cosapi_interfaces.ContactService
	ContractService             cosapi_interfaces.ContractService
	CountryService              cosapi_interfaces.CountryService
	CustomFieldService          cosapi_interfaces.CustomFieldService
	CustomFieldsTemplateService cosapi_interfaces.CustomFieldTemplateService
	DashboardService            cosapi_interfaces.DashboardService
	EmailService                cosapi_interfaces.EmailService
	ExternalSystemService       cosapi_interfaces.ExternalSystemService
	InvoiceService              cosapi_interfaces.InvoiceService
	IssueService                cosapi_interfaces.IssueService
	JWTService                  api_jwt.JWTService
	LocationService             cosapi_interfaces.LocationService
	LogEntryService             cosapi_interfaces.LogEntryService
	MeetingService              cosapi_interfaces.MeetingService
	NoteService                 cosapi_interfaces.NoteService
	OAuthUserSettingsService    cosapi_interfaces.OAuthUserSettingsService
	OpportunityService          cosapi_interfaces.OpportunityService
	OrganizationService         cosapi_interfaces.OrganizationService
	PersonalIntegrationsService cosapi_interfaces.PersonalIntegrationsService
	SearchService               cosapi_interfaces.SearchService
	ServiceLineItemService      cosapi_interfaces.ServiceLineItemService
	SlackService                cosapi_interfaces.SlackService
	TenantSettingsService       cosapi_interfaces.TenantSettingsService
	TimelineEventService        cosapi_interfaces.TimelineEventService
	UserService                 cosapi_interfaces.UserService
	WebhookService              cosapi_interfaces.WebhookService
}

func InitServices(log logger.Logger, driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB, cfg *config.Config, grpcClients *grpc_client.Clients) *Services {
	repositories := repository.InitRepos(driver, cfg.CommonServices.Neo4jConfig.Database, postgresDB)

	commonServices := commonService.InitCommonServices(
		log,
		repositories.Neo4jRepositories,
		repositories.PostgresRepositories,
		cfg.CommonServices,
		grpcClients,
	)

	services := Services{
		CommonServices:              commonServices,
		Repositories:                repositories,
		Cache:                       commonServices.Cache,
		Cfg:                         cfg,
		Log:                         log,
		ActionItemService:           api_action_item.NewActionItemService(log, repositories),
		BillableService:             api_billable.NewBillableService(log, repositories),
		BillingProfileService:       api_billing_profile.NewBillingProfileService(log, repositories, grpcClients),
		CacheService:                api_cache.NewCacheService(repositories.Neo4jRepositories),
		CalendarService:             api_calendar.NewCalendarService(log, repositories),
		CommentService:              api_comment.NewCommentService(log, repositories),
		CountryService:              api_country.NewCountryService(log, repositories),
		CustomFieldService:          api_customfields.NewCustomFieldService(log, repositories),
		CustomFieldsTemplateService: api_custom_fields_template.NewCustomFieldTemplateService(log, repositories),
		DashboardService:            api_dashboard.NewDashboardService(log, repositories),
		EmailService:                api_email.NewEmailService(log, repositories, grpcClients),
		ExternalSystemService:       api_external_system.NewExternalSystemService(log, repositories),
		JWTService:                  *api_jwt.NewJWTTenantUserService(&cfg.CommonServices.InternalServices.FileStoreConfig),
		LocationService:             api_location.NewLocationService(log, repositories),
		LogEntryService:             api_log_entry.NewLogEntryService(log, repositories),
		NoteService:                 api_note.NewNoteService(log, repositories),
		OAuthUserSettingsService:    api_oauthuser.NewUserSettingsService(log, repositories.PostgresRepositories),
		PersonalIntegrationsService: api_personal_integrations.NewPersonalIntegrationsService(log, repositories.PostgresRepositories),
		SearchService:               api_search.NewSearchService(log, repositories),
		TenantSettingsService:       api_tenant_settings.NewTenantSettingsService(log, cfg, repositories.PostgresRepositories),
		TimelineEventService:        api_timeline_event.NewTimelineEventService(log, repositories),
		WebhookService:              api_webhook.NewWebhookService(log, repositories),
		UserService:                 api_user.NewUserService(log, repositories, grpcClients),
	}

	// has dependencies
	services.ContractService = api_contract.NewContractService(
		log,
		repositories,
		grpcClients,
		commonServices.TenantSettingsService,
		commonServices.ContractService,
		commonServices.OpportunityService,
	)

	services.InvoiceService = api_invoice.NewInvoiceService(
		log,
		repositories,
		grpcClients,
		commonServices.InvoiceService,
	)

	services.OpportunityService = api_opportunity.NewOpportunityService(
		log,
		repositories,
		grpcClients,
		commonServices.OpportunityService,
		commonServices.OrganizationService,
	)

	services.OrganizationService = api_organization.NewOrganizationService(
		log,
		repositories,
		grpcClients,
		commonServices.Events,
		commonServices.OrganizationService,
	)

	services.ContactService = api_contact.NewContactService(
		log,
		repositories,
		grpcClients,
		commonServices.ContactService,
		commonServices.EmailService,
		commonServices.PhoneNumberService,
		services.OrganizationService,
	)

	services.ServiceLineItemService = api_sli.NewServiceLineItemService(
		log,
		repositories,
		grpcClients,
		commonServices.ServiceLineItemService,
		services.ContractService,
		commonServices.TenantSettingsService,
	)

	services.SlackService = api_slack.NewSlackService(
		log,
		repositories,
		grpcClients,
		commonServices.SlackService,
	)

	services.MeetingService = api_meeting.NewMeetingService(
		log,
		repositories,
		services.OrganizationService,
		services.NoteService,
	)

	return &services
}
