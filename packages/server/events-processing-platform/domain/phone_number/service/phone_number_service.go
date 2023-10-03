package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	phone_number_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	grpcErrors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"strings"
)

type phoneNumberService struct {
	phone_number_grpc_service.UnimplementedPhoneNumberGrpcServiceServer
	log                 logger.Logger
	repositories        *repository.Repositories
	phoneNumberCommands *command_handler.PhoneNumberCommands
}

func NewPhoneNumberService(log logger.Logger, repositories *repository.Repositories, phoneNumberCommands *command_handler.PhoneNumberCommands) *phoneNumberService {
	return &phoneNumberService{
		log:                 log,
		repositories:        repositories,
		phoneNumberCommands: phoneNumberCommands,
	}
}

func (s *phoneNumberService) UpsertPhoneNumber(ctx context.Context, request *phone_number_grpc_service.UpsertPhoneNumberGrpcRequest) (*phone_number_grpc_service.PhoneNumberIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "PhoneNumberService.UpsertPhoneNumber")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)

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

	sourceFields := common_models.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	cmd := command.NewUpsertPhoneNumberCommand(objectID, request.Tenant, request.LoggedInUserId, request.PhoneNumber,
		sourceFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err = s.phoneNumberCommands.UpsertPhoneNumber.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(UpsertPhoneNumber) tenant:{%s}, phoneNumber ID: {%s}, err: {%v}", request.Tenant, objectID, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("(UpsertPhoneNumber): {%s}", objectID)

	return &phone_number_grpc_service.PhoneNumberIdGrpcResponse{Id: objectID}, nil
}

func (s *phoneNumberService) errResponse(err error) error {
	return grpcErrors.ErrResponse(err)
}
