package service

import (
	"context"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	email_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/errors"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/utils"
)

type emailService struct {
	email_grpc_service.UnimplementedEmailGrpcServiceServer
	log           logger.Logger
	repositories  *repository.Repositories
	emailCommands *commands.EmailCommands
}

func NewEmailService(log logger.Logger, repositories *repository.Repositories, emailCommands *commands.EmailCommands) *emailService {
	return &emailService{
		log:           log,
		repositories:  repositories,
		emailCommands: emailCommands,
	}
}

func (s *emailService) UpsertEmail(ctx context.Context, request *email_grpc_service.UpsertEmailGrpcRequest) (*email_grpc_service.EmailIdGrpcResponse, error) {
	aggregateID := request.Id

	if len(aggregateID) == 0 {
		return &email_grpc_service.EmailIdGrpcResponse{}, email_errors.ErrEmailMissingId
	}

	command := commands.NewUpsertEmailCommand(aggregateID, request.Tenant, request.RawEmail, request.Source, request.SourceOfTruth, request.AppSource, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.emailCommands.UpsertEmail.Handle(ctx, command); err != nil {
		s.log.Errorf("(UpsertSyncEmail.Handle) tenant:{%s}, email ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created existing Email): {%s}", aggregateID)

	return &email_grpc_service.EmailIdGrpcResponse{Id: aggregateID}, nil
}

func (s *emailService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
