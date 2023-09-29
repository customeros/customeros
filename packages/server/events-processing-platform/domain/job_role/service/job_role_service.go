package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	job_role_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
)

type jobRoleService struct {
	job_role_grpc_service.UnimplementedJobRoleGrpcServiceServer
	log             logger.Logger
	repositories    *repository.Repositories
	jobRoleCommands *commands.JobRoleCommands
}

func NewJobRoleService(log logger.Logger, repositories *repository.Repositories, jobRoleCommands *commands.JobRoleCommands) *jobRoleService {
	return &jobRoleService{
		log:             log,
		repositories:    repositories,
		jobRoleCommands: jobRoleCommands,
	}
}

func (jobRoleService *jobRoleService) CreateJobRole(ctx context.Context, request *job_role_grpc_service.CreateJobRoleGrpcRequest) (*job_role_grpc_service.JobRoleIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "JobRoleService.CreateJobRole")
	defer span.Finish()

	newObjectId, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new object ID: %w", err)
	}
	objectID := newObjectId.String()

	primary := false
	if request.Primary != nil {
		primary = *request.Primary
	}
	command := model.NewCreateJobRoleCommand(objectID, request.Tenant, request.JobTitle, request.Description, primary, request.Source, request.SourceOfTruth, request.AppSource, utils.TimestampProtoToTime(request.StartedAt), utils.TimestampProtoToTime(request.EndedAt), utils.TimestampProtoToTime(request.CreatedAt))
	if err := jobRoleService.jobRoleCommands.CreateJobRoleCommand.Handle(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to create job role: %w", err)
	}

	jobRoleService.log.Infof("(Created New Job Role): {%s}", objectID)
	return &job_role_grpc_service.JobRoleIdGrpcResponse{
		Id: objectID,
	}, nil
}
