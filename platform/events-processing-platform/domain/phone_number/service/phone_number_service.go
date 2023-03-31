package service

import (
	grpcErrors "github.com/AleksK1NG/es-microservice/pkg/grpc_errors"
	"github.com/google/uuid"
	phoneNumberGrpcService "github.com/openline-ai/openline-customer-os/platform/events-processing-common/proto/phone_number"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/commands"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/repository"
	"golang.org/x/net/context"
)

type phoneNumberService struct {
	phoneNumberGrpcService.UnimplementedPhoneNumberGrpcServiceServer
	log                        logger.Logger
	repositories               *repository.Repositories
	phoneNumberCommandsService *PhoneNumberCommandsService
}

func NewPhoneNumberService(log logger.Logger, repositories *repository.Repositories, phoneNumberCommandsService *PhoneNumberCommandsService) *phoneNumberService {
	return &phoneNumberService{
		log:                        log,
		repositories:               repositories,
		phoneNumberCommandsService: phoneNumberCommandsService,
	}
}

func (s *phoneNumberService) CreatePhoneNumber(ctx context.Context, request *phoneNumberGrpcService.CreatePhoneNumberGrpcRequest) (*phoneNumberGrpcService.CreatePhoneNumberGrpcResponse, error) {
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

	command := commands.NewCreatePhoneNumberCommand(aggregateID, request.GetTenant(), request.GetPhoneNumber())
	if err := s.phoneNumberCommandsService.Commands.CreatePhoneNumber.Handle(ctx, command); err != nil {
		//phoneNumberService.log.Errorf("(CreatePhoneNumber.Handle) phoneNumber ID: {%s}, err: {%v}", aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created PhoneNumber): {%s}", aggregateID)
	return &phoneNumberGrpcService.CreatePhoneNumberGrpcResponse{UUID: aggregateID}, nil
}

func (s *phoneNumberService) errResponse(err error) error {
	return grpcErrors.ErrResponse(err)
}
