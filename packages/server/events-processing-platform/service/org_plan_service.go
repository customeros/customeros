package service

import (
	"github.com/google/uuid"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/event_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	orgplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	"golang.org/x/net/context"
)

type organizationPlanService struct {
	orgplanpb.UnimplementedOrganizationPlanGrpcServiceServer
	log            logger.Logger
	eventHandlers  *event_handler.EventHandlers
	aggregateStore eventstore.AggregateStore
}

func NewOrganizationPlanService(log logger.Logger, commandHandlers *event_handler.EventHandlers, aggregateStore eventstore.AggregateStore) *organizationPlanService {
	return &organizationPlanService{
		log:            log,
		eventHandlers:  commandHandlers,
		aggregateStore: aggregateStore,
	}
}

func (s *organizationPlanService) CreateOrganizationPlan(ctx context.Context, request *orgplanpb.CreateOrganizationPlanGrpcRequest) (*orgplanpb.OrganizationPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationPlanService.CreateOrganizationPlan")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	organizationPlanId := uuid.New().String()

	baseRequest := eventstore.NewBaseRequest(organizationPlanId, request.Tenant, request.LoggedInUserId, commonmodel.SourceFromGrpc(request.SourceFields))

	if err := s.eventHandlers.CreateOrganizationPlan.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateOrganizationPlan.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created master plan
	return &orgplanpb.OrganizationPlanIdGrpcResponse{Id: organizationPlanId}, nil
}

func (s *organizationPlanService) CreateOrganizationPlanMilestone(ctx context.Context, request *orgplanpb.CreateOrganizationPlanMilestoneGrpcRequest) (*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationPlanService.CreateOrganizationPlanMilestone")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	milestoneId := uuid.New().String()

	baseRequest := eventstore.NewBaseRequest(milestoneId, request.Tenant, request.LoggedInUserId, commonmodel.SourceFromGrpc(request.SourceFields))

	if err := s.eventHandlers.CreateOrganizationPlanMilestone.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateOrganizationPlanMilestone.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orgplanpb.OrganizationPlanMilestoneIdGrpcResponse{Id: milestoneId}, nil
}

func (s *organizationPlanService) UpdateOrganizationPlan(ctx context.Context, request *orgplanpb.UpdateOrganizationPlanGrpcRequest) (*orgplanpb.OrganizationPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationPlanService.UpdateOrganizationPlan")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationPlanId)

	if request.OrganizationPlanId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("organizationPlanId"))
	}

	srcFields := commonmodel.Source{AppSource: request.AppSource}

	baseRequest := eventstore.NewBaseRequest(request.OrganizationPlanId, request.Tenant, request.LoggedInUserId, srcFields)

	if err := s.eventHandlers.UpdateOrganizationPlan.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOrganizationPlan.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orgplanpb.OrganizationPlanIdGrpcResponse{Id: request.OrganizationPlanId}, nil
}

func (s *organizationPlanService) UpdateOrganizationPlanMilestone(ctx context.Context, request *orgplanpb.UpdateOrganizationPlanMilestoneGrpcRequest) (*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationPlanService.UpdateOrganizationPlanMilestone")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationPlanId)

	if request.OrganizationPlanId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("organizationPlanId"))
	}
	if request.OrganizationPlanMilestoneId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("organizationPlanMilestoneId"))
	}

	srcFields := commonmodel.Source{AppSource: request.AppSource}

	baseRequest := eventstore.NewBaseRequest(request.OrganizationPlanId, request.Tenant, request.LoggedInUserId, srcFields)

	if err := s.eventHandlers.UpdateOrganizationPlanMilestone.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOrganizationPlanMilestone.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created master plan
	return &orgplanpb.OrganizationPlanMilestoneIdGrpcResponse{Id: request.OrganizationPlanMilestoneId}, nil
}

func (s *organizationPlanService) ReorderOrganizationPlanMilestones(ctx context.Context, request *orgplanpb.ReorderOrganizationPlanMilestonesGrpcRequest) (*orgplanpb.OrganizationPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationPlanService.ReorderOrganizationPlanMilestones")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.OrganizationPlanId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("organizationPlanId"))
	}
	if len(request.OrganizationPlanMilestoneIds) == 0 {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("organizationPlanMilestoneIds"))
	}

	srcFields := commonmodel.Source{AppSource: request.AppSource}

	baseRequest := eventstore.NewBaseRequest(request.OrganizationPlanId, request.Tenant, request.LoggedInUserId, srcFields)

	if err := s.eventHandlers.ReorderOrganizationPlanMilestones.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOrganizationPlanMilestone.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orgplanpb.OrganizationPlanIdGrpcResponse{Id: request.OrganizationPlanId}, nil
}
