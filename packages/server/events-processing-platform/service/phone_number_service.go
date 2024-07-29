package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"strings"
)

type phoneNumberService struct {
	phonenumberpb.UnimplementedPhoneNumberGrpcServiceServer
	log                 logger.Logger
	neo4jRepositories   *neo4jrepository.Repositories
	phoneNumberCommands *command_handler.CommandHandlers
	services            *Services
}

func NewPhoneNumberService(log logger.Logger, neo4jRepositories *neo4jrepository.Repositories, phoneNumberCommands *command_handler.CommandHandlers, services *Services) *phoneNumberService {
	return &phoneNumberService{
		log:                 log,
		neo4jRepositories:   neo4jRepositories,
		phoneNumberCommands: phoneNumberCommands,
		services:            services,
	}
}

func (s *phoneNumberService) UpsertPhoneNumber(ctx context.Context, request *phonenumberpb.UpsertPhoneNumberGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "PhoneNumberService.UpsertPhoneNumber")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	phoneNumberId := strings.TrimSpace(request.Id)
	if phoneNumberId == "" {
		phoneNumberId = utils.NewUUIDIfEmpty(phoneNumberId)
	}

	sourceFields := common.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	cmd := command.NewUpsertPhoneNumberCommand(phoneNumberId, request.Tenant, request.LoggedInUserId, request.PhoneNumber,
		sourceFields, utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err := s.phoneNumberCommands.UpsertPhoneNumber.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(UpsertPhoneNumber) tenant:{%s}, phoneNumber ID: {%s}, err: {%v}", request.Tenant, phoneNumberId, err.Error())
		return nil, s.errResponse(err)
	}

	return &phonenumberpb.PhoneNumberIdGrpcResponse{Id: phoneNumberId}, nil
}

func (s *phoneNumberService) FailPhoneNumberValidation(ctx context.Context, request *phonenumberpb.FailPhoneNumberValidationGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "PhoneNumberService.FailPhoneNumberValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.SetTag(tracing.SpanTagEntityId, request.PhoneNumberId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewPhoneNumberAggregateWithTenantAndID(request.Tenant, request.PhoneNumberId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(FailPhoneNumberValidation.HandleGRPCRequest) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &phonenumberpb.PhoneNumberIdGrpcResponse{Id: request.PhoneNumberId}, nil
}

func (s *phoneNumberService) PassPhoneNumberValidation(ctx context.Context, request *phonenumberpb.PassPhoneNumberValidationGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "PhoneNumberService.PassPhoneNumberValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.SetTag(tracing.SpanTagEntityId, request.PhoneNumberId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewPhoneNumberAggregateWithTenantAndID(request.Tenant, request.PhoneNumberId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(FailPhoneNumberValidation.HandleGRPCRequest) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &phonenumberpb.PhoneNumberIdGrpcResponse{Id: request.PhoneNumberId}, nil
}

func (s *phoneNumberService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}

func (s *phoneNumberService) RequestPhoneNumberValidation(ctx context.Context, request *phonenumberpb.RequestPhoneNumberValidationGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "PhoneNumberService.RequestPhoneNumberValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewPhoneNumberTempAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{
		SkipLoadEvents: true,
	}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RequestPhoneNumberValidation.HandleTemp) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &phonenumberpb.PhoneNumberIdGrpcResponse{Id: request.Id}, nil
}
