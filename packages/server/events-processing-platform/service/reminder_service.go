package service

import (
	"context"

	"github.com/google/uuid"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/event_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
)

type reminderService struct {
	reminderpb.UnimplementedReminderGrpcServiceServer
	log           logger.Logger
	eventHandlers *event_handler.EventHandlers
}

func NewReminderService(log logger.Logger, commandHandlers *event_handler.EventHandlers) *reminderService {
	return &reminderService{
		log:           log,
		eventHandlers: commandHandlers,
	}
}

func (s *reminderService) CreateReminder(ctx context.Context, request *reminderpb.CreateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ReminderService.CreateReminder")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.UserId)
	tracing.LogObjectAsJson(span, "request", request)

	reminderId := uuid.New().String()

	baseRequest := eventstore.NewBaseRequest(reminderId, request.Tenant, request.UserId, commonmodel.SourceFromGrpc(request.SourceFields))

	if err := s.eventHandlers.CreateReminderHandler.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateReminder.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &reminderpb.ReminderGrpcResponse{Id: reminderId}, nil
}

func (s *reminderService) UpdateReminder(ctx context.Context, request *reminderpb.UpdateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ReminderService.UpdateReminder")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.ReminderId)
	tracing.LogObjectAsJson(span, "request", request)

	srcFields := commonmodel.Source{AppSource: request.AppSource}

	baseRequest := eventstore.NewBaseRequest(request.ReminderId, request.Tenant, request.ReminderId, srcFields)

	if err := s.eventHandlers.UpdateReminderHandler.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateReminder.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &reminderpb.ReminderGrpcResponse{Id: request.ReminderId}, nil
}
