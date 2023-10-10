package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	pb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
)

type contactService struct {
	pb.UnimplementedContactGrpcServiceServer
	log             logger.Logger
	repositories    *repository.Repositories
	contactCommands *command_handler.ContactCommands
}

func NewContactService(log logger.Logger, repositories *repository.Repositories, contactCommands *command_handler.ContactCommands) *contactService {
	return &contactService{
		log:             log,
		repositories:    repositories,
		contactCommands: contactCommands,
	}
}

func (s *contactService) UpsertContact(ctx context.Context, request *pb.UpsertContactGrpcRequest) (*pb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.UpsertContact")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request.contactId", request.Id))

	contactId := utils.NewUUIDIfEmpty(request.Id)

	dataFields := models.ContactDataFields{
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Name:            request.Name,
		Description:     request.Description,
		Prefix:          request.Prefix,
		Timezone:        request.Timezone,
		ProfilePhotoUrl: request.ProfilePhotoUrl,
	}
	sourceFields := cmnmod.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.Source = utils.StringFirstNonEmpty(sourceFields.Source, request.Source)
	sourceFields.SourceOfTruth = utils.StringFirstNonEmpty(sourceFields.SourceOfTruth, request.SourceOfTruth)
	sourceFields.AppSource = utils.StringFirstNonEmpty(sourceFields.AppSource, request.AppSource)

	externalSystem := cmnmod.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	cmd := command.NewUpsertContactCommand(contactId, request.Tenant, request.LoggedInUserId, sourceFields, externalSystem,
		dataFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt), request.Id == "")
	if err := s.contactCommands.UpsertContact.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertContact.Handle) tenant:%s, contactID: %s, err: {%v}", request.Tenant, contactId, err)
		return nil, s.errResponse(err)
	}

	return &pb.ContactIdGrpcResponse{Id: contactId}, nil
}

func (s *contactService) LinkPhoneNumberToContact(ctx context.Context, request *pb.LinkPhoneNumberToContactGrpcRequest) (*pb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.LinkPhoneNumberToContact")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("phoneNumberId", request.PhoneNumberId), log.String("contactId", request.ContactId))

	cmd := command.NewLinkPhoneNumberCommand(request.ContactId, request.Tenant, request.LoggedInUserId, request.PhoneNumberId, request.Label, request.Primary)
	if err := s.contactCommands.LinkPhoneNumberCommand.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(LinkPhoneNumberCommand.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, request.ContactId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) LinkEmailToContact(ctx context.Context, request *pb.LinkEmailToContactGrpcRequest) (*pb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.LinkEmailToContact")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("emailId", request.EmailId), log.String("contactId", request.ContactId))

	cmd := command.NewLinkEmailCommand(request.ContactId, request.Tenant, request.LoggedInUserId, request.EmailId, request.Label, request.Primary)
	if err := s.contactCommands.LinkEmailCommand.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(LinkEmailCommand.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, request.ContactId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) LinkLocationToContact(ctx context.Context, request *pb.LinkLocationToContactGrpcRequest) (*pb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.LinkLocationToContact")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("locationId", request.LocationId), log.String("contactId", request.ContactId))

	cmd := command.NewLinkLocationCommand(request.ContactId, request.Tenant, request.LoggedInUserId, request.LocationId)
	if err := s.contactCommands.LinkLocationCommand.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(LinkLocationCommand.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, request.ContactId, err.Error())
		return nil, s.errResponse(err)
	}

	return &pb.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
