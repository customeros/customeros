package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	repository "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type invoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
	log                   logger.Logger
	invoiceRequestHandler invoice.InvoiceRequestHandler
	aggregateStore        eventstore.AggregateStore
	invoiceRepository     repository.InvoiceRepository
}

func NewInvoiceService(log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config, invoiceRepository repository.InvoiceRepository) *invoiceService {
	return &invoiceService{
		log:                   log,
		invoiceRequestHandler: invoice.NewInvoiceRequestHandler(log, aggregateStore, cfg.Utils),
		aggregateStore:        aggregateStore,
		invoiceRepository:     invoiceRepository,
	}
}

func (s *invoiceService) NewInvoiceForContract(ctx context.Context, request *invoicepb.NewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.NewInvoiceForContract")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Currency == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("currency"))
	} else if request.InvoicePeriodStart == nil {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoicePeriodStart"))
	} else if request.InvoicePeriodEnd == nil {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoicePeriodEnd"))
	} else if request.ContractId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("contractId"))
	}

	// Check if the contract aggregate exists
	contractExists, err := s.checkContractExists(ctx, request.Tenant, request.ContractId)
	if err != nil {
		s.log.Error(err, "error checking contract existence")
		return nil, status.Errorf(codes.Internal, "error checking contract existence: %v", err)
	}
	if !contractExists {
		return nil, status.Errorf(codes.NotFound, "contract with ID %s not found", request.ContractId)
	}

	invoiceId := uuid.New().String()

	extraParams := map[string]any{
		invoice.PARAM_INVOICE_NUMBER: s.prepareInvoiceNumber(request.Tenant),
	}

	if _, err := s.invoiceRequestHandler.Handle(ctx, request.Tenant, invoiceId, request, extraParams); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewOnCycleInvoiceForContract) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: invoiceId}, nil
}

func (s *invoiceService) prepareInvoiceNumber(tenant string) string {
	maxAttempts := 20
	var invoiceNumber string
	for attempt := 1; attempt < maxAttempts+1; attempt++ {
		invoiceNumber = generateNewRandomInvoiceNumber()
		invoiceNumberEntity := postgresentity.InvoiceNumberEntity{
			InvoiceNumber: invoiceNumber,
			Tenant:        tenant,
			Attempts:      attempt,
		}
		innerErr := s.invoiceRepository.Reserve(invoiceNumberEntity)
		if innerErr == nil {
			break
		}
	}

	return invoiceNumber
}

func generateNewRandomInvoiceNumber() string {
	digits := "0123456789"
	consonants := "BCDFGHJKLMNPQRSTVWXYZ"
	invoiceNumber := utils.GenerateRandomStringFromCharset(3, consonants) + "-" + utils.GenerateRandomStringFromCharset(5, digits)
	return invoiceNumber
}

func (s *invoiceService) FillInvoice(ctx context.Context, request *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.FillInvoice")
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

func (s *invoiceService) GenerateInvoicePdf(ctx context.Context, request *invoicepb.GenerateInvoicePdfRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.GenerateInvoicePdf")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	if _, err := s.invoiceRequestHandler.HandleWithRetry(ctx, request.Tenant, request.InvoiceId, true, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(GenerateInvoicePdf) tenant:{%v}, err: %v", request.Tenant, err.Error())
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

func (s *invoiceService) UpdateInvoice(ctx context.Context, request *invoicepb.UpdateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.UpdateInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	if _, err := s.invoiceRequestHandler.HandleWithRetry(ctx, request.Tenant, request.InvoiceId, true, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
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

func (s *invoiceService) PayInvoiceNotification(ctx context.Context, request *invoicepb.PayInvoiceNotificationRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.PayInvoiceNotification")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	if _, err := s.invoiceRequestHandler.HandleWithRetry(ctx, request.Tenant, request.InvoiceId, true, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PayInvoiceNotification) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) checkContractExists(ctx context.Context, tenant, contractId string) (bool, error) {
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	err := s.aggregateStore.Exists(ctx, contractAggregate.GetID())
	if err != nil {
		if errors.Is(err, eventstore.ErrAggregateNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil // The contract exists
}

func (s *invoiceService) RequestFillInvoice(ctx context.Context, request *invoicepb.RequestFillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.RequestFillInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	if _, err := s.invoiceRequestHandler.HandleWithRetry(ctx, request.Tenant, request.InvoiceId, true, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RequestFillInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}
