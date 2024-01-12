package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	orgplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	"golang.org/x/net/context"
)

type orgPlanService struct {
	orgplanpb.UnimplementedOrgPlanGrpcServiceServer
	log                    logger.Logger
	orgPlanCommandHandlers *command_handler.CommandHandlers
	aggregateStore         eventstore.AggregateStore
}

func NewOrgPlanService(log logger.Logger, commandHandlers *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore) *orgPlanService {
	return &orgPlanService{
		log:                    log,
		orgPlanCommandHandlers: commandHandlers,
		aggregateStore:         aggregateStore,
	}
}

func (s *orgPlanService) CreateOrgPlan(ctx context.Context, request *orgplanpb.CreateOrgPlanGrpcRequest) (*orgplanpb.OrgPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrgPlanService.CreateOrgPlan")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	orgPlanId := uuid.New().String()

	createdAt := utils.TimestampProtoToTimePtr(request.CreatedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createOrgPlanCommand := command.NewCreateOrgPlanCommand(
		orgPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.Name,
		sourceFields,
		createdAt,
	)

	if err := s.orgPlanCommandHandlers.CreateOrgPlan.Handle(ctx, createOrgPlanCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateOrgPlan.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created master plan
	return &orgplanpb.OrgPlanIdGrpcResponse{Id: orgPlanId}, nil
}

func (s *orgPlanService) CreateOrgPlanMilestone(ctx context.Context, request *orgplanpb.CreateOrgPlanMilestoneGrpcRequest) (*orgplanpb.OrgPlanMilestoneIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrgPlanService.CreateOrgPlanMilestone")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	milestoneId := uuid.New().String()

	createdAt := utils.TimestampProtoToTimePtr(request.CreatedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createOrgPlanMilestoneCommand := command.NewCreateOrgPlanMilestoneCommand(
		request.OrgPlanId,
		request.Tenant,
		request.LoggedInUserId,
		milestoneId,
		request.Name,
		request.Order,
		request.DurationHours,
		request.Items,
		request.Optional,
		sourceFields,
		createdAt,
	)

	if err := s.orgPlanCommandHandlers.CreateOrgPlanMilestone.Handle(ctx, createOrgPlanMilestoneCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateOrgPlanMilestone.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orgplanpb.OrgPlanMilestoneIdGrpcResponse{Id: milestoneId}, nil
}

func (s *orgPlanService) UpdateOrgPlan(ctx context.Context, request *orgplanpb.UpdateOrgPlanGrpcRequest) (*orgplanpb.OrgPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrgPlanService.UpdateOrgPlan")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrgPlanId)

	if request.OrgPlanId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("orgPlanId"))
	}

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	cmd := command.NewUpdateOrgPlanCommand(
		request.OrgPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.AppSource,
		request.Name,
		request.Retired,
		updatedAt,
		extractOrgPlanFieldsMask(request.FieldsMask),
	)

	if err := s.orgPlanCommandHandlers.UpdateOrgPlan.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOrgPlan.Handle) tenant:{%v}, orgPlanId:{%v}, err: %v", request.Tenant, request.OrgPlanId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orgplanpb.OrgPlanIdGrpcResponse{Id: request.OrgPlanId}, nil
}

func (s *orgPlanService) UpdateOrgPlanMilestone(ctx context.Context, request *orgplanpb.UpdateOrgPlanMilestoneGrpcRequest) (*orgplanpb.OrgPlanMilestoneIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrgPlanService.UpdateOrgPlanMilestone")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.OrgPlanId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("orgPlanId"))
	}
	if request.OrgPlanMilestoneId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("orgPlanMilestoneId"))
	}

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	updateOrgPlanMilestoneCommand := command.NewUpdateOrgPlanMilestoneCommand(
		request.OrgPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.OrgPlanMilestoneId,
		request.Name,
		request.AppSource,
		request.Order,
		request.DurationHours,
		request.Items, // FIXME(@max-openline)
		request.Optional,
		request.Retired,
		updatedAt,
		extractOrgPlanMilestoneFieldsMask(request.FieldsMask),
	)

	if err := s.orgPlanCommandHandlers.UpdateOrgPlanMilestone.Handle(ctx, updateOrgPlanMilestoneCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOrgPlanMilestone.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created master plan
	return &orgplanpb.OrgPlanMilestoneIdGrpcResponse{Id: request.OrgPlanMilestoneId}, nil
}

func extractOrgPlanFieldsMask(fields []orgplanpb.OrgPlanFieldMask) []string {
	fieldsMask := make([]string, 0)
	if fields == nil || len(fields) == 0 {
		return fieldsMask
	}
	if containsOrgPlanMaskFieldAll(fields) {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case orgplanpb.OrgPlanFieldMask_ORG_PLAN_PROPERTY_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case orgplanpb.OrgPlanFieldMask_ORG_PLAN_PROPERTY_RETIRED:
			fieldsMask = append(fieldsMask, event.FieldMaskRetired)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}

func containsOrgPlanMaskFieldAll(fields []orgplanpb.OrgPlanFieldMask) bool {
	for _, field := range fields {
		if field == orgplanpb.OrgPlanFieldMask_ORG_PLAN_PROPERTY_ALL {
			return true
		}
	}
	return false
}

func extractOrgPlanMilestoneFieldsMask(fields []orgplanpb.OrgPlanMilestoneFieldMask) []string {
	fieldsMask := make([]string, 0)
	if fields == nil || len(fields) == 0 {
		return fieldsMask
	}
	if containsOrgPlanMilestoneMaskFieldAll(fields) {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case orgplanpb.OrgPlanMilestoneFieldMask_ORG_PLAN_MILESTONE_PROPERTY_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case orgplanpb.OrgPlanMilestoneFieldMask_ORG_PLAN_MILESTONE_PROPERTY_RETIRED:
			fieldsMask = append(fieldsMask, event.FieldMaskRetired)
		case orgplanpb.OrgPlanMilestoneFieldMask_ORG_PLAN_MILESTONE_PROPERTY_ORDER:
			fieldsMask = append(fieldsMask, event.FieldMaskOrder)
		case orgplanpb.OrgPlanMilestoneFieldMask_ORG_PLAN_MILESTONE_PROPERTY_OPTIONAL:
			fieldsMask = append(fieldsMask, event.FieldMaskOptional)
		case orgplanpb.OrgPlanMilestoneFieldMask_ORG_PLAN_MILESTONE_PROPERTY_DURATION_HOURS:
			fieldsMask = append(fieldsMask, event.FieldMaskDurationHours)
		case orgplanpb.OrgPlanMilestoneFieldMask_ORG_PLAN_MILESTONE_PROPERTY_ITEMS:
			fieldsMask = append(fieldsMask, event.FieldMaskItems)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}

func containsOrgPlanMilestoneMaskFieldAll(fields []orgplanpb.OrgPlanMilestoneFieldMask) bool {
	for _, field := range fields {
		if field == orgplanpb.OrgPlanMilestoneFieldMask_ORG_PLAN_MILESTONE_PROPERTY_ALL {
			return true
		}
	}
	return false
}

func (s *orgPlanService) ReorderOrgPlanMilestones(ctx context.Context, request *orgplanpb.ReorderOrgPlanMilestonesGrpcRequest) (*orgplanpb.OrgPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrgPlanService.ReorderOrgPlanMilestones")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.OrgPlanId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("orgPlanId"))
	}
	if len(request.OrgPlanMilestoneIds) == 0 {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("orgPlanMilestoneIds"))
	}

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	reorderOrgPlanMilestonesCommand := command.NewReorderOrgPlanMilestonesCommand(
		request.OrgPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.AppSource,
		request.OrgPlanMilestoneIds,
		updatedAt,
	)

	if err := s.orgPlanCommandHandlers.ReorderOrgPlanMilestones.Handle(ctx, reorderOrgPlanMilestonesCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(ReorderOrgPlanMilestones.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orgplanpb.OrgPlanIdGrpcResponse{Id: request.OrgPlanId}, nil
}
