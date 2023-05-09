package service

import (
	"context"
	contact_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
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
	aggregateID := request.Id

	coreFields := commands.ContactCoreFields{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Name:      request.Name,
		Prefix:    request.Prefix,
	}
	command := commands.NewUpsertContactCommand(aggregateID, request.Tenant, request.Source, request.SourceOfTruth, request.AppSource,
		coreFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.contactCommands.UpsertContact.Handle(ctx, command); err != nil {
		s.log.Errorf("(UpsertSyncContact.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created existing Contact): {%s}", aggregateID)

	return &contact_grpc_service.ContactIdGrpcResponse{Id: aggregateID}, nil
}

func (s *contactService) LinkPhoneNumberToContact(ctx context.Context, request *contact_grpc_service.LinkPhoneNumberToContactGrpcRequest) (*contact_grpc_service.ContactIdGrpcResponse, error) {
	aggregateID := request.ContactId

	command := commands.NewLinkPhoneNumberCommand(aggregateID, request.Tenant, request.PhoneNumberId, request.Label, request.Primary)
	if err := s.contactCommands.LinkPhoneNumberCommand.Handle(ctx, command); err != nil {
		s.log.Errorf("(LinkPhoneNumberToContact.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked phone number {%s} to contact {%s}", request.PhoneNumberId, aggregateID)

	return &contact_grpc_service.ContactIdGrpcResponse{Id: aggregateID}, nil
}

func (s *contactService) LinkEmailToContact(ctx context.Context, request *contact_grpc_service.LinkEmailToContactGrpcRequest) (*contact_grpc_service.ContactIdGrpcResponse, error) {
	aggregateID := request.ContactId

	command := commands.NewLinkEmailCommand(aggregateID, request.Tenant, request.EmailId, request.Label, request.Primary)
	if err := s.contactCommands.LinkEmailCommand.Handle(ctx, command); err != nil {
		s.log.Errorf("(LinkEmailToContact.Handle) tenant:{%s}, contact ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked email {%s} to contact {%s}", request.EmailId, aggregateID)

	return &contact_grpc_service.ContactIdGrpcResponse{Id: aggregateID}, nil
}

// FIXME alexb implement correctly
//func (contactService *contactService) CreateContact(ctx context.Context, request *contact_grpc_service.CreateContactGrpcRequest) (*contact_grpc_service.CreateContactGrpcResponse, error) {
//	/*ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "contactService.CreateContact")
//	defer span.Finish()
//	span.LogFields(log.String("request", request.String()))
//	*/
//
//	aggregateID := uuid.NewV4().String()
//	command := commands.NewCreateContactCommand(aggregateID, request.GetUUID(), request.GetFirstName(), request.GetLastName())
//
//	// add validation
//
//	if err := contactService.contactCommandsService.Commands.CreateContact.Handle(ctx, command); err != nil {
//		contactService.log.Errorf("(CreateContact.Handle) contact ID: {%s}, err: {%v}", aggregateID, err)
//		return nil, contactService.errResponse(err)
//	}
//
//	contactService.log.Infof("(created contact): contact: {%s}", aggregateID)
//	return &contact_grpc_service.CreateContactGrpcResponse{AggregateID: aggregateID}, nil
//}

func (contactService *contactService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
