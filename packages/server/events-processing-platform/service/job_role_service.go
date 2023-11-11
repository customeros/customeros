package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
)

type jobRoleService struct {
	jobrolepb.UnimplementedJobRoleGrpcServiceServer
	log             logger.Logger
	repositories    *repository.Repositories
	jobRoleCommands *commands.CommandHandlers
}

func NewJobRoleService(log logger.Logger, repositories *repository.Repositories, jobRoleCommands *commands.CommandHandlers) *jobRoleService {
	return &jobRoleService{
		log:             log,
		repositories:    repositories,
		jobRoleCommands: jobRoleCommands,
	}
}

func (jobRoleService *jobRoleService) CreateJobRole(ctx context.Context, request *jobrolepb.CreateJobRoleGrpcRequest) (*jobrolepb.JobRoleIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "JobRoleService.CreateJobRole")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, "") // TODO enhance request with LoggedInUserId
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	newObjectId, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new object ID: %w", err)
	}
	objectID := newObjectId.String()

	primary := false
	if request.Primary != nil {
		primary = *request.Primary
	}
	command := model.NewCreateJobRoleCommand(objectID, request.Tenant, request.JobTitle, request.Description, primary, request.Source, request.SourceOfTruth, request.AppSource, utils.TimestampProtoToTimePtr(request.StartedAt), utils.TimestampProtoToTimePtr(request.EndedAt), utils.TimestampProtoToTimePtr(request.CreatedAt))
	if err := jobRoleService.jobRoleCommands.CreateJobRoleCommand.Handle(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to create job role: %w", err)
	}

	jobRoleService.log.Infof("(Created New Job Role): {%s}", objectID)
	return &jobrolepb.JobRoleIdGrpcResponse{
		Id: objectID,
	}, nil
}
