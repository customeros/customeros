package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/command_handler"
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
