package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicingcyclepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoicing_cycle"
	"golang.org/x/net/context"
)

type invoicingCycleService struct {
	invoicingcyclepb.UnimplementedInvoicingCycleServiceServer
	log                           logger.Logger
	invoicingCycleCommandHandlers *command_handler.CommandHandlers
	aggregateStore                eventstore.AggregateStore
}

func NewInvoicingCycleService(log logger.Logger, commandHandlers *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore) *invoicingCycleService {
	return &invoicingCycleService{
		log:                           log,
		invoicingCycleCommandHandlers: commandHandlers,
		aggregateStore:                aggregateStore,
	}
}

func (s *invoicingCycleService) CreateInvoicingCycleType(ctx context.Context, request *invoicingcyclepb.CreateInvoicingCycleTypeRequest) (*invoicingcyclepb.InvoicingCycleTypeResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoicingCycleService.CreateInvoicingCycleType")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	invoicingCycleTypeId := uuid.New().String()

	createdAt := utils.TimestampProtoToTimePtr(request.CreatedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createInvoicingCycleTypeCommand := command.NewCreateInvoicingCycleTypeCommand(
		invoicingCycleTypeId,
		request.Tenant,
		request.LoggedInUserId,
		sourceFields,
		createdAt,
		model.InvoicingCycleType(request.Type),
	)

	if err := s.invoicingCycleCommandHandlers.CreateInvoicingCycle.Handle(ctx, createInvoicingCycleTypeCommand); err != nil {
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

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	cmd := command.NewUpdateInvoicingCycleCommand(
		request.InvoicingCycleTypeId,
		request.Tenant,
		request.LoggedInUserId,
		sourceFields,
		updatedAt,
		model.InvoicingCycleType(request.Type),
	)

	if err := s.invoicingCycleCommandHandlers.UpdateInvoicingCycle.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateInvoicingCycle.Handle) tenant:{%v}, invoicingCycleId:{%v}, err: %v", request.Tenant, request.InvoicingCycleTypeId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicingcyclepb.InvoicingCycleTypeResponse{Id: request.InvoicingCycleTypeId}, nil
}
