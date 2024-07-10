package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type commentService struct {
	commentpb.UnimplementedCommentGrpcServiceServer
	services       *Services
	log            logger.Logger
	aggregateStore eventstore.AggregateStore
}

func NewCommentService(services *Services, log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config) *commentService {
	return &commentService{
		services:       services,
		log:            log,
		aggregateStore: aggregateStore,
	}
}

func (s *commentService) UpsertComment(ctx context.Context, request *commentpb.UpsertCommentGrpcRequest) (*commentpb.CommentIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "CommentService.UpsertComment")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.UserId)
	tracing.LogObjectAsJson(span, "request", request)

	commentId := utils.NewUUIDIfEmpty(request.Id)

	initAggregateFunc := func() eventstore.Aggregate {
		return comment.NewCommentAggregateWithTenantAndID(request.Tenant, commentId)
	}
	_, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptions(), request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertComment) tenant:{%v}, commentId:{%s} ,err: %v", request.Tenant, commentId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commentpb.CommentIdGrpcResponse{Id: commentId}, nil
}

func (s *commentService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
