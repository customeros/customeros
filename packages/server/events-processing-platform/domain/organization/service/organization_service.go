package service

import (
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/commands"
	grpc_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/utils"
	"golang.org/x/net/context"
)

type organizationService struct {
	organization_grpc_service.UnimplementedOrganizationGrpcServiceServer
	log                  logger.Logger
	repositories         *repository.Repositories
	organizationCommands *commands.OrganizationCommands
}

func NewOrganizationService(log logger.Logger, repositories *repository.Repositories, organizationCommands *commands.OrganizationCommands) *organizationService {
	return &organizationService{
		log:                  log,
		repositories:         repositories,
		organizationCommands: organizationCommands,
	}
}

func (s *organizationService) UpsertOrganization(ctx context.Context, request *organization_grpc_service.UpsertOrganizationGrpcRequest) (*organization_grpc_service.OrganizationIdGrpcResponse, error) {
	aggregateID := request.Id

	coreFields := commands.OrganizationCoreFields{
		Name: request.Name,
	}
	command := commands.NewUpsertOrganizationCommand(aggregateID, request.Tenant, request.Source, request.SourceOfTruth, request.AppSource,
		coreFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.organizationCommands.UpsertOrganization.Handle(ctx, command); err != nil {
		s.log.Errorf("(UpsertSyncOrganization.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created existing Organization): {%s}", aggregateID)

	return &organization_grpc_service.OrganizationIdGrpcResponse{Id: aggregateID}, nil
}

func (s *organizationService) LinkPhoneNumberToOrganization(ctx context.Context, request *organization_grpc_service.LinkPhoneNumberToOrganizationGrpcRequest) (*organization_grpc_service.OrganizationIdGrpcResponse, error) {
	aggregateID := request.OrganizationId

	command := commands.NewLinkPhoneNumberCommand(aggregateID, request.Tenant, request.PhoneNumberId, request.Label, request.Primary)
	if err := s.organizationCommands.LinkPhoneNumberCommand.Handle(ctx, command); err != nil {
		s.log.Errorf("(LinkPhoneNumberToOrganization.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked phone number {%s} to organization {%s}", request.PhoneNumberId, aggregateID)

	return &organization_grpc_service.OrganizationIdGrpcResponse{Id: aggregateID}, nil
}

func (s *organizationService) LinkEmailToOrganization(ctx context.Context, request *organization_grpc_service.LinkEmailToOrganizationGrpcRequest) (*organization_grpc_service.OrganizationIdGrpcResponse, error) {
	aggregateID := request.OrganizationId

	command := commands.NewLinkEmailCommand(aggregateID, request.Tenant, request.EmailId, request.Label, request.Primary)
	if err := s.organizationCommands.LinkEmailCommand.Handle(ctx, command); err != nil {
		s.log.Errorf("(LinkEmailToOrganization.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked email {%s} to organization {%s}", request.EmailId, aggregateID)

	return &organization_grpc_service.OrganizationIdGrpcResponse{Id: aggregateID}, nil
}

func (organizationService *organizationService) errResponse(err error) error {
	return grpc_errors.ErrResponse(err)
}
