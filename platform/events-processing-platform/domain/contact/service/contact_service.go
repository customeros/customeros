package service

import (
	"context"
	grpcErrors "github.com/AleksK1NG/es-microservice/pkg/grpc_errors"
	contactGrpcService "github.com/openline-ai/openline-customer-os/platform/events-processing-common/proto/contact"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/contact/commands"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
	uuid "github.com/satori/go.uuid"
)

type contactService struct {
	contactGrpcService.UnimplementedContactGrpcServiceServer
	log                    logger.Logger
	contactCommandsService *ContactCommandsService
}

func NewContactService(log logger.Logger, contactCommandsService *ContactCommandsService) *contactService {
	return &contactService{log: log, contactCommandsService: contactCommandsService}
}

func (contactService *contactService) CreateContact(ctx context.Context, request *contactGrpcService.CreateContactGrpcRequest) (*contactGrpcService.CreateContactGrpcResponse, error) {
	/*ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "contactService.CreateContact")
	defer span.Finish()
	span.LogFields(log.String("request", request.String()))
	*/

	aggregateID := uuid.NewV4().String()
	command := commands.NewCreateContactCommand(aggregateID, request.GetUUID(), request.GetFirstName(), request.GetLastName())

	// add validation

	if err := contactService.contactCommandsService.Commands.CreateContact.Handle(ctx, command); err != nil {
		contactService.log.Errorf("(CreateContact.Handle) contact ID: {%s}, err: {%v}", aggregateID, err)
		return nil, contactService.errResponse(err)
	}

	contactService.log.Infof("(created contact): contact: {%s}", aggregateID)
	return &contactGrpcService.CreateContactGrpcResponse{AggregateID: aggregateID}, nil
}

func (contactService *contactService) errResponse(err error) error {
	return grpcErrors.ErrResponse(err)
}
