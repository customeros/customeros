package service

import (
	"context"
	user_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/commands"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/utils"
)

type userService struct {
	user_grpc_service.UnimplementedUserGrpcServiceServer
	log          logger.Logger
	repositories *repository.Repositories
	userCommands *commands.UserCommands
}

func NewUserService(log logger.Logger, repositories *repository.Repositories, userCommands *commands.UserCommands) *userService {
	return &userService{
		log:          log,
		repositories: repositories,
		userCommands: userCommands,
	}
}

func (s *userService) UpsertUser(ctx context.Context, request *user_grpc_service.UpsertUserGrpcRequest) (*user_grpc_service.UserIdGrpcResponse, error) {
	aggregateID := request.Id

	coreFields := commands.UserCoreFields{
		Name: request.Name,
	}
	command := commands.NewUpsertUserCommand(aggregateID, request.Tenant, request.Source, request.SourceOfTruth, request.AppSource,
		coreFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.userCommands.UpsertUser.Handle(ctx, command); err != nil {
		s.log.Errorf("(UpsertSyncUser.Handle) tenant:{%s}, user ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created existing User): {%s}", aggregateID)

	return &user_grpc_service.UserIdGrpcResponse{Id: aggregateID}, nil
}

func (s *userService) LinkPhoneNumberToUser(ctx context.Context, request *user_grpc_service.LinkPhoneNumberToUserGrpcRequest) (*user_grpc_service.UserIdGrpcResponse, error) {
	aggregateID := request.UserId

	command := commands.NewLinkPhoneNumberCommand(aggregateID, request.Tenant, request.PhoneNumberId, request.Label, request.Primary)
	if err := s.userCommands.LinkPhoneNumberCommand.Handle(ctx, command); err != nil {
		s.log.Errorf("(LinkPhoneNumberToUser.Handle) tenant:{%s}, user ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked phone number {%s} to user {%s}", request.PhoneNumberId, aggregateID)

	return &user_grpc_service.UserIdGrpcResponse{Id: aggregateID}, nil
}

func (s *userService) LinkEmailToUser(ctx context.Context, request *user_grpc_service.LinkEmailToUserGrpcRequest) (*user_grpc_service.UserIdGrpcResponse, error) {
	aggregateID := request.UserId

	command := commands.NewLinkEmailCommand(aggregateID, request.Tenant, request.EmailId, request.Label, request.Primary)
	if err := s.userCommands.LinkEmailCommand.Handle(ctx, command); err != nil {
		s.log.Errorf("(LinkEmailToUser.Handle) tenant:{%s}, user ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked email {%s} to user {%s}", request.EmailId, aggregateID)

	return &user_grpc_service.UserIdGrpcResponse{Id: aggregateID}, nil
}

func (userService *userService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
