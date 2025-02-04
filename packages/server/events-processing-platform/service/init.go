package service

import (
	"github.com/customeros/customeros/packages/server/events/eventstore"

	"github.com/customeros/customeros/packages/server/events-processing-platform/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform/logger"
)

type Services struct {
	es eventstore.AggregateStore

	RequestHandler *requestHandler // generic grpc request handler

	// GRPC services
	InvoiceService *invoiceService
}

func InitServices(
	cfg *config.Config,
	aggregateStore eventstore.AggregateStore,
	log logger.Logger,
) *Services {
	services := Services{}

	services.es = aggregateStore
	services.RequestHandler = NewRequestHandler(log, aggregateStore, &cfg.Utils)
	services.InvoiceService = NewInvoiceService(
		log,
		aggregateStore,
		services.RequestHandler,
	)

	return &services
}
