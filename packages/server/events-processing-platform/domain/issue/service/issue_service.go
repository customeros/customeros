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
		Subject:                  request.Subject,
		Description:              request.Description,
		Status:                   request.Status,
		Priority:                 request.Priority,
		ReportedByOrganizationId: request.ReportedByOrganizationId,
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

func (s *issueService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
