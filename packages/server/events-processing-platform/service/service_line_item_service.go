package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	sliaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serviceLineItemService struct {
	servicelineitempb.UnimplementedServiceLineItemGrpcServiceServer
	log            logger.Logger
	aggregateStore eventstore.AggregateStore
	services       *Services
}

func NewServiceLineItemService(log logger.Logger, aggregateStore eventstore.AggregateStore, services *Services) *serviceLineItemService {
	return &serviceLineItemService{
		log:            log,
		aggregateStore: aggregateStore,
		services:       services,
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

	initAggregateFunc := func() eventstore.Aggregate {
		return sliaggregate.NewServiceLineItemAggregateWithTenantAndID(request.Tenant, serviceLineItemId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateServiceLineItem) tenant:{%v}, err: %v", request.Tenant, err.Error())
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

	initAggregateFunc := func() eventstore.Aggregate {
		return sliaggregate.NewServiceLineItemAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(DeleteServiceLineItem) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the updated service line item
	return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: request.Id}, nil
}

func (s *serviceLineItemService) DeleteServiceLineItem(ctx context.Context, request *servicelineitempb.DeleteServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ServiceLineItemService.DeleteServiceLineItem")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// Validate service line item ID
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return sliaggregate.NewServiceLineItemAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(DeleteServiceLineItem) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the deleted service line item
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

	extraParams := map[string]any{}
	extraParams[model.PARAM_CANCELLED] = true

	initAggregateFunc := func() eventstore.Aggregate {
		return sliaggregate.NewServiceLineItemAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request, extraParams); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CloseServiceLineItem) tenant:{%v}, err: %v", request.Tenant, err.Error())
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

func (s *serviceLineItemService) checkSLINotEnded(ctx context.Context, tenant, id string) (bool, error) {
	sliAggregate := sliaggregate.NewServiceLineItemAggregateWithTenantAndID(tenant, id)
	err := s.aggregateStore.Exists(ctx, sliAggregate.GetID())
	if err != nil {
		if errors.Is(err, eventstore.ErrAggregateNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	return sliAggregate.ServiceLineItem.EndedAt != nil, nil
}

func (s *serviceLineItemService) PauseServiceLineItem(ctx context.Context, request *servicelineitempb.PauseServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ServiceLineItemService.PauseServiceLineItem")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Validate service line item ID
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return sliaggregate.NewServiceLineItemAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PauseServiceLineItem) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the updated service line item
	return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: request.Id}, nil
}

func (s *serviceLineItemService) ResumeServiceLineItem(ctx context.Context, request *servicelineitempb.ResumeServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ServiceLineItemService.ResumeServiceLineItem")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Validate service line item ID
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return sliaggregate.NewServiceLineItemAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(ResumeServiceLineItem) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the updated service line item
	return &servicelineitempb.ServiceLineItemIdGrpcResponse{Id: request.Id}, nil
}
