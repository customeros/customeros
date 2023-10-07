package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"strings"
)

type emailService struct {
	email_grpc_service.UnimplementedEmailGrpcServiceServer
	log           logger.Logger
	repositories  *repository.Repositories
	emailCommands *command_handler.EmailCommands
}

func NewEmailService(log logger.Logger, repositories *repository.Repositories, emailCommands *command_handler.EmailCommands) *emailService {
	return &emailService{
		log:           log,
		repositories:  repositories,
		emailCommands: emailCommands,
	}
}

func (s *emailService) UpsertEmail(ctx context.Context, request *email_grpc_service.UpsertEmailGrpcRequest) (*email_grpc_service.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.UpsertEmail")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)

	emailId := strings.TrimSpace(request.Id)
	var err error
	if emailId == "" {
		emailId, err = s.repositories.EmailRepository.GetIdIfExists(ctx, request.Tenant, request.RawEmail)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("(UpsertEmail) tenant:{%s}, email: {%s}, err: {%v}", request.Tenant, request.RawEmail, err)
			return nil, s.errResponse(err)
		}
		emailId = utils.NewUUIDIfEmpty(emailId)
	}

	sourceFields := common_models.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	cmd := command.NewUpsertEmailCommand(emailId, request.Tenant, request.LoggedInUserId, request.RawEmail, sourceFields,
		utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.emailCommands.UpsertEmail.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(UpsertSyncEmail.Handle) tenant:{%s}, email ID: {%s}, err: {%v}", request.Tenant, emailId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(UpsertEmail): {%s}", request.RawEmail)

	return &email_grpc_service.EmailIdGrpcResponse{Id: emailId}, nil
}

func (s *emailService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
