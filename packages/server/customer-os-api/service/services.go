package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type Services struct {
	Cfg          *config.Config
	Log          logger.Logger
	Cache        CacheService // todo move this to cache
	Repositories *repository.Repositories

	CommonServices      *commonService.Services
	FileStoreApiService fsc.FileStoreApiService

	BankAccountService         BankAccountService
	ContactService             ContactService
	OrganizationService        OrganizationService
	CustomFieldService         CustomFieldService
	EmailService               EmailService
	UserService                UserService
	CustomFieldTemplateService CustomFieldTemplateService
	TimelineEventService       TimelineEventService
	NoteService                NoteService
	CalendarService            CalendarService
	LocationService            LocationService
	SearchService              SearchService
	QueryService               DashboardService
	IssueService               IssueService
	PageViewService            PageViewService
	MeetingService             MeetingService
	ExternalSystemService      ExternalSystemService
	CountryService             CountryService
	ActionItemService          ActionItemService
	BillableService            BillableService
	LogEntryService            LogEntryService
	CommentService             CommentService
	ContractService            ContractService
	ServiceLineItemService     ServiceLineItemService
	OpportunityService         OpportunityService
	BillingProfileService      BillingProfileService
	InvoiceService             InvoiceService
	SlackService               SlackService
	EnrichmentService          EnrichmentService
	WebhookService             WebhookService
}

func InitServices(log logger.Logger, driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB, cfg *config.Config, commonServices *commonService.Services, grpcClients *grpc_client.Clients) *Services {
	repositories := repository.InitRepos(driver, cfg.Database.Neo4j.Database, postgresDB)

	services := Services{
		CommonServices:             commonServices,
		CustomFieldService:         NewCustomFieldService(log, repositories),
		CustomFieldTemplateService: NewCustomFieldTemplateService(log, repositories),
		LocationService:            NewLocationService(log, repositories),
		PageViewService:            NewPageViewService(log, repositories),
		ExternalSystemService:      NewExternalSystemService(log, repositories),
		CountryService:             NewCountryService(log, repositories),
		ActionItemService:          NewActionItemService(log, repositories),
		BillableService:            NewBillableService(log, repositories),
		LogEntryService:            NewLogEntryService(log, repositories),
		CommentService:             NewCommentService(log, repositories),
	}
	services.Repositories = repositories
	services.BankAccountService = NewBankAccountService(log, repositories, grpcClients, &services)
	services.UserService = NewUserService(log, repositories, grpcClients, &services)
	services.OrganizationService = NewOrganizationService(log, repositories, grpcClients, &services)
	services.IssueService = NewIssueService(log, repositories, &services)
	services.CalendarService = NewCalendarService(log, repositories, &services)
	services.EmailService = NewEmailService(log, repositories, &services, grpcClients)
	services.ContactService = NewContactService(log, repositories, grpcClients, &services)
	services.NoteService = NewNoteService(log, repositories, &services)
	services.TimelineEventService = NewTimelineEventService(log, repositories, &services)
	services.SearchService = NewSearchService(log, repositories, &services)
	services.QueryService = NewDashboardService(log, repositories, &services)
	services.MeetingService = NewMeetingService(log, repositories, &services)
	services.ContractService = NewContractService(log, repositories, grpcClients, &services)
	services.ServiceLineItemService = NewServiceLineItemService(log, repositories, grpcClients, &services)
	services.OpportunityService = NewOpportunityService(log, repositories, grpcClients, &services)
	services.BillingProfileService = NewBillingProfileService(log, repositories, grpcClients)
	services.InvoiceService = NewInvoiceService(log, repositories, grpcClients, &services)
	services.SlackService = NewSlackService(log, repositories, grpcClients, &services)
	services.FileStoreApiService = fsc.NewFileStoreApiService(&cfg.InternalServices.FileStoreApi)
	services.EnrichmentService = NewEnrichmentService(log, &services, cfg)
	services.WebhookService = NewWebhookService(log, repositories, &services)

	log.Info("Init cache service")
	services.Cache = NewCacheService(&services)
	services.Cache.InitCache()
	log.Info("Init cache service done")

	services.Cfg = cfg
	services.Log = log
	return &services
}
