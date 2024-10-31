package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"gorm.io/gorm"
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
	PhoneNumberService         PhoneNumberService
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
	TenantService              TenantService
	ExternalSystemService      ExternalSystemService
	ActionService              ActionService
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
	ReminderService            ReminderService
	CloudflareService          CloudflareService
	EnrichmentService          EnrichmentService
	NamecheapService           NamecheapService
	OpensrsService             OpensrsService
}

func InitServices(log logger.Logger, driver *neo4j.DriverWithContext, cfg *config.Config, commonServices *commonService.Services, grpcClients *grpc_client.Clients, gormDb *gorm.DB) *Services {
	repositories := repository.InitRepos(driver, cfg.Neo4j.Database, gormDb)

	services := Services{
		CommonServices:             commonServices,
		BankAccountService:         NewBankAccountService(log, repositories, grpcClients),
		CustomFieldService:         NewCustomFieldService(log, repositories),
		CustomFieldTemplateService: NewCustomFieldTemplateService(log, repositories),
		LocationService:            NewLocationService(log, repositories),
		PageViewService:            NewPageViewService(log, repositories),
		TenantService:              NewTenantService(log, repositories, grpcClients),
		ExternalSystemService:      NewExternalSystemService(log, repositories),
		ActionService:              NewActionService(log, repositories),
		CountryService:             NewCountryService(log, repositories),
		ActionItemService:          NewActionItemService(log, repositories),
		BillableService:            NewBillableService(log, repositories),
		LogEntryService:            NewLogEntryService(log, repositories),
		CommentService:             NewCommentService(log, repositories),
		ReminderService:            NewReminderService(log, repositories, grpcClients),
	}
	services.Repositories = repositories
	services.UserService = NewUserService(log, repositories, grpcClients, &services)
	services.OrganizationService = NewOrganizationService(log, repositories, grpcClients, &services)
	services.IssueService = NewIssueService(log, repositories, &services)
	services.PhoneNumberService = NewPhoneNumberService(log, repositories, grpcClients, &services)
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
	services.FileStoreApiService = fsc.NewFileStoreApiService(&cfg.InternalServices.FileStoreApiConfig)
	services.CloudflareService = NewCloudflareService(log, &services, cfg)
	services.EnrichmentService = NewEnrichmentService(log, &services, cfg)
	services.OpensrsService = NewOpensrsService(log, &services, cfg)
	services.NamecheapService = NewNamecheapService(log, cfg, repositories)

	log.Info("Init cache service")
	services.Cache = NewCacheService(&services)
	services.Cache.InitCache()
	log.Info("Init cache service done")

	services.Cfg = cfg
	services.Log = log
	return &services
}
