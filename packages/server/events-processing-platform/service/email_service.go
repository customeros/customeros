package service

import (
	"context"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go/log"
)

type emailService struct {
	emailpb.UnimplementedEmailGrpcServiceServer
	log               logger.Logger
	neo4jRepositories *neo4jrepository.Repositories
	services          *Services
}

func NewEmailService(log logger.Logger, neo4jRepositories *neo4jrepository.Repositories, services *Services) *emailService {
	return &emailService{
		log:               log,
		neo4jRepositories: neo4jRepositories,
		services:          services,
	}
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
		return email.NewEmailTempAggregateWithTenantAndID(request.Tenant, request.Id)
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

func (s *emailService) UpdateEmailValidation(ctx context.Context, request *emailpb.EmailValidationGrpcRequest) (*emailpb.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.UpdateEmailValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.SetTag(tracing.SpanTagEntityId, request.EmailId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return email.NewEmailAggregateWithTenantAndID(request.Tenant, request.EmailId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateEmailValidation.HandleGRPCRequest) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &emailpb.EmailIdGrpcResponse{Id: request.EmailId}, nil
}

func (s *emailService) UpsertEmailV2(ctx context.Context, request *emailpb.UpsertEmailRequest) (*emailpb.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.UpsertEmailV2")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Validate email ID is present
	if request.EmailId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("emailId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return email.NewEmailAggregateWithTenantAndID(request.Tenant, request.EmailId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertEmailV2) tenant:%s, err: %s", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &emailpb.EmailIdGrpcResponse{Id: request.EmailId}, nil
}

func (s *emailService) DeleteEmail(ctx context.Context, request *emailpb.DeleteEmailRequest) (*emailpb.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.DeleteEmail")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Validate email ID is present
	if request.EmailId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("emailId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return email.NewEmailAggregateWithTenantAndID(request.Tenant, request.EmailId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertEmailV2) tenant:%s, err: %s", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &emailpb.EmailIdGrpcResponse{Id: request.EmailId}, nil
}
