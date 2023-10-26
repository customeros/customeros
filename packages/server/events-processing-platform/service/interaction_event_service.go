package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	iepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/model"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type interactionEventService struct {
	iepb.UnimplementedInteractionEventGrpcServiceServer
	log                              logger.Logger
	interactionEventsCommandHandlers *cmdhnd.InteractionEventCommandHandlers
}

func NewInteractionEventService(log logger.Logger, commands *cmdhnd.InteractionEventCommandHandlers) *interactionEventService {
	return &interactionEventService{
		log:                              log,
		interactionEventsCommandHandlers: commands,
	}
}

func (s *interactionEventService) UpsertInteractionEvent(ctx context.Context, request *iepb.UpsertInteractionEventGrpcRequest) (*iepb.InteractionEventIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionEventService.UpsertInteractionEvent")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	interactionEventId := strings.TrimSpace(utils.NewUUIDIfEmpty(request.Id))

	dataFields := model.InteractionEventDataFields{
		Content:         request.Content,
		ContentType:     request.ContentType,
		Channel:         request.Channel,
		ChannelData:     request.ChannelData,
		EventType:       request.EventType,
		Identifier:      request.Identifier,
		PartOfIssueId:   request.PartOfIssueId,
		PartOfSessionId: request.PartOfSessionId,
		Hide:            request.Hide,
	}

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	cmd := command.NewUpsertInteractionEventCommand(interactionEventId, request.Tenant, request.LoggedInUserId, dataFields, source, externalSystem, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.interactionEventsCommandHandlers.UpsertInteractionEvent.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewUpsertInteractionEventCommand.Handle) tenant:{%v}, interactionEventId:{%v} , err: %v", request.Tenant, request.Id, err.Error())
		return nil, s.errResponse(err)
	}

	return &iepb.InteractionEventIdGrpcResponse{Id: interactionEventId}, nil
}

func (s *interactionEventService) RequestGenerateSummary(ctx context.Context, request *iepb.RequestGenerateSummaryGrpcRequest) (*iepb.InteractionEventIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionEventService.RequestGenerateSummary")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewRequestSummaryCommand(request.Tenant, request.InteractionEventId)
	if err := s.interactionEventsCommandHandlers.RequestSummary.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error handling RequestSummary command: %v", err.Error())
		return nil, s.errResponse(err)
	}

	return &iepb.InteractionEventIdGrpcResponse{Id: request.InteractionEventId}, nil
}

func (s *interactionEventService) RequestGenerateActionItems(ctx context.Context, request *iepb.RequestGenerateActionItemsGrpcRequest) (*iepb.InteractionEventIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionEventService.RequestGenerateActionItems")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewRequestActionItemsCommand(request.Tenant, request.InteractionEventId)
	if err := s.interactionEventsCommandHandlers.RequestActionItems.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error handling RequestActionItems command: %v", err.Error())
		return nil, s.errResponse(err)
	}

	return &iepb.InteractionEventIdGrpcResponse{Id: request.InteractionEventId}, nil
}

func (s *interactionEventService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
