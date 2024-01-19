package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
)

type invoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
	log                   logger.Logger
	invoiceRequestHandler invoice.InvoiceRequestHandler
	aggregateStore        eventstore.AggregateStore
}

func NewInvoiceService(log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config) *invoiceService {
	return &invoiceService{
		log:                   log,
		invoiceRequestHandler: invoice.NewInvoiceRequestHandler(log, aggregateStore, cfg.Utils),
		aggregateStore:        aggregateStore,
	}
}

func (s *invoiceService) NewOnCycleInvoiceForContract(ctx context.Context, request *invoicepb.NewOnCycleInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.NewOnCycleInvoiceForContract")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	invoiceId := uuid.New().String()

	if _, err := s.invoiceRequestHandler.Handle(ctx, request.Tenant, invoiceId, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewOnCycleInvoiceForContract) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: invoiceId}, nil
}

func (s *invoiceService) FillInvoice(ctx context.Context, request *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.FillInvoiceRequest")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	if _, err := s.invoiceRequestHandler.HandleWithRetry(ctx, request.Tenant, request.InvoiceId, true, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(FillInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) PdfGeneratedInvoice(ctx context.Context, request *invoicepb.PdfGeneratedInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.PdfGeneratedInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	if _, err := s.invoiceRequestHandler.HandleWithRetry(ctx, request.Tenant, request.InvoiceId, true, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PdfGeneratedInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) PayInvoice(ctx context.Context, request *invoicepb.PayInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.PayInvoiceRequest")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	if _, err := s.invoiceRequestHandler.HandleWithRetry(ctx, request.Tenant, request.InvoiceId, true, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PayInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) SimulateInvoice(ctx context.Context, request *invoicepb.SimulateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.NewInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	invoiceId := uuid.New().String()

	if _, err := s.invoiceRequestHandler.HandleTemp(ctx, request.Tenant, invoiceId, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: invoiceId}, nil
}
