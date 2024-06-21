package service

import (
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

type Services struct {
	FileStoreApiService fsc.FileStoreApiService
	CommonServices      *commonService.Services

	RequestHandler *requestHandler // generic grpc request handler

	//GRPC services
	ContactService            *contactService
	OrganizationService       *organizationService
	PhoneNumberService        *phoneNumberService
	EmailService              *emailService
	UserService               *userService
	LocationService           *locationService
	JobRoleService            *jobRoleService
	InteractionEventService   *interactionEventService
	InteractionSessionService *interactionSessionService
	LogEntryService           *logEntryService
	IssueService              *issueService
	CommentService            *commentService
	OpportunityService        *opportunityService
	ContractService           *contractService
	ServiceLineItemService    *serviceLineItemService
	MasterPlanService         *masterPlanService
	OrganizationPlanService   *organizationPlanService
	InvoicingCycleService     *invoicingCycleService
	InvoiceService            *invoiceService
	TenantService             *tenantService
	CountryService            *countryService
	ReminderService           *reminderService
	OrderService              *orderService
	EventStoreService         *eventStoreService
}

func InitServices(cfg *config.Config, repositories *repository.Repositories, aggregateStore eventstore.AggregateStore, commandHandlers *command.CommandHandlers, log logger.Logger, ebs *eventstore.EventBufferService) *Services {
	services := Services{}

	services.FileStoreApiService = fsc.NewFileStoreApiService(&cfg.Services.FileStoreApiConfig)
	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{}, repositories.Drivers.GormDb, repositories.Drivers.Neo4jDriver, cfg.Neo4j.Database, nil)

	services.RequestHandler = NewRequestHandler(log, aggregateStore, cfg.Utils)

	//GRPC services
	services.ContactService = NewContactService(log, commandHandlers.Contact, aggregateStore, cfg, &services)
	services.OrganizationService = NewOrganizationService(log, commandHandlers.Organization, aggregateStore, cfg, &services)
	services.PhoneNumberService = NewPhoneNumberService(log, repositories.Neo4jRepositories, commandHandlers.PhoneNumber, &services)
	services.EmailService = NewEmailService(log, repositories.Neo4jRepositories, commandHandlers.Email, &services)
	services.UserService = NewUserService(log, aggregateStore, cfg, commandHandlers.User)
	services.LocationService = NewLocationService(log, commandHandlers.Location)
	services.JobRoleService = NewJobRoleService(log, commandHandlers.JobRole)
	services.InteractionEventService = NewInteractionEventService(log, commandHandlers.InteractionEvent)
	services.InteractionSessionService = NewInteractionSessionService(log, commandHandlers.InteractionSession)
	services.LogEntryService = NewLogEntryService(log, commandHandlers.LogEntry)
	services.IssueService = NewIssueService(log, commandHandlers.Issue)
	services.CommentService = NewCommentService(&services, log, aggregateStore, cfg)
	services.OpportunityService = NewOpportunityService(log, commandHandlers.Opportunity, aggregateStore, &services)
	services.ContractService = NewContractService(log, aggregateStore, &services)
	services.ServiceLineItemService = NewServiceLineItemService(log, aggregateStore, &services)
	services.MasterPlanService = NewMasterPlanService(log, commandHandlers.MasterPlan, aggregateStore)
	services.OrganizationPlanService = NewOrganizationPlanService(log, commandHandlers.OrganizationPlan, aggregateStore)
	services.InvoicingCycleService = NewInvoicingCycleService(log, commandHandlers.InvoicingCycle, aggregateStore)
	services.InvoiceService = NewInvoiceService(repositories, &services, log, aggregateStore)
	services.TenantService = NewTenantService(&services, log, aggregateStore, cfg)
	services.CountryService = NewCountryService(&services, log, aggregateStore, cfg)
	services.ReminderService = NewReminderService(log, aggregateStore, cfg, ebs)
	services.OrderService = NewOrderService(log, aggregateStore, cfg)
	services.EventStoreService = NewEventStoreService(&services, log, aggregateStore)

	return &services
}
