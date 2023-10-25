package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	pb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/models"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type logEntryService struct {
	pb.UnimplementedLogEntryGrpcServiceServer
	log              logger.Logger
	logEntryCommands *cmdhnd.LogEntryCommandHandlers
}

func NewLogEntryService(log logger.Logger, logEntryCommands *cmdhnd.LogEntryCommandHandlers) *logEntryService {
	return &logEntryService{
		log:              log,
		logEntryCommands: logEntryCommands,
	}
}

func (s *logEntryService) UpsertLogEntry(ctx context.Context, request *pb.UpsertLogEntryGrpcRequest) (*pb.LogEntryIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "LogEntryService.UpsertLogEntry")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.UserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	logEntryId := request.Id
	if strings.TrimSpace(logEntryId) == "" {
		logEntryId = uuid.New().String()
	}

	dataFields := models.LogEntryDataFields{
		Content:              request.Content,
		ContentType:          request.ContentType,
		StartedAt:            utils.TimestampProtoToTime(request.StartedAt),
		AuthorUserId:         request.AuthorUserId,
		LoggedOrganizationId: request.LoggedOrganizationId,
	}
	source := cmnmod.Source{}
	source.FromGrpc(request.SourceFields)
	externalSystem := cmnmod.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	command := cmd.NewUpsertLogEntryCommand(logEntryId, request.Tenant, request.UserId, source, externalSystem, dataFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.logEntryCommands.UpsertLogEntry.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertLogEntryCommand.Handle) tenant:{%s}, logEntryId:{%s} , err: %s", request.Tenant, request.Id, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Upserted logEntry %s", logEntryId)

	return &pb.LogEntryIdGrpcResponse{Id: logEntryId}, nil
}

func (s *logEntryService) AddTag(ctx context.Context, request *pb.AddTagGrpcRequest) (*pb.LogEntryIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "LogEntryService.Addtag")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.UserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	command := cmd.NewAddTagCommand(request.Id, request.Tenant, request.UserId, request.TagId, utils.TimePtr(utils.Now()))
	if err := s.logEntryCommands.AddTag.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddTag.Handle) tenant:%s, logEntryId: %s, tagId , err: %s", request.Tenant, request.Id, request.TagId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.LogEntryIdGrpcResponse{Id: request.Id}, nil
}

func (s *logEntryService) RemoveTag(ctx context.Context, request *pb.RemoveTagGrpcRequest) (*pb.LogEntryIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "LogEntryService.RemoveTag")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.UserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	command := cmd.NewRemoveTagCommand(request.Id, request.Tenant, request.UserId, request.TagId)
	if err := s.logEntryCommands.RemoveTag.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemoveTag.Handle) tenant:%s, logEntryId: %s, tagId , err: %s", request.Tenant, request.Id, request.TagId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.LogEntryIdGrpcResponse{Id: request.Id}, nil
}

func (s *logEntryService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
