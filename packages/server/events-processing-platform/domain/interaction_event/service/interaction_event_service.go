package service

import (
	"context"
	"fmt"
	interaction_event_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/commands"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
)

type interactionEventService struct {
	interaction_event_grpc_service.UnimplementedInteractionEventGrpcServiceServer
	log                       logger.Logger
	repositories              *repository.Repositories
	interactionEventsCommands *commands.InteractionEventCommands
}

func NewInteractionEventService(log logger.Logger, repositories *repository.Repositories, commands *commands.InteractionEventCommands) *interactionEventService {
	return &interactionEventService{
		log:                       log,
		repositories:              repositories,
		interactionEventsCommands: commands,
	}
}

func (s *interactionEventService) RequestGenerateSummary(ctx context.Context, request *interaction_event_grpc_service.RequestGenerateSummaryGrpcRequest) (*interaction_event_grpc_service.InteractionEventIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionEventService.RequestGenerateSummary")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, "") // TODO enhance request with LoggedInUserId
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	command := commands.NewRequestSummaryCommand(request.Tenant, request.InteractionEventId)
	if err := s.interactionEventsCommands.RequestSummary.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error handling RequestSummary command: %v", err.Error())
		return nil, s.errResponse(err)
	}

	return &interaction_event_grpc_service.InteractionEventIdGrpcResponse{Id: request.InteractionEventId}, nil
}

func (s *interactionEventService) RequestGenerateActionItems(ctx context.Context, request *interaction_event_grpc_service.RequestGenerateActionItensGrpcRequest) (*interaction_event_grpc_service.InteractionEventIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionEventService.RequestGenerateActionItems")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, "") // TODO enhance request with LoggedInUserId
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	command := commands.NewRequestActionItemsCommand(request.Tenant, request.InteractionEventId)
	if err := s.interactionEventsCommands.RequestActionItems.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error handling RequestActionItems command: %v", err.Error())
		return nil, s.errResponse(err)
	}

	return &interaction_event_grpc_service.InteractionEventIdGrpcResponse{Id: request.InteractionEventId}, nil
}

func (s *interactionEventService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
