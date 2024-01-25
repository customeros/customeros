package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
)

type tenantService struct {
	tenantpb.UnimplementedTenantGrpcServiceServer
	log                  logger.Logger
	tenantRequestHandler tenant.TenantRequestHandler
	aggregateStore       eventstore.AggregateStore
}

func NewTenantService(log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config) *tenantService {
	return &tenantService{
		log:                  log,
		tenantRequestHandler: tenant.NewTenantRequestHandler(log, aggregateStore, cfg.Utils),
		aggregateStore:       aggregateStore,
	}
}

func (s *tenantService) AddBillingProfile(ctx context.Context, request *tenantpb.AddBillingProfileRequest) (*commonpb.IdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "TenantService.AddBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	billingProfileId, err := s.tenantRequestHandler.HandleWithRetry(ctx, request.Tenant, false, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewTenant) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commonpb.IdResponse{Id: billingProfileId.(string)}, nil
}

func (s *tenantService) UpdateBillingProfile(ctx context.Context, request *tenantpb.UpdateBillingProfileRequest) (*commonpb.IdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "TenantService.UpdateBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.tenantRequestHandler.HandleWithRetry(ctx, request.Tenant, false, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewTenant) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commonpb.IdResponse{Id: request.Id}, nil
}
