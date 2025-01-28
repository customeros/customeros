package service

import (
	"github.com/customeros/customeros/packages/server/events/eventstore"

	"github.com/customeros/customeros/packages/server/events-processing-platform/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform/repository"
)

type Services struct {
	es eventstore.AggregateStore

	RequestHandler *requestHandler // generic grpc request handler

	// GRPC services
	OrganizationService *organizationService
	InvoiceService      *invoiceService
}

func InitServices(
	cfg *config.Config,
	repositories *repository.Repositories,
	aggregateStore eventstore.AggregateStore,
	log logger.Logger,
) *Services {
	services := Services{}

	services.es = aggregateStore
	services.RequestHandler = NewRequestHandler(log, aggregateStore, &cfg.Utils)

	services.OrganizationService = NewOrganizationService(log, aggregateStore, cfg, services.RequestHandler)

	services.InvoiceService = NewInvoiceService(
		repositories,
		log,
		aggregateStore,
		services.RequestHandler,
	)

	return &services
}
