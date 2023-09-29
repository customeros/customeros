package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
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
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.UpsertEmail")
	defer span.Finish()

	objectID := request.Id
	var err error
	if len(objectID) == 0 {
		objectID, err = s.repositories.EmailRepository.GetIdIfExists(ctx, request.Tenant, request.RawEmail)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("(UpsertSyncEmail.Handle) tenant:{%s}, email: {%s}, err: {%v}", request.Tenant, request.RawEmail, err)
			return nil, s.errResponse(err)
		}
		if len(objectID) == 0 {
			newId, err := uuid.NewUUID()
			if err != nil {
				s.log.Errorf("(UpsertSyncEmail.Handle) tenant:{%s}, email: {%s}, err: {%v}", request.Tenant, request.RawEmail, err)
				return nil, s.errResponse(err)
			}
			objectID = newId.String()
		}
	}

	command := commands.NewUpsertEmailCommand(objectID, request.Tenant, request.RawEmail, request.Source, request.SourceOfTruth, request.AppSource, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.emailCommands.UpsertEmail.Handle(ctx, command); err != nil {
		s.log.Errorf("(UpsertSyncEmail.Handle) tenant:{%s}, email ID: {%s}, err: {%v}", request.Tenant, objectID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created existing Email): {%s}", objectID)

	return &email_grpc_service.EmailIdGrpcResponse{Id: objectID}, nil
}

func (s *emailService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
