package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"google.golang.org/protobuf/types/known/emptypb"
)

type tenantService struct {
	tenantpb.UnimplementedTenantGrpcServiceServer
	services       *Services
	log            logger.Logger
	aggregateStore eventstore.AggregateStore
}

func NewTenantService(services *Services, log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config) *tenantService {
	return &tenantService{
		services:       services,
		log:            log,
		aggregateStore: aggregateStore,
	}
}

func (s *tenantService) AddBillingProfile(ctx context.Context, request *tenantpb.AddBillingProfileRequest) (*commonpb.IdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "TenantService.AddBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return tenant.NewTenantAggregate(request.Tenant)
	}
	billingProfileId, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptions(), request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddBillingProfile) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commonpb.IdResponse{Id: billingProfileId.(string)}, nil
}

func (s *tenantService) UpdateBillingProfile(ctx context.Context, request *tenantpb.UpdateBillingProfileRequest) (*commonpb.IdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "TenantService.UpdateBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return tenant.NewTenantAggregate(request.Tenant)
	}
	_, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptions(), request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateBillingProfile) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commonpb.IdResponse{Id: request.Id}, nil
}

func (s *tenantService) UpdateTenantSettings(ctx context.Context, request *tenantpb.UpdateTenantSettingsRequest) (*emptypb.Empty, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "TenantService.UpdateTenantSettings")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return tenant.NewTenantAggregate(request.Tenant)
	}
	_, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptions(), request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateTenantSettings) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *tenantService) AddBankAccount(ctx context.Context, request *tenantpb.AddBankAccountGrpcRequest) (*commonpb.IdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "TenantService.AddBankAccount")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return tenant.NewTenantAggregate(request.Tenant)
	}
	bankAccountId, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptions(), request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddBankAccount) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commonpb.IdResponse{Id: bankAccountId.(string)}, nil
}

func (s *tenantService) UpdateBankAccount(ctx context.Context, request *tenantpb.UpdateBankAccountGrpcRequest) (*commonpb.IdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "TenantService.UpdateBankAccount")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return tenant.NewTenantAggregate(request.Tenant)
	}
	_, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptions(), request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateBankAccount) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commonpb.IdResponse{Id: request.Id}, nil
}

func (s *tenantService) DeleteBankAccount(ctx context.Context, request *tenantpb.DeleteBankAccountGrpcRequest) (*emptypb.Empty, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "TenantService.DeleteBankAccount")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return tenant.NewTenantAggregate(request.Tenant)
	}
	_, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptions(), request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(DeleteBankAccount) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &emptypb.Empty{}, nil
}
