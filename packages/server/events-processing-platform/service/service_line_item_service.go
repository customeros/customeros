package service

import (
	"github.com/google/uuid"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/service_line_item"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serviceLineItemService struct {
	servicelineitempb.UnimplementedServiceLineItemGrpcServiceServer
	log                            logger.Logger
	serviceLineItemCommandHandlers *command_handler.CommandHandlers
	aggregateStore                 eventstore.AggregateStore
}

func NewServiceLineItemService(log logger.Logger, commandHandlers *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore) *serviceLineItemService {
	return &serviceLineItemService{
		log:                            log,
		serviceLineItemCommandHandlers: commandHandlers,
		aggregateStore:                 aggregateStore,
	}
}

func (s *serviceLineItemService) CreateServiceLineItem(ctx context.Context, request *servicelineitempb.CreateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ServiceLineItemService.CreateServiceLineItem")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Validate contract ID
	if request.ContractId == "" {
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

	serviceLineItemId := uuid.New().String()

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	createdAt, updatedAt := convertCreateAndUpdateProtoTimestampsToTime(request.CreatedAt, request.UpdatedAt)

	createServiceLineItemCommand := command.NewCreateServiceLineItemCommand(
		serviceLineItemId,
		request.Tenant,
		request.LoggedInUserId,
		model.ServiceLineItemDataFields{
			Billed:      model.BilledType(request.Billed),
			Licenses:    request.Licenses,
			Price:       request.Price,
			Description: request.Description,
			ContractId:  request.ContractId,
		},
		source,
		createdAt,
		updatedAt,
	)

	if err := s.serviceLineItemCommandHandlers.CreateServiceLineItem.Handle(ctx, createServiceLineItemCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateServiceLineItem.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created service line item
	return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: serviceLineItemId}, nil
}

func (s *serviceLineItemService) checkContractExists(ctx context.Context, tenant, contractId string) (bool, error) {
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
