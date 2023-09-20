package service

import (
	"context"
	"github.com/google/uuid"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	pb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/models"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/utils"
	"strings"
)

type logEntryService struct {
	pb.UnimplementedLogEntryGrpcServiceServer
	log              logger.Logger
	repositories     *repository.Repositories
	logEntryCommands *cmdhnd.LogEntryCommands
}

func NewLogEntryService(log logger.Logger, repositories *repository.Repositories, logEntryCommands *cmdhnd.LogEntryCommands) *logEntryService {
	return &logEntryService{
		log:              log,
		repositories:     repositories,
		logEntryCommands: logEntryCommands,
	}
}

func (s *logEntryService) UpsertLogEntry(ctx context.Context, request *pb.UpsertLogEntryGrpcRequest) (*pb.LogEntryIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "LogEntryService.UpsertLogEntry")
	defer span.Finish()

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
		s.log.Errorf("(UpsertSyncLogEntry.Handle) tenant:%s, logEntryId: %s , err: %s", request.Tenant, request.Id, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Upserted logEntry %s", logEntryId)

	return &pb.LogEntryIdGrpcResponse{Id: logEntryId}, nil
}

func (s *logEntryService) AddTag(ctx context.Context, request *pb.AddTagGrpcRequest) (*pb.LogEntryIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "LogEntryService.Addtag")
	defer span.Finish()

	command := cmd.NewAddTagCommand(request.Id, request.Tenant, request.UserId, request.TagId, common_utils.TimePtr(common_utils.Now()))
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
