package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type contractService struct {
	contractpb.UnimplementedContractGrpcServiceServer
	services       *Services
	log            logger.Logger
	aggregateStore eventstore.AggregateStore
}

func NewContractService(log logger.Logger, aggregateStore eventstore.AggregateStore, services *Services) *contractService {
	return &contractService{
		services:       services,
		log:            log,
		aggregateStore: aggregateStore,
	}
}

func (s *contractService) RefreshContractStatus(ctx context.Context, request *contractpb.RefreshContractStatusGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContractService.RefreshContractStatus")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewContractTempAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{
		SkipLoadEvents: true,
	}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RefreshContractStatus.HandleTemp) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &contractpb.ContractIdGrpcResponse{Id: request.Id}, nil
}

func (s *contractService) RefreshContractLtv(ctx context.Context, request *contractpb.RefreshContractLtvGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContractService.RefreshContractLtv")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewContractTempAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{
		SkipLoadEvents: true,
	}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RefreshContractLtv.HandleGRPCRequest) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &contractpb.ContractIdGrpcResponse{Id: request.Id}, nil
}

func (s *contractService) RolloutRenewalOpportunityOnExpiration(ctx context.Context, request *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContractService.RefreshContractStatus")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewContractAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RolloutRenewalOpportunityOnExpiration.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &contractpb.ContractIdGrpcResponse{Id: request.Id}, nil
}

func (s *contractService) SoftDeleteContract(ctx context.Context, request *contractpb.SoftDeleteContractGrpcRequest) (*emptypb.Empty, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContractService.SoftDeleteContract")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// Validate contract ID
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewContractAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(DeleteContract.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &emptypb.Empty{}, nil
}
