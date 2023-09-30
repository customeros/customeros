package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	pb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type userService struct {
	pb.UnimplementedUserGrpcServiceServer
	log          logger.Logger
	repositories *repository.Repositories
	userCommands *command_handler.UserCommands
}

func NewUserService(log logger.Logger, repositories *repository.Repositories, userCommands *command_handler.UserCommands) *userService {
	return &userService{
		log:          log,
		repositories: repositories,
		userCommands: userCommands,
	}
}

func (s *userService) UpsertUser(ctx context.Context, request *pb.UpsertUserGrpcRequest) (*pb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.UpsertUser")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("userRequestId", request.Id))

	userInputId := request.Id
	if strings.TrimSpace(userInputId) == "" {
		userInputId = uuid.New().String()
	}

	dataFields := models.UserDataFields{
		Name:            request.Name,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Internal:        request.Internal,
		ProfilePhotoUrl: request.ProfilePhotoUrl,
		Timezone:        request.Timezone,
	}
	sourceFields := common_models.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	if sourceFields.Source == "" && request.Source != "" {
		sourceFields.Source = request.Source
	}
	if sourceFields.SourceOfTruth == "" && request.SourceOfTruth != "" {
		sourceFields.SourceOfTruth = request.SourceOfTruth
	}
	if sourceFields.AppSource == "" && request.AppSource != "" {
		sourceFields.AppSource = request.AppSource
	}
	externalSystem := common_models.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	cmd := command.NewUpsertUserCommand(userInputId, request.Tenant, request.LoggedInUserId, sourceFields, externalSystem,
		dataFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.userCommands.UpsertUser.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertUserCommand.Handle) tenant:{%s}, user input id:{%s}, err: %s", request.Tenant, userInputId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Upserted user {%s}", userInputId)

	return &pb.UserIdGrpcResponse{Id: userInputId}, nil
}

func (s *userService) AddPlayerInfo(ctx context.Context, request *pb.AddPlayerInfoGrpcRequest) (*pb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.AddPlayerInfo")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)

	sourceFields := common_models.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	cmd := command.NewAddPlayerInfoCommand(request.UserId, request.Tenant, request.LoggedInUserId, sourceFields,
		request.Provider, request.AuthId, request.IdentityId, utils.TimestampProtoToTime(request.Timestamp))
	if err := s.userCommands.AddPlayerInfo.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddPlayerInfoCommand.Handle) tenant:{%s}, user input id:{%s}, err: %s", request.Tenant, request.UserId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Added player info to user {%s}", request.UserId)

	return &pb.UserIdGrpcResponse{Id: request.UserId}, nil
}

func (s *userService) LinkJobRoleToUser(ctx context.Context, request *pb.LinkJobRoleToUserGrpcRequest) (*pb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.AddPlayerInfo")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, "")

	aggregateID := request.UserId

	cmd := command.NewLinkJobRoleCommand(aggregateID, request.Tenant, request.JobRoleId)
	if err := s.userCommands.LinkJobRoleCommand.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(LinkJobRoleToUser.Handle) tenant:{%s}, user ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked job role {%s} to user {%s}", request.JobRoleId, aggregateID)

	return &pb.UserIdGrpcResponse{Id: aggregateID}, nil
}

func (s *userService) LinkPhoneNumberToUser(ctx context.Context, request *pb.LinkPhoneNumberToUserGrpcRequest) (*pb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.LinkPhoneNumberToUser")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("phoneNumberId", request.PhoneNumberId))

	objectId := request.UserId

	cmd := command.NewLinkPhoneNumberCommand(objectId, request.Tenant, request.LoggedInUserId, request.PhoneNumberId, request.Label, request.Primary)
	if err := s.userCommands.LinkPhoneNumberCommand.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(LinkPhoneNumberToUser.Handle) tenant:{%s}, user ID: {%s}, err: {%v}", request.Tenant, objectId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked phone number {%s} to user {%s}", request.PhoneNumberId, objectId)

	return &pb.UserIdGrpcResponse{Id: objectId}, nil
}

func (s *userService) LinkEmailToUser(ctx context.Context, request *pb.LinkEmailToUserGrpcRequest) (*pb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.LinkPhoneNumberToUser")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("emailId", request.EmailId))

	aggregateID := request.UserId

	cmd := command.NewLinkEmailCommand(aggregateID, request.Tenant, request.LoggedInUserId, request.EmailId, request.Label, request.Primary)
	if err := s.userCommands.LinkEmailCommand.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(LinkEmailToUser.Handle) tenant:{%s}, user ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked email {%s} to user {%s}", request.EmailId, aggregateID)

	return &pb.UserIdGrpcResponse{Id: aggregateID}, nil
}

func (userService *userService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
