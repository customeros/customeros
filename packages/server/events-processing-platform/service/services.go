package service

import (
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	genericServices "github.com/openline-ai/openline-customer-os/packages/server/events/services"
)

type Services struct {
	es eventstore.AggregateStore

	FileStoreApiService fsc.FileStoreApiService
	CommonServices      *commonService.Services

	EventStoreGenericService genericServices.EventStoreGenericService
	RequestHandler           *requestHandler // generic grpc request handler

	//GRPC services
	ContactService          *contactService
	OrganizationService     *organizationService
	PhoneNumberService      *phoneNumberService
	EmailService            *emailService
	UserService             *userService
	LocationService         *locationService
	JobRoleService          *jobRoleService
	LogEntryService         *logEntryService
	IssueService            *issueService
	CommentService          *commentService
	OpportunityService      *opportunityService
	ContractService         *contractService
	ServiceLineItemService  *serviceLineItemService
	MasterPlanService       *masterPlanService
	OrganizationPlanService *organizationPlanService
	InvoiceService          *invoiceService
	TenantService           *tenantService
	CountryService          *countryService
	EventStoreService       *eventStoreService
	EventCompletionService  *eventCompletionService
}

func InitServices(cfg *config.Config, repositories *repository.Repositories, aggregateStore eventstore.AggregateStore, commandHandlers *command.CommandHandlers, log logger.Logger, ebs *eventbuffer.EventBufferStoreService) *Services {
	services := Services{}

	services.es = aggregateStore

	services.FileStoreApiService = fsc.NewFileStoreApiService(&cfg.Services.FileStoreApiConfig)
	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{}, repositories.Drivers.GormDb, repositories.Drivers.Neo4jDriver, cfg.Neo4j.Database, nil, log)

	services.RequestHandler = NewRequestHandler(log, aggregateStore, cfg.Utils)

	//GRPC services
	services.ContactService = NewContactService(log, &services)
	services.OrganizationService = NewOrganizationService(log, commandHandlers.Organization, aggregateStore, cfg, &services)
	services.PhoneNumberService = NewPhoneNumberService(log, repositories.Neo4jRepositories, commandHandlers.PhoneNumber, &services)
	services.EmailService = NewEmailService(log, repositories.Neo4jRepositories, &services)
	services.UserService = NewUserService(log, aggregateStore, cfg, commandHandlers.User, &services)
	services.LocationService = NewLocationService(log, commandHandlers.Location)
	services.JobRoleService = NewJobRoleService(log, commandHandlers.JobRole)
	services.LogEntryService = NewLogEntryService(log, commandHandlers.LogEntry)
	services.IssueService = NewIssueService(log, commandHandlers.Issue)
	services.CommentService = NewCommentService(&services, log, aggregateStore, cfg)
	services.OpportunityService = NewOpportunityService(log, commandHandlers.Opportunity, aggregateStore, &services)
	services.ContractService = NewContractService(log, aggregateStore, &services)
	services.ServiceLineItemService = NewServiceLineItemService(log, aggregateStore, &services)
	services.MasterPlanService = NewMasterPlanService(log, commandHandlers.MasterPlan, aggregateStore)
	services.OrganizationPlanService = NewOrganizationPlanService(log, commandHandlers.OrganizationPlan, aggregateStore)
	services.InvoiceService = NewInvoiceService(repositories, &services, log, aggregateStore)
	services.TenantService = NewTenantService(&services, log, aggregateStore, cfg)
	services.CountryService = NewCountryService(&services, log, aggregateStore, cfg)
	services.EventStoreService = NewEventStoreService(&services, log, aggregateStore)
	services.EventCompletionService = NewEventCompletionService(&services, log, aggregateStore, cfg)

	services.EventStoreGenericService = genericServices.NewEventStoreGenericService(log, aggregateStore)

	return &services
}
