package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/model"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	interactionsessionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_session"
	"strings"
)

type interactionSessionService struct {
	interactionsessionpb.UnimplementedInteractionSessionGrpcServiceServer
	log                                logger.Logger
	interactionSessionsCommandHandlers *cmdhnd.CommandHandlers
}

func NewInteractionSessionService(log logger.Logger, commands *cmdhnd.CommandHandlers) *interactionSessionService {
	return &interactionSessionService{
		log:                                log,
		interactionSessionsCommandHandlers: commands,
	}
}

func (s *interactionSessionService) UpsertInteractionSession(ctx context.Context, request *interactionsessionpb.UpsertInteractionSessionGrpcRequest) (*interactionsessionpb.InteractionSessionIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InteractionSessionService.UpsertInteractionSession")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	interactionSessionId := strings.TrimSpace(utils.NewUUIDIfEmpty(request.Id))

	dataFields := model.InteractionSessionDataFields{
		Channel:     request.Channel,
		ChannelData: request.ChannelData,
		Identifier:  request.Identifier,
		Type:        request.Type,
		Name:        request.Name,
		Status:      request.Status,
	}

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	cmd := command.NewUpsertInteractionSessionCommand(interactionSessionId, request.Tenant, request.LoggedInUserId, dataFields, source, externalSystem, utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err := s.interactionSessionsCommandHandlers.UpsertInteractionSession.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewUpsertInteractionSessionCommand.Handle) tenant:{%v}, interactionSessionId:{%v} , err: %v", request.Tenant, request.Id, err.Error())
		return nil, s.errResponse(err)
	}

	return &interactionsessionpb.InteractionSessionIdGrpcResponse{Id: interactionSessionId}, nil
}

func (s *interactionSessionService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
