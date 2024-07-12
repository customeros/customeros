package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type reminderService struct {
	reminderpb.UnimplementedReminderGrpcServiceServer
	log            logger.Logger
	requestHandler reminder.ReminderRequestHandler
}

func NewReminderService(log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config, ebs *eventbuffer.EventBufferStoreService) *reminderService {
	return &reminderService{
		log:            log,
		requestHandler: reminder.NewReminderRequestHandler(log, aggregateStore, ebs, cfg.Utils),
	}
}

func (s *reminderService) CreateReminder(ctx context.Context, request *reminderpb.CreateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ReminderService.CreateReminder")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	reminderId := uuid.New().String()

	_, err := s.requestHandler.Handle(ctx, request.Tenant, reminderId, request)
	if err != nil {
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

	_, err := s.requestHandler.HandleWithRetry(ctx, request.Tenant, request.ReminderId, true, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateReminder.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &reminderpb.ReminderGrpcResponse{Id: request.ReminderId}, nil
}
