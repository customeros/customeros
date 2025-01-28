package service

import (
	"context"

	organizationpb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/customeros/customeros/packages/server/events-processing-platform/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization"
	grpcerr "github.com/customeros/customeros/packages/server/events-processing-platform/grpc_errors"
	"github.com/customeros/customeros/packages/server/events-processing-platform/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform/tracing"
)

type organizationService struct {
	organizationpb.UnimplementedOrganizationGrpcServiceServer
	log                        logger.Logger
	organizationRequestHandler organization.OrganizationRequestHandler
	requestHandlerService      RequestHandler
}

func NewOrganizationService(log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config, req RequestHandler) *organizationService {
	return &organizationService{
		log:                        log,
		organizationRequestHandler: organization.NewOrganizationRequestHandler(log, aggregateStore, cfg.Utils),
		requestHandlerService:      req,
	}
}

func (s *organizationService) RefreshRenewalSummary(ctx context.Context, request *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RefreshRenewalSummary")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	_, err := s.organizationRequestHandler.HandleTempWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to refresh renewal summary for organization with id {%s} for tenant {%s}, err: %s", request.OrganizationId, request.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}

func (s *organizationService) CreateBillingProfile(ctx context.Context, request *organizationpb.CreateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.CreateBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	result, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	billingProfileId := ""
	if result != nil {
		billingProfileId = result.(string)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: billingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: billingProfileId}, nil
}

func (s *organizationService) UpdateBillingProfile(ctx context.Context, request *organizationpb.UpdateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) LinkEmailToBillingProfile(ctx context.Context, request *organizationpb.LinkEmailToBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkEmailToBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkEmailToBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) UnlinkEmailFromBillingProfile(ctx context.Context, request *organizationpb.UnlinkEmailFromBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UnlinkEmailFromBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UnlinkEmailFromBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) LinkLocationToBillingProfile(ctx context.Context, request *organizationpb.LinkLocationToBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkLocationToBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkLocationToBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) UnlinkLocationFromBillingProfile(ctx context.Context, request *organizationpb.UnlinkLocationFromBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UnlinkLocationFromBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UnlinkLocationFromBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}
