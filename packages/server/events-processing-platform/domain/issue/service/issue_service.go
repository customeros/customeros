package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	pb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type issueService struct {
	pb.UnimplementedIssueGrpcServiceServer
	log                  logger.Logger
	issueCommandHandlers *cmdhnd.IssueCommandHandlers
}

func NewIssueService(log logger.Logger, issueCommandHandlers *cmdhnd.IssueCommandHandlers) *issueService {
	return &issueService{
		log:                  log,
		issueCommandHandlers: issueCommandHandlers,
	}
}

func (s *issueService) UpsertIssue(ctx context.Context, request *pb.UpsertIssueGrpcRequest) (*pb.IssueIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "IssueService.UpsertIssue")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	issueId := strings.TrimSpace(utils.NewUUIDIfEmpty(request.Id))

	dataFields := model.IssueDataFields{
		Subject:                   request.Subject,
		Description:               request.Description,
		Status:                    request.Status,
		Priority:                  request.Priority,
		ReportedByOrganizationId:  request.ReportedByOrganizationId,
		SubmittedByOrganizationId: request.SubmittedByOrganizationId,
		SubmittedByUserId:         request.SubmittedByUserId,
	}

	source := cmnmod.Source{}
	source.FromGrpc(request.SourceFields)

	externalSystem := cmnmod.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	cmd := command.NewUpsertIssueCommand(issueId, request.Tenant, request.LoggedInUserId, dataFields, source, externalSystem, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.issueCommandHandlers.UpsertIssue.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertIssueCommand.Handle) tenant:{%v}, issueId:{%v} , err: %v", request.Tenant, request.Id, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.IssueIdGrpcResponse{Id: issueId}, nil
}

func (s *issueService) AddUserAssignee(ctx context.Context, request *pb.AddUserAssigneeToIssueGrpcRequest) (*pb.IssueIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "IssueService.AddUserAssignee")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewAddUserAssigneeCommand(request.IssueId, request.Tenant, request.LoggedInUserId, request.UserId, request.AppSource, nil)
	if err := s.issueCommandHandlers.AddUserAssignee.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddUserAssigneeCommand.Handle) tenant:{%v}, issueId:{%v} , err: %v", request.Tenant, request.IssueId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.IssueIdGrpcResponse{Id: request.IssueId}, nil
}

func (s *issueService) RemoveUserAssignee(ctx context.Context, request *pb.RemoveUserAssigneeFromIssueGrpcRequest) (*pb.IssueIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "IssueService.RemoveUserAssignee")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewRemoveUserAssigneeCommand(request.IssueId, request.Tenant, request.LoggedInUserId, request.UserId, request.AppSource, nil)
	if err := s.issueCommandHandlers.RemoveUserAssignee.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemoveUserAssigneeCommand.Handle) tenant:{%v}, issueId:{%v} , err: %v", request.Tenant, request.IssueId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.IssueIdGrpcResponse{Id: request.IssueId}, nil
}

func (s *issueService) AddUserFollower(ctx context.Context, request *pb.AddUserFollowerToIssueGrpcRequest) (*pb.IssueIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "IssueService.AddUserFollower")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewAddUserFollowerCommand(request.IssueId, request.Tenant, request.LoggedInUserId, request.UserId, request.AppSource, nil)
	if err := s.issueCommandHandlers.AddUserFollower.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddUserFollowerCommand.Handle) tenant:{%v}, issueId:{%v} , err: %v", request.Tenant, request.IssueId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.IssueIdGrpcResponse{Id: request.IssueId}, nil
}

func (s *issueService) RemoveUserFollower(ctx context.Context, request *pb.RemoveUserFollowerFromIssueGrpcRequest) (*pb.IssueIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "IssueService.RemoveUserFollower")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewRemoveUserFollowerCommand(request.IssueId, request.Tenant, request.LoggedInUserId, request.UserId, request.AppSource, nil)
	if err := s.issueCommandHandlers.RemoveUserFollower.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemoveUserFollowerCommand.Handle) tenant:{%v}, issueId:{%v} , err: %v", request.Tenant, request.IssueId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.IssueIdGrpcResponse{Id: request.IssueId}, nil
}

func (s *issueService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
