package service

import (
	"context"
	"github.com/google/uuid"
	invoicingcycleEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicingcyclepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoicing_cycle"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type invoicingCycleService struct {
	invoicingcyclepb.UnimplementedInvoicingCycleGrpcServiceServer
	log            logger.Logger
	eventHandlers  *invoicingcycleEvents.EventHandlers
	aggregateStore eventstore.AggregateStore
}

func NewInvoicingCycleService(log logger.Logger, eventHandlers *invoicingcycleEvents.EventHandlers, aggregateStore eventstore.AggregateStore) *invoicingCycleService {
	return &invoicingCycleService{
		log:            log,
		eventHandlers:  eventHandlers,
		aggregateStore: aggregateStore,
	}
}

func (s *invoicingCycleService) CreateInvoicingCycleType(ctx context.Context, request *invoicingcyclepb.CreateInvoicingCycleTypeRequest) (*invoicingcyclepb.InvoicingCycleTypeResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoicingCycleService.CreateInvoicingCycleType")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	invoicingCycleTypeId := uuid.New().String()

	baseRequest := events.NewBaseRequest(invoicingCycleTypeId, request.Tenant, request.LoggedInUserId, events.SourceFromGrpc(request.SourceFields))

	if err := s.eventHandlers.CreateInvoicingCycle.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateInvoicingCycleType.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicingcyclepb.InvoicingCycleTypeResponse{Id: invoicingCycleTypeId}, nil
}

func (s *invoicingCycleService) UpdateInvoicingCycleType(ctx context.Context, request *invoicingcyclepb.UpdateInvoicingCycleTypeRequest) (*invoicingcyclepb.InvoicingCycleTypeResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoicingCycleService.UpdateInvoicingCycleType")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.InvoicingCycleTypeId)

	if request.InvoicingCycleTypeId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoicingCycleTypeId"))
	}

	baseRequest := events.NewBaseRequest(request.InvoicingCycleTypeId, request.Tenant, request.LoggedInUserId, events.SourceFromGrpc(request.SourceFields))

	if err := s.eventHandlers.UpdateInvoicingCycle.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateInvoicingCycle.Handle) tenant:{%v}, invoicingCycleId:{%v}, err: %v", request.Tenant, request.InvoicingCycleTypeId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicingcyclepb.InvoicingCycleTypeResponse{Id: request.InvoicingCycleTypeId}, nil
}
