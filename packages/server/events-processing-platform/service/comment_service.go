package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/command"
	commentcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/model"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
)

type commentService struct {
	commentpb.UnimplementedCommentGrpcServiceServer
	log                    logger.Logger
	commentCommandHandlers *commentcmdhandler.CommandHandlers
}

func NewCommentService(log logger.Logger, commentCommandHandlers *commentcmdhandler.CommandHandlers) *commentService {
	return &commentService{
		log:                    log,
		commentCommandHandlers: commentCommandHandlers,
	}
}

func (s *commentService) UpsertComment(ctx context.Context, request *commentpb.UpsertCommentGrpcRequest) (*commentpb.CommentIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "CommentService.UpsertComment")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.UserId)
	tracing.LogObjectAsJson(span, "request", request)

	commentId := utils.NewUUIDIfEmpty(request.Id)

	dataFields := model.CommentDataFields{
		Content:          request.Content,
		ContentType:      request.ContentType,
		AuthorUserId:     request.AuthorUserId,
		CommentedIssueId: request.CommentedIssueId,
	}
	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)
	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	cmd := command.NewUpsertCommentCommand(commentId, request.Tenant, request.UserId, source, externalSystem, dataFields, utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err := s.commentCommandHandlers.Upsert.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertCommentCommand.Handle) tenant:{%s}, commentId:{%s} , err: %s", request.Tenant, commentId, err.Error())
		return nil, s.errResponse(err)
	}

	return &commentpb.CommentIdGrpcResponse{Id: commentId}, nil
}

func (s *commentService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
