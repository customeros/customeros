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
	OrganizationService    *organizationService
	LocationService        *locationService
	CommentService         *commentService
	OpportunityService     *opportunityService
	ServiceLineItemService *serviceLineItemService
	InvoiceService         *invoiceService
	EventStoreService      *eventStoreService
}

func InitServices(cfg *config.Config, repositories *repository.Repositories, aggregateStore eventstore.AggregateStore, commandHandlers *command.CommandHandlers, log logger.Logger, ebs *eventbuffer.EventBufferStoreService) *Services {
	services := Services{}

	services.es = aggregateStore

	services.FileStoreApiService = fsc.NewFileStoreApiService(&cfg.Services.FileStoreApiConfig)
	services.CommonServices = commonService.InitServices(&commonConfig.GlobalConfig{}, repositories.Drivers.PostgresDB, repositories.Drivers.Neo4jDriver, cfg.Neo4j.Database, nil, log)

	services.RequestHandler = NewRequestHandler(log, aggregateStore, cfg.Utils)

	//GRPC services
	services.OrganizationService = NewOrganizationService(log, commandHandlers.Organization, aggregateStore, cfg, &services)
	services.LocationService = NewLocationService(log, commandHandlers.Location)
	services.CommentService = NewCommentService(&services, log, aggregateStore, cfg)
	services.OpportunityService = NewOpportunityService(log, commandHandlers.Opportunity, aggregateStore, &services)
	services.ServiceLineItemService = NewServiceLineItemService(log, aggregateStore, &services)
	services.InvoiceService = NewInvoiceService(repositories, &services, log, aggregateStore)
	services.EventStoreService = NewEventStoreService(&services, log, aggregateStore)

	services.EventStoreGenericService = genericServices.NewEventStoreGenericService(log, aggregateStore)

	return &services
}
