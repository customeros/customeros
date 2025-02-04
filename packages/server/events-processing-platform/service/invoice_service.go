package service

import (
	"context"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/invoice"
	grpcerr "github.com/customeros/customeros/packages/server/events-processing-platform/grpc_errors"
	"github.com/customeros/customeros/packages/server/events-processing-platform/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/customeros/customeros/packages/server/events/eventstore"
)

type invoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
	log                   logger.Logger
	aggregateStore        eventstore.AggregateStore
	requestHandlerService RequestHandler
}

func NewInvoiceService(log logger.Logger, aggregateStore eventstore.AggregateStore, req RequestHandler) *invoiceService {
	return &invoiceService{
		log:                   log,
		aggregateStore:        aggregateStore,
		requestHandlerService: req,
	}
}

func (s *invoiceService) PayInvoiceNotification(ctx context.Context, request *invoicepb.PayInvoiceNotificationRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.PayInvoiceNotification")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.requestHandlerService.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PayInvoiceNotification) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) PermanentlyDeleteInitializedInvoice(ctx context.Context, request *invoicepb.PermanentlyDeleteInitializedInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.PermanentlyDeleteInitializedInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.requestHandlerService.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PermanentlyDeleteInitializedInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) VoidInvoice(ctx context.Context, request *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.VoidInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.requestHandlerService.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(VoidInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) RemindInvoiceNotification(ctx context.Context, request *invoicepb.RemindInvoiceNotificationRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.RemindInvoiceNotification")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.requestHandlerService.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemindInvoiceNotification) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}
