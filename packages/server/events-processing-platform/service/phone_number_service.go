package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	"strings"
)

type phoneNumberService struct {
	phonenumberpb.UnimplementedPhoneNumberGrpcServiceServer
	log                 logger.Logger
	repositories        *repository.Repositories
	phoneNumberCommands *command_handler.CommandHandlers
}

func NewPhoneNumberService(log logger.Logger, repositories *repository.Repositories, phoneNumberCommands *command_handler.CommandHandlers) *phoneNumberService {
	return &phoneNumberService{
		log:                 log,
		repositories:        repositories,
		phoneNumberCommands: phoneNumberCommands,
	}
}

func (s *phoneNumberService) UpsertPhoneNumber(ctx context.Context, request *phonenumberpb.UpsertPhoneNumberGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "PhoneNumberService.UpsertPhoneNumber")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	objectID := strings.TrimSpace(request.Id)
	var err error
	if objectID == "" {
		objectID, err = s.repositories.PhoneNumberRepository.GetIdIfExists(ctx, request.Tenant, request.PhoneNumber)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("(UpsertPhoneNumber) tenant:{%s}, email: {%s}, err: {%v}", request.Tenant, request.PhoneNumber, err.Error())
			return nil, s.errResponse(err)
		}
		objectID = utils.NewUUIDIfEmpty(objectID)
	}

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	cmd := command.NewUpsertPhoneNumberCommand(objectID, request.Tenant, request.LoggedInUserId, request.PhoneNumber,
		sourceFields, utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err = s.phoneNumberCommands.UpsertPhoneNumber.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(UpsertPhoneNumber) tenant:{%s}, phoneNumber ID: {%s}, err: {%v}", request.Tenant, objectID, err.Error())
		return nil, s.errResponse(err)
	}

	return &phonenumberpb.PhoneNumberIdGrpcResponse{Id: objectID}, nil
}

func (s *phoneNumberService) FailPhoneNumberValidation(ctx context.Context, request *phonenumberpb.FailPhoneNumberValidationGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "PhoneNumberService.FailPhoneNumberValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.SetTag(tracing.SpanTagEntityId, request.PhoneNumberId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewFailedPhoneNumberValidationCommand(request.PhoneNumberId, request.Tenant, request.LoggedInUserId, request.AppSource, request.PhoneNumber, request.CountryCodeA2, request.ErrorMessage)
	if err := s.phoneNumberCommands.FailedPhoneNumberValidation.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(FailPhoneNumberValidation) tenant:{%s}, phoneNumber ID: {%s}, err: %s", request.Tenant, request.PhoneNumberId, err.Error())
		return nil, s.errResponse(err)
	}

	return &phonenumberpb.PhoneNumberIdGrpcResponse{Id: request.PhoneNumberId}, nil
}

func (s *phoneNumberService) PassPhoneNumberValidation(ctx context.Context, request *phonenumberpb.PassPhoneNumberValidationGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "PhoneNumberService.PassPhoneNumberValidation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.SetTag(tracing.SpanTagEntityId, request.PhoneNumberId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewPhoneNumberValidatedCommand(request.PhoneNumberId, request.Tenant, request.LoggedInUserId, request.AppSource, request.PhoneNumber, request.E164, request.CountryCodeA2)
	if err := s.phoneNumberCommands.PhoneNumberValidated.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PhoneNumberValidated) tenant:{%s}, phoneNumber ID: {%s}, err: %s", request.Tenant, request.PhoneNumberId, err.Error())
		return nil, s.errResponse(err)
	}

	return &phonenumberpb.PhoneNumberIdGrpcResponse{Id: request.PhoneNumberId}, nil
}

func (s *phoneNumberService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
