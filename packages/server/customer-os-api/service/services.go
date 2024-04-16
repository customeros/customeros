package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	commonAuthService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"gorm.io/gorm"
)

type Services struct {
	cfg          *config.Config
	Cache        CacheService
	Repositories *repository.Repositories

	CommonServices     *commonService.Services
	CommonAuthServices *commonAuthService.Services

	BankAccountService         BankAccountService
	ContactService             ContactService
	OrganizationService        OrganizationService
	CustomFieldService         CustomFieldService
	PhoneNumberService         PhoneNumberService
	EmailService               EmailService
	UserService                UserService
	FieldSetService            FieldSetService
	EntityTemplateService      EntityTemplateService
	FieldSetTemplateService    FieldSetTemplateService
	CustomFieldTemplateService CustomFieldTemplateService
	TimelineEventService       TimelineEventService
	NoteService                NoteService
	JobRoleService             JobRoleService
	CalendarService            CalendarService
	LocationService            LocationService
	TagService                 TagService
	SearchService              SearchService
	QueryService               DashboardService
	DomainService              DomainService
	IssueService               IssueService
	InteractionSessionService  InteractionSessionService
	InteractionEventService    InteractionEventService
	PageViewService            PageViewService
	AnalysisService            AnalysisService
	AttachmentService          AttachmentService
	MeetingService             MeetingService
	TenantService              TenantService
	WorkspaceService           WorkspaceService
	SocialService              SocialService
	PlayerService              PlayerService
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
	MasterPlanService          MasterPlanService
	BillingProfileService      BillingProfileService
	InvoiceService             InvoiceService
	OrganizationPlanService    OrganizationPlanService
	SlackService               SlackService
	ReminderService            ReminderService
	OrderService               OrderService
	OfferingService            OfferingService
}

func InitServices(log logger.Logger, driver *neo4j.DriverWithContext, cfg *config.Config, commonServices *commonService.Services, commonAuthServices *commonAuthService.Services, grpcClients *grpc_client.Clients, gormDb *gorm.DB) *Services {
	Repositories := repository.InitRepos(driver, cfg.Neo4j.Database, gormDb)

	services := Services{
		CommonServices:             commonServices,
		CommonAuthServices:         commonAuthServices,
		BankAccountService:         NewBankAccountService(log, Repositories, grpcClients),
		OrganizationService:        NewOrganizationService(log, Repositories, grpcClients),
		CustomFieldService:         NewCustomFieldService(log, Repositories),
		UserService:                NewUserService(log, Repositories, grpcClients),
		FieldSetService:            NewFieldSetService(log, Repositories),
		EntityTemplateService:      NewEntityTemplateService(log, Repositories),
		FieldSetTemplateService:    NewFieldSetTemplateService(log, Repositories),
		CustomFieldTemplateService: NewCustomFieldTemplateService(log, Repositories),
		LocationService:            NewLocationService(log, Repositories),
		TagService:                 NewTagService(log, Repositories),
		DomainService:              NewDomainService(log, Repositories),
		PageViewService:            NewPageViewService(log, Repositories),
		AttachmentService:          NewAttachmentService(log, Repositories),
		TenantService:              NewTenantService(log, Repositories, grpcClients),
		WorkspaceService:           NewWorkspaceService(log, Repositories),
		SocialService:              NewSocialService(log, Repositories),
		ExternalSystemService:      NewExternalSystemService(log, Repositories),
		ActionService:              NewActionService(log, Repositories),
		CountryService:             NewCountryService(log, Repositories),
		ActionItemService:          NewActionItemService(log, Repositories),
		BillableService:            NewBillableService(log, Repositories),
		LogEntryService:            NewLogEntryService(log, Repositories),
		CommentService:             NewCommentService(log, Repositories),
		MasterPlanService:          NewMasterPlanService(log, Repositories, grpcClients),
		OrganizationPlanService:    NewOrganizationPlanService(log, Repositories, grpcClients),
		ReminderService:            NewReminderService(log, Repositories, grpcClients),
		OrderService:               NewOrderService(log, Repositories),
		OfferingService:            NewOfferingService(log, Repositories, grpcClients),
	}
	services.IssueService = NewIssueService(log, Repositories, &services)
	services.PhoneNumberService = NewPhoneNumberService(log, Repositories, grpcClients, &services)
	services.JobRoleService = NewJobRoleService(log, Repositories, &services)
	services.CalendarService = NewCalendarService(log, Repositories, &services)
	services.EmailService = NewEmailService(log, Repositories, &services, grpcClients)
	services.ContactService = NewContactService(log, Repositories, grpcClients, &services)
	services.NoteService = NewNoteService(log, Repositories, &services)
	services.TimelineEventService = NewTimelineEventService(log, Repositories, &services)
	services.SearchService = NewSearchService(log, Repositories, &services)
	services.QueryService = NewDashboardService(log, Repositories, &services)
	services.InteractionEventService = NewInteractionEventService(log, Repositories, &services)
	services.InteractionSessionService = NewInteractionSessionService(log, Repositories, &services)
	services.AnalysisService = NewAnalysisService(log, Repositories, &services)
	services.MeetingService = NewMeetingService(log, Repositories, &services)
	services.PlayerService = NewPlayerService(Repositories, &services)
	services.ContractService = NewContractService(log, Repositories, grpcClients, &services)
	services.ServiceLineItemService = NewServiceLineItemService(log, Repositories, grpcClients, &services)
	services.OpportunityService = NewOpportunityService(log, Repositories, grpcClients, &services)
	services.BillingProfileService = NewBillingProfileService(log, Repositories, grpcClients)
	services.InvoiceService = NewInvoiceService(log, Repositories, grpcClients, &services)
	services.SlackService = NewSlackService(log, Repositories, grpcClients, &services)

	log.Info("Init cache service")
	services.Cache = NewCacheService(&services)
	services.Cache.InitCache()
	log.Info("Init cache service done")

	services.cfg = cfg
	return &services
}
