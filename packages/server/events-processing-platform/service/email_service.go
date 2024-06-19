package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"strings"
)

type emailService struct {
	emailpb.UnimplementedEmailGrpcServiceServer
	log                  logger.Logger
	neo4jRepositories    *neo4jrepository.Repositories
	emailCommandHandlers *command_handler.CommandHandlers
	services             *Services
}

func NewEmailService(log logger.Logger, neo4jRepositories *neo4jrepository.Repositories, emailCommandHandlers *command_handler.CommandHandlers, services *Services) *emailService {
	return &emailService{
		log:                  log,
		neo4jRepositories:    neo4jRepositories,
		emailCommandHandlers: emailCommandHandlers,
		services:             services,
	}
}

func (s *emailService) UpsertEmail(ctx context.Context, request *emailpb.UpsertEmailGrpcRequest) (*emailpb.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.UpsertEmail")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	emailId := strings.TrimSpace(request.Id)
	if emailId == "" {
		emailId = utils.NewUUIDIfEmpty(emailId)
	}

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	cmd := command.NewUpsertEmailCommand(emailId, request.Tenant, request.LoggedInUserId, request.RawEmail, sourceFields,
		utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err := s.emailCommandHandlers.Upsert.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(UpsertSyncEmail.Handle) tenant:{%s}, email ID: {%s}, err: {%v}", request.Tenant, emailId, err)
		return nil, s.errResponse(err)
	}

	return &emailpb.EmailIdGrpcResponse{Id: emailId}, nil
}

func (s *emailService) FailEmailValidation(ctx context.Context, request *emailpb.FailEmailValidationGrpcRequest) (*emailpb.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.FailEmailValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.SetTag(tracing.SpanTagEntityId, request.EmailId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewFailedEmailValidationCommand(request.EmailId, request.Tenant, request.LoggedInUserId, request.AppSource, request.ErrorMessage)
	if err := s.emailCommandHandlers.FailEmailValidation.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(FailEmailValidation) tenant:{%s}, email ID: {%s}, err: {%v}", request.Tenant, request.EmailId, err.Error())
		return nil, s.errResponse(err)
	}

	return &emailpb.EmailIdGrpcResponse{Id: request.EmailId}, nil
}

func (s *emailService) PassEmailValidation(ctx context.Context, request *emailpb.PassEmailValidationGrpcRequest) (*emailpb.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.PassEmailValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.SetTag(tracing.SpanTagEntityId, request.EmailId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewEmailValidatedCommand(request.EmailId, request.Tenant, request.LoggedInUserId, request.AppSource, request.RawEmail,
		request.IsReachable, request.ErrorMessage, request.Domain, request.Username, request.Email, request.AcceptsMail, request.CanConnectSmtp,
		request.HasFullInbox, request.IsCatchAll, request.IsDisabled, request.IsValidSyntax)
	if err := s.emailCommandHandlers.EmailValidated.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(EmailValidated) tenant:{%s}, email ID: {%s}, err: %s", request.Tenant, request.EmailId, err.Error())
		return nil, s.errResponse(err)
	}

	return &emailpb.EmailIdGrpcResponse{Id: request.EmailId}, nil
}

func (s *emailService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}

func (s *emailService) RequestEmailValidation(ctx context.Context, request *emailpb.RequestEmailValidationGrpcRequest) (*emailpb.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.RequestEmailValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewEmailTempAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{
		SkipLoadEvents: true,
	}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RequestEmailValidation.HandleTemp) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &emailpb.EmailIdGrpcResponse{Id: request.Id}, nil
}
