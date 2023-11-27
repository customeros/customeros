package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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

	createCommand := command.NewCreateServiceLineItemCommand(
		serviceLineItemId,
		request.Tenant,
		request.LoggedInUserId,
		model.ServiceLineItemDataFields{
			Billed:     model.BilledType(request.Billed),
			Quantity:   request.Quantity,
			Price:      float64(request.Price),
			Name:       request.Name,
			ContractId: request.ContractId,
			ParentId:   serviceLineItemId,
		},
		source,
		createdAt,
		updatedAt,
	)
	createCommand.StartedAt = utils.TimestampProtoToTimePtr(request.StartedAt)
	createCommand.EndedAt = utils.TimestampProtoToTimePtr(request.EndedAt)

	if err = s.serviceLineItemCommandHandlers.CreateServiceLineItem.Handle(ctx, createCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateServiceLineItem.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created service line item
	return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: serviceLineItemId}, nil
}

func (s *serviceLineItemService) UpdateServiceLineItem(ctx context.Context, request *servicelineitempb.UpdateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ServiceLineItemService.UpdateServiceLineItem")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Validate service line item ID
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	if request.IsRetroactiveCorrection {
		updateServiceLineItemCommand := command.NewUpdateServiceLineItemCommand(
			request.Id,
			request.Tenant,
			request.LoggedInUserId,
			model.ServiceLineItemDataFields{
				Billed:   model.BilledType(request.Billed),
				Quantity: request.Quantity,
				Price:    float64(request.Price),
				Name:     request.Name,
				Comments: request.Comments,
			},
			source,
			updatedAt,
		)

		if err := s.serviceLineItemCommandHandlers.UpdateServiceLineItem.Handle(ctx, updateServiceLineItemCommand); err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("(UpdateServiceLineItem.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
			return nil, grpcerr.ErrResponse(err)
		}
		// Return the ID of the updated service line item
		return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: request.Id}, nil
	} else {
		versionDate := utils.NowPtr()

		// Validate contract ID
		if request.ContractId == "" {
			return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("contractId"))
		}

		// Check if the contract aggregate exists prior to closing the service line item
		contractExists, err := s.checkContractExists(ctx, request.Tenant, request.ContractId)
		if err != nil {
			s.log.Error(err, "error checking contract existence")
			return nil, status.Errorf(codes.Internal, "error checking contract existence: %v", err)
		}
		if !contractExists {
			return nil, status.Errorf(codes.NotFound, "contract with ID %s not found", request.ContractId)
		}

		//Close service line item
		closeSliCommand := command.NewCloseServiceLineItemCommand(request.Id, request.Tenant, request.LoggedInUserId, source.AppSource,
			versionDate, utils.TimestampProtoToTimePtr(request.UpdatedAt))

		if err := s.serviceLineItemCommandHandlers.CloseServiceLineItem.Handle(ctx, closeSliCommand); err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("(CloseServiceLineItem.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
			return nil, grpcerr.ErrResponse(err)
		}

		//Create new service line item
		serviceLineItemId := uuid.New().String()

		createCommand := command.NewCreateServiceLineItemCommand(
			serviceLineItemId,
			request.Tenant,
			request.LoggedInUserId,
			model.ServiceLineItemDataFields{
				Billed:     model.BilledType(request.Billed),
				Quantity:   request.Quantity,
				Price:      float64(request.Price),
				Name:       request.Name,
				ContractId: request.ContractId,
				ParentId:   request.Id,
			},
			source,
			utils.NowPtr(),
			updatedAt,
		)
		createCommand.StartedAt = versionDate

		if err = s.serviceLineItemCommandHandlers.CreateServiceLineItem.Handle(ctx, createCommand); err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("(CreateServiceLineItem.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
			return nil, grpcerr.ErrResponse(err)
		}
		return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: serviceLineItemId}, nil
	}
}

func (s *serviceLineItemService) DeleteServiceLineItem(ctx context.Context, request *servicelineitempb.DeleteServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ServiceLineItemService.DeleteServiceLineItem")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Validate service line item ID
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	deleteSliCommand := command.NewDeleteServiceLineItemCommand(request.Id, request.Tenant, request.LoggedInUserId, request.AppSource)

	if err := s.serviceLineItemCommandHandlers.DeleteServiceLineItem.Handle(ctx, deleteSliCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(DeleteServiceLineItem.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the updated service line item
	return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: request.Id}, nil
}

func (s *serviceLineItemService) CloseServiceLineItem(ctx context.Context, request *servicelineitempb.CloseServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ServiceLineItemService.CloseServiceLineItem")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Validate service line item ID
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	closeSliCommand := command.NewCloseServiceLineItemCommand(request.Id, request.Tenant, request.LoggedInUserId, request.AppSource,
		utils.TimestampProtoToTimePtr(request.EndedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))

	if err := s.serviceLineItemCommandHandlers.CloseServiceLineItem.Handle(ctx, closeSliCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CloseServiceLineItem.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the updated service line item
	return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: request.Id}, nil
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
