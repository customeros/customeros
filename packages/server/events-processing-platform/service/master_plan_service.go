package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	"golang.org/x/net/context"
)

type masterPlanService struct {
	masterplanpb.UnimplementedMasterPlanGrpcServiceServer
	log                       logger.Logger
	masterPlanCommandHandlers *command_handler.CommandHandlers
	aggregateStore            eventstore.AggregateStore
}

func NewMasterPlanService(log logger.Logger, commandHandlers *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore) *masterPlanService {
	return &masterPlanService{
		log:                       log,
		masterPlanCommandHandlers: commandHandlers,
		aggregateStore:            aggregateStore,
	}
}

func (s *masterPlanService) CreateMasterPlan(ctx context.Context, request *masterplanpb.CreateMasterPlanGrpcRequest) (*masterplanpb.MasterPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "MasterPlanService.CreateMasterPlan")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	masterPlanId := uuid.New().String()

	createdAt := utils.TimestampProtoToTimePtr(request.CreatedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createMasterPlanCommand := command.NewCreateMasterPlanCommand(
		masterPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.Name,
		sourceFields,
		createdAt,
	)

	if err := s.masterPlanCommandHandlers.CreateMasterPlan.Handle(ctx, createMasterPlanCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateMasterPlan.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created master plan
	return &masterplanpb.MasterPlanIdGrpcResponse{Id: masterPlanId}, nil
}

func (s *masterPlanService) CreateMasterPlanMilestone(ctx context.Context, request *masterplanpb.CreateMasterPlanMilestoneGrpcRequest) (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "MasterPlanService.CreateMasterPlanMilestone")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	milestoneId := uuid.New().String()

	createdAt := utils.TimestampProtoToTimePtr(request.CreatedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createMasterPlanMilestoneCommand := command.NewCreateMasterPlanMilestoneCommand(
		request.MasterPlanId,
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

	if err := s.masterPlanCommandHandlers.CreateMasterPlanMilestone.Handle(ctx, createMasterPlanMilestoneCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateMasterPlanMilestone.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created master plan
	return &masterplanpb.MasterPlanMilestoneIdGrpcResponse{Id: milestoneId}, nil
}

func (s *masterPlanService) UpdateMasterPlan(ctx context.Context, request *masterplanpb.UpdateMasterPlanGrpcRequest) (*masterplanpb.MasterPlanIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "MasterPlanService.UpdateMasterPlan")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.MasterPlanId)

	if request.MasterPlanId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("masterPlanId"))
	}

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	cmd := command.NewUpdateMasterPlanCommand(
		request.MasterPlanId,
		request.Tenant,
		request.LoggedInUserId,
		request.AppSource,
		request.Name,
		request.Retired,
		updatedAt,
		extractMasterPlanFieldsMask(request.FieldsMask),
	)

	if err := s.masterPlanCommandHandlers.UpdateMasterPlan.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateMasterPlan.Handle) tenant:{%v}, masterPlanId:{%v}, err: %v", request.Tenant, request.MasterPlanId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &masterplanpb.MasterPlanIdGrpcResponse{Id: request.MasterPlanId}, nil
}

func extractMasterPlanFieldsMask(fields []masterplanpb.MasterPlanMaskField) []string {
	fieldsMask := make([]string, 0)
	if fields == nil || len(fields) == 0 {
		return fieldsMask
	}
	if containsMasterPlanMaskFieldAll(fields) {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case masterplanpb.MasterPlanMaskField_MASTER_PLAN_PROPERTY_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case masterplanpb.MasterPlanMaskField_MASTER_PLAN_PROPERTY_RETIRED:
			fieldsMask = append(fieldsMask, event.FieldMaskRetired)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}

func containsMasterPlanMaskFieldAll(fields []masterplanpb.MasterPlanMaskField) bool {
	for _, field := range fields {
		if field == masterplanpb.MasterPlanMaskField_MASTER_PLAN_PROPERTY_ALL {
			return true
		}
	}
	return false
}
