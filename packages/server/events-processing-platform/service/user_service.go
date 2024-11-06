package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go/log"
)

type userService struct {
	userpb.UnimplementedUserGrpcServiceServer
	log                logger.Logger
	userCommands       *command_handler.CommandHandlers
	userRequestHandler user.UserRequestHandler
	services           *Services
}

func NewUserService(log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config, userCommands *command_handler.CommandHandlers, services *Services) *userService {
	return &userService{
		log:                log,
		userCommands:       userCommands,
		userRequestHandler: user.NewUserRequestHandler(log, aggregateStore, cfg.Utils),
		services:           services,
	}
}

func (s *userService) UpsertUser(ctx context.Context, request *userpb.UpsertUserGrpcRequest) (*userpb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.UpsertUser")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	userInputId := utils.NewUUIDIfEmpty(request.Id)

	dataFields := models.UserDataFields{
		Name:            request.Name,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Internal:        request.Internal,
		Bot:             request.Bot,
		Test:            request.Test,
		ProfilePhotoUrl: request.ProfilePhotoUrl,
		Timezone:        request.Timezone,
	}
	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.Source = utils.StringFirstNonEmpty(sourceFields.Source, request.Source)
	sourceFields.SourceOfTruth = utils.StringFirstNonEmpty(sourceFields.SourceOfTruth, request.SourceOfTruth)
	sourceFields.AppSource = utils.StringFirstNonEmpty(sourceFields.AppSource, request.AppSource)
	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	cmd := command.NewUpsertUserCommand(userInputId, request.Tenant, request.LoggedInUserId, sourceFields, externalSystem,
		dataFields, utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err := s.userCommands.UpsertUser.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertUserCommand.Handle) tenant:{%s}, user input id:{%s}, err: %s", request.Tenant, userInputId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Upserted user {%s}", userInputId)

	return &userpb.UserIdGrpcResponse{Id: userInputId}, nil
}

func (s *userService) LinkJobRoleToUser(ctx context.Context, request *userpb.LinkJobRoleToUserGrpcRequest) (*userpb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.AddPlayerInfo")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, "")
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	aggregateID := request.UserId

	cmd := command.NewLinkJobRoleCommand(aggregateID, request.Tenant, request.JobRoleId)
	if err := s.userCommands.LinkJobRoleCommand.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(LinkJobRoleToUser.Handle) tenant:{%s}, user ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked job role {%s} to user {%s}", request.JobRoleId, aggregateID)

	return &userpb.UserIdGrpcResponse{Id: aggregateID}, nil
}

func (s *userService) LinkPhoneNumberToUser(ctx context.Context, request *userpb.LinkPhoneNumberToUserGrpcRequest) (*userpb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.LinkPhoneNumberToUser")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if _, err := s.userRequestHandler.HandleWithRetry(ctx, request.Tenant, request.UserId, false, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkPhoneNumberToUser.Handle) tenant:{%s}, user ID: {%s}, err: {%v}", request.Tenant, request.UserId, err)
		return nil, grpcerr.ErrResponse(err)
	}

	return &userpb.UserIdGrpcResponse{Id: request.UserId}, nil
}

func (s *userService) AddRole(ctx context.Context, request *userpb.AddRoleGrpcRequest) (*userpb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.AddRole")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewAddRole(request.UserId, request.Tenant, request.LoggedInUserId, request.Role)
	if err := s.userCommands.AddRole.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddRoleCommand.Handle) tenant:{%s}, user id:{%s}, role: {%s}, err: %s", request.Tenant, request.UserId, request.Role, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Added role {%s} for user {%s}", request.Role, request.UserId)

	return &userpb.UserIdGrpcResponse{Id: request.UserId}, nil
}

func (s *userService) RemoveRole(ctx context.Context, request *userpb.RemoveRoleGrpcRequest) (*userpb.UserIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "UserService.RemoveRole")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewRemoveRole(request.UserId, request.Tenant, request.LoggedInUserId, request.Role)
	if err := s.userCommands.RemoveRole.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemoveRoleCommand.Handle) tenant:{%s}, user id:{%s}, role: {%s}, err: %s", request.Tenant, request.UserId, request.Role, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Removed role {%s} from user {%s}", request.Role, request.UserId)

	return &userpb.UserIdGrpcResponse{Id: request.UserId}, nil
}

func (s *userService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
