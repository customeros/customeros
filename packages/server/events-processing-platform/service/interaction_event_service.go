package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	interactioneventpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/model"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"strings"
)

type interactionEventService struct {
	interactioneventpb.UnimplementedInteractionEventGrpcServiceServer
	log                              logger.Logger
	interactionEventsCommandHandlers *cmdhnd.CommandHandlers
}

func NewInteractionEventService(log logger.Logger, commands *cmdhnd.CommandHandlers) *interactionEventService {
	return &interactionEventService{
		log:                              log,
		interactionEventsCommandHandlers: commands,
	}
}

func (s *interactionEventService) UpsertInteractionEvent(ctx context.Context, request *interactioneventpb.UpsertInteractionEventGrpcRequest) (*interactioneventpb.InteractionEventIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionEventService.UpsertInteractionEvent")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	interactionEventId := strings.TrimSpace(utils.NewUUIDIfEmpty(request.Id))

	dataFields := model.InteractionEventDataFields{
		Content:            request.Content,
		ContentType:        request.ContentType,
		Channel:            request.Channel,
		ChannelData:        request.ChannelData,
		EventType:          request.EventType,
		Identifier:         request.Identifier,
		BelongsToIssueId:   request.BelongsToIssueId,
		BelongsToSessionId: request.BelongsToSessionId,
		Hide:               request.Hide,
	}
	if request.Sender != nil {
		dataFields.Sender = model.Sender{
			Participant: commonmodel.Participant{
				ID:              request.Sender.Participant.Id,
				ParticipantType: GetParticipantTypeFromPB(request.Sender.Participant),
			},
			RelationType: request.Sender.RelationType,
		}
	}
	dataFields.Receivers = make([]model.Receiver, len(request.Receivers))
	for _, receiver := range request.Receivers {
		dataFields.Receivers = append(dataFields.Receivers, model.Receiver{
			Participant: commonmodel.Participant{
				ID:              receiver.Participant.Id,
				ParticipantType: GetParticipantTypeFromPB(receiver.Participant),
			},
			RelationType: receiver.RelationType,
		})
	}

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	cmd := command.NewUpsertInteractionEventCommand(interactionEventId, request.Tenant, request.LoggedInUserId, dataFields, source, externalSystem, utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err := s.interactionEventsCommandHandlers.UpsertInteractionEvent.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewUpsertInteractionEventCommand.Handle) tenant:{%v}, interactionEventId:{%v} , err: %v", request.Tenant, request.Id, err.Error())
		return nil, s.errResponse(err)
	}

	return &interactioneventpb.InteractionEventIdGrpcResponse{Id: interactionEventId}, nil
}

func (s *interactionEventService) RequestGenerateSummary(ctx context.Context, request *interactioneventpb.RequestGenerateSummaryGrpcRequest) (*interactioneventpb.InteractionEventIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionEventService.RequestGenerateSummary")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewRequestSummaryCommand(request.Tenant, request.InteractionEventId)
	if err := s.interactionEventsCommandHandlers.RequestSummary.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error handling RequestSummary command: %v", err.Error())
		return nil, s.errResponse(err)
	}

	return &interactioneventpb.InteractionEventIdGrpcResponse{Id: request.InteractionEventId}, nil
}

func (s *interactionEventService) RequestGenerateActionItems(ctx context.Context, request *interactioneventpb.RequestGenerateActionItemsGrpcRequest) (*interactioneventpb.InteractionEventIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionEventService.RequestGenerateActionItems")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewRequestActionItemsCommand(request.Tenant, request.InteractionEventId)
	if err := s.interactionEventsCommandHandlers.RequestActionItems.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error handling RequestActionItems command: %v", err.Error())
		return nil, s.errResponse(err)
	}

	return &interactioneventpb.InteractionEventIdGrpcResponse{Id: request.InteractionEventId}, nil
}

func (s *interactionEventService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}

func GetParticipantTypeFromPB(participant *interactioneventpb.Participant) commonmodel.ParticipantType {
	if participant == nil {
		return ""
	}
	switch participant.ParticipantType.(type) {
	case *interactioneventpb.Participant_User:
		return commonmodel.UserType
	case *interactioneventpb.Participant_Contact:
		return commonmodel.ContactType
	case *interactioneventpb.Participant_Organization:
		return commonmodel.OrganizationType
	default:
		return ""
	}
}
