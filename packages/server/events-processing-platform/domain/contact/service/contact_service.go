package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	contact_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/utils"
)

type contactService struct {
	contact_grpc_service.UnimplementedContactGrpcServiceServer
	log             logger.Logger
	repositories    *repository.Repositories
	contactCommands *commands.ContactCommands
}

func NewContactService(log logger.Logger, repositories *repository.Repositories, contactCommands *commands.ContactCommands) *contactService {
	return &contactService{
		log:             log,
		repositories:    repositories,
		contactCommands: contactCommands,
	}
}

func (s *contactService) UpsertContact(ctx context.Context, request *contact_grpc_service.UpsertContactGrpcRequest) (*contact_grpc_service.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.UpsertContact")
	defer span.Finish()

	objectID := request.Id

	coreFields := commands.ContactDataFields{
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Prefix:          request.Prefix,
		Description:     request.Description,
		Timezone:        request.Timezone,
		ProfilePhotoUrl: request.ProfilePhotoUrl,
		Name:            request.Name,
	}
	command := commands.NewUpsertContactCommand(objectID, request.Tenant, request.Source, request.SourceOfTruth, request.AppSource,
		coreFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.contactCommands.UpsertContact.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertContact.Handle) tenant:%s, contactID: %s, err: {%v}", request.Tenant, objectID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Upserted contact: %s", objectID)

	return &contact_grpc_service.ContactIdGrpcResponse{Id: objectID}, nil
}

func (s *contactService) LinkPhoneNumberToContact(ctx context.Context, request *contact_grpc_service.LinkPhoneNumberToContactGrpcRequest) (*contact_grpc_service.ContactIdGrpcResponse, error) {
	command := commands.NewLinkPhoneNumberCommand(request.ContactId, request.Tenant, request.PhoneNumberId, request.Label, request.Primary)
	if err := s.contactCommands.LinkPhoneNumberCommand.Handle(ctx, command); err != nil {
		s.log.Errorf("(LinkPhoneNumberToContact.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, request.ContactId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked phone number {%s} to contact {%s}", request.PhoneNumberId, request.ContactId)

	return &contact_grpc_service.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) LinkEmailToContact(ctx context.Context, request *contact_grpc_service.LinkEmailToContactGrpcRequest) (*contact_grpc_service.ContactIdGrpcResponse, error) {
	command := commands.NewLinkEmailCommand(request.ContactId, request.Tenant, request.EmailId, request.Label, request.Primary)
	if err := s.contactCommands.LinkEmailCommand.Handle(ctx, command); err != nil {
		s.log.Errorf("(LinkEmailToContact.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, request.ContactId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked email {%s} to contact {%s}", request.EmailId, request.ContactId)

	return &contact_grpc_service.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) CreateContact(ctx context.Context, request *contact_grpc_service.CreateContactGrpcRequest) (*contact_grpc_service.CreateContactGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.CreateContact")
	defer span.Finish()

	newObjectId, err := uuid.NewUUID()
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to generate new object ID: %w", err)
	}
	objectID := newObjectId.String()

	command := commands.NewContactCreateCommand(objectID, request.Tenant, request.FirstName, request.LastName, request.Prefix, request.Description, request.Timezone, request.ProfilePhotoUrl, request.Source, request.SourceOfTruth, request.AppSource, utils.TimestampProtoToTime(request.CreatedAt))
	if err := s.contactCommands.CreateContactCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(ContactCreateCommand.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, objectID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created new Contact): {%s}", objectID)

	return &contact_grpc_service.CreateContactGrpcResponse{Id: objectID}, nil
}

func (s *contactService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
