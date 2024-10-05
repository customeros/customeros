package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
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

func (s *emailService) UpsertEmail(ctx context.Context, request *emailpb.UpsertEmailGrpcRequest) (*emailpb.EmailIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "EmailService.UpsertEmail")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	emailId := utils.NewUUIDIfEmpty(request.Id)

	sourceFields := common.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	emailAggregate, err := email.LoadEmailAggregate(ctx, s.services.es, request.Tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)

	}

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), createdAtNotNil)

	var evt eventstore.Event

	if eventstore.IsAggregateNotFound(emailAggregate) {
		evt, err = event.NewEmailCreateEvent(emailAggregate, request.Tenant, request.RawEmail, sourceFields, createdAtNotNil, updatedAtNotNil, request.LinkWithType, request.LinkWithId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, s.errResponse(err)
		}
	} else {
		evt, err = event.NewEmailUpdateEvent(emailAggregate, request.Tenant, request.RawEmail, sourceFields.Source, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, s.errResponse(err)
		}
	}

	eventstore.EnrichEventWithMetadataExtended(&evt, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	err = emailAggregate.Apply(evt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	err = s.services.es.Save(ctx, emailAggregate)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	return &emailpb.EmailIdGrpcResponse{Id: emailId}, nil
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
