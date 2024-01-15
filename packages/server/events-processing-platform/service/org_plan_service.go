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

type organizationPlanService struct {
	orgplanpb.UnimplementedOrganizationPlanGrpcServiceServer
	log                             logger.Logger
	organizationPlanCommandHandlers *command_handler.CommandHandlers
	aggregateStore                  eventstore.AggregateStore
}

func NewOrganizationPlanService(log logger.Logger, commandHandlers *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore) *organizationPlanService {
	return &organizationPlanService{
		log:                             log,
		organizationPlanCommandHandlers: commandHandlers,
		aggregateStore:                  aggregateStore,
	}
}

func (s *organizationPlanService) CreateOrganizationPlan(ctx context.Context, request *orgplanpb.CreateOrganizationPlanGrpcRequest) (*orgplanpb.OrganizationPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationPlanService.CreateOrganizationPlan")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	organizationPlanId := uuid.New().String()

	createdAt := utils.TimestampProtoToTimePtr(request.CreatedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createOrganizationPlanCommand := command.NewCreateOrganizationPlanCommand(
		organizationPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.Name,
		sourceFields,
		createdAt,
	)

	if err := s.organizationPlanCommandHandlers.CreateOrganizationPlan.Handle(ctx, createOrganizationPlanCommand); err != nil {
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

	createdAt := utils.TimestampProtoToTimePtr(request.CreatedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createOrganizationPlanMilestoneCommand := command.NewCreateOrganizationPlanMilestoneCommand(
		request.OrganizationPlanId,
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

	if err := s.organizationPlanCommandHandlers.CreateOrganizationPlanMilestone.Handle(ctx, createOrganizationPlanMilestoneCommand); err != nil {
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

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	cmd := command.NewUpdateOrganizationPlanCommand(
		request.OrganizationPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.AppSource,
		request.Name,
		request.Retired,
		updatedAt,
		extractOrganizationPlanFieldsMask(request.FieldsMask),
	)

	if err := s.organizationPlanCommandHandlers.UpdateOrganizationPlan.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOrganizationPlan.Handle) tenant:{%v}, organizationPlanId:{%v}, err: %v", request.Tenant, request.OrganizationPlanId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orgplanpb.OrganizationPlanIdGrpcResponse{Id: request.OrganizationPlanId}, nil
}

func (s *organizationPlanService) UpdateOrganizationPlanMilestone(ctx context.Context, request *orgplanpb.UpdateOrganizationPlanMilestoneGrpcRequest) (*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationPlanService.UpdateOrganizationPlanMilestone")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.OrganizationPlanId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("organizationPlanId"))
	}
	if request.OrganizationPlanMilestoneId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("organizationPlanMilestoneId"))
	}

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	updateOrganizationPlanMilestoneCommand := command.NewUpdateOrganizationPlanMilestoneCommand(
		request.OrganizationPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.OrganizationPlanMilestoneId,
		request.Name,
		request.AppSource,
		request.Order,
		request.DurationHours,
		request.Items, // FIXME(@max-openline)
		request.Optional,
		request.Retired,
		updatedAt,
		extractOrganizationPlanMilestoneFieldsMask(request.FieldsMask),
	)

	if err := s.organizationPlanCommandHandlers.UpdateOrganizationPlanMilestone.Handle(ctx, updateOrganizationPlanMilestoneCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOrganizationPlanMilestone.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created master plan
	return &orgplanpb.OrganizationPlanMilestoneIdGrpcResponse{Id: request.OrganizationPlanMilestoneId}, nil
}

func extractOrganizationPlanFieldsMask(fields []orgplanpb.OrganizationPlanFieldMask) []string {
	fieldsMask := make([]string, 0)
	if fields == nil || len(fields) == 0 {
		return fieldsMask
	}
	if containsOrganizationPlanMaskFieldAll(fields) {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_RETIRED:
			fieldsMask = append(fieldsMask, event.FieldMaskRetired)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}

func containsOrganizationPlanMaskFieldAll(fields []orgplanpb.OrganizationPlanFieldMask) bool {
	for _, field := range fields {
		if field == orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_ALL {
			return true
		}
	}
	return false
}

func extractOrganizationPlanMilestoneFieldsMask(fields []orgplanpb.OrganizationPlanMilestoneFieldMask) []string {
	fieldsMask := make([]string, 0)
	if fields == nil || len(fields) == 0 {
		return fieldsMask
	}
	if containsOrganizationPlanMilestoneMaskFieldAll(fields) {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_RETIRED:
			fieldsMask = append(fieldsMask, event.FieldMaskRetired)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ORDER:
			fieldsMask = append(fieldsMask, event.FieldMaskOrder)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_OPTIONAL:
			fieldsMask = append(fieldsMask, event.FieldMaskOptional)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_DURATION_HOURS:
			fieldsMask = append(fieldsMask, event.FieldMaskDurationHours)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ITEMS:
			fieldsMask = append(fieldsMask, event.FieldMaskItems)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}

func containsOrganizationPlanMilestoneMaskFieldAll(fields []orgplanpb.OrganizationPlanMilestoneFieldMask) bool {
	for _, field := range fields {
		if field == orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ALL {
			return true
		}
	}
	return false
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

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	reorderOrganizationPlanMilestonesCommand := command.NewReorderOrganizationPlanMilestonesCommand(
		request.OrganizationPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.AppSource,
		request.OrganizationPlanMilestoneIds,
		updatedAt,
	)

	if err := s.organizationPlanCommandHandlers.ReorderOrganizationPlanMilestones.Handle(ctx, reorderOrganizationPlanMilestonesCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(ReorderOrganizationPlanMilestones.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orgplanpb.OrganizationPlanIdGrpcResponse{Id: request.OrganizationPlanId}, nil
}
