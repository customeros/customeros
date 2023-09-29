package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	phone_number_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	grpcErrors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

type phoneNumberService struct {
	phone_number_grpc_service.UnimplementedPhoneNumberGrpcServiceServer
	log                 logger.Logger
	repositories        *repository.Repositories
	phoneNumberCommands *commands.PhoneNumberCommands
}

func NewPhoneNumberService(log logger.Logger, repositories *repository.Repositories, phoneNumberCommands *commands.PhoneNumberCommands) *phoneNumberService {
	return &phoneNumberService{
		log:                 log,
		repositories:        repositories,
		phoneNumberCommands: phoneNumberCommands,
	}
}

func (s *phoneNumberService) UpsertPhoneNumber(ctx context.Context, request *phone_number_grpc_service.UpsertPhoneNumberGrpcRequest) (*phone_number_grpc_service.PhoneNumberIdGrpcResponse, error) {
	aggregateID := request.Id

	command := commands.NewUpsertPhoneNumberCommand(aggregateID, request.Tenant, request.PhoneNumber, request.Source, request.SourceOfTruth, request.AppSource, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.phoneNumberCommands.UpsertPhoneNumber.Handle(ctx, command); err != nil {
		s.log.Errorf("(UpsertSyncPhoneNumber.Handle) tenant:{%s}, phoneNumber ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created existing PhoneNumber): {%s}", aggregateID)

	return &phone_number_grpc_service.PhoneNumberIdGrpcResponse{Id: aggregateID}, nil
}

// FIXME alexb finish implementation
func (s *phoneNumberService) CreatePhoneNumber(ctx context.Context, request *phone_number_grpc_service.CreatePhoneNumberGrpcRequest) (*phone_number_grpc_service.PhoneNumberIdGrpcResponse, error) {
	id, err := s.repositories.PhoneNumberRepository.GetIdIfExists(ctx, request.Tenant, request.PhoneNumber)
	if err != nil {
		return nil, s.errResponse(err)
	}

	var aggregateID string
	if id != "" {
		aggregateID = id
	} else {
		aggregateID = uuid.New().String()
	}

	// FIXME alexb if phoneNumber exists proceed with creation but return error if aggregate already exists

	// FIXME alexb re-implement
	//command := commands.NewCreatePhoneNumberCommand(aggregateID, request.GetTenant(), request.GetPhoneNumber())
	//if err := s.phoneNumberCommandsService.Commands.CreatePhoneNumber.Handle(ctx, command); err != nil {
	//	s.log.Errorf("(CreatePhoneNumber.Handle) phoneNumber ID: {%s}, err: {%v}", aggregateID, err)
	//	return nil, s.errResponse(err)
	//}

	s.log.Infof("(created PhoneNumber): {%s}", aggregateID)
	return &phone_number_grpc_service.PhoneNumberIdGrpcResponse{Id: aggregateID}, nil
}

func (s *phoneNumberService) errResponse(err error) error {
	return grpcErrors.ErrResponse(err)
}
