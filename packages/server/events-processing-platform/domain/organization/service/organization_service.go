package service

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/utils"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type organizationService struct {
	pb.UnimplementedOrganizationGrpcServiceServer
	log                  logger.Logger
	repositories         *repository.Repositories
	organizationCommands *command_handler.OrganizationCommands
}

func NewOrganizationService(log logger.Logger, repositories *repository.Repositories, organizationCommands *command_handler.OrganizationCommands) *organizationService {
	return &organizationService{
		log:                  log,
		repositories:         repositories,
		organizationCommands: organizationCommands,
	}
}

func (s *organizationService) UpsertOrganization(ctx context.Context, request *pb.UpsertOrganizationGrpcRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpsertOrganization")
	defer span.Finish()

	organizationId := request.Id
	if organizationId == "" {
		organizationId = uuid.New().String()
	}

	coreFields := models.OrganizationDataFields{
		Name:              request.Name,
		Hide:              request.Hide,
		Description:       request.Description,
		Website:           request.Website,
		Industry:          request.Industry,
		SubIndustry:       request.SubIndustry,
		IndustryGroup:     request.IndustryGroup,
		TargetAudience:    request.TargetAudience,
		ValueProposition:  request.ValueProposition,
		IsPublic:          request.IsPublic,
		Employees:         request.Employees,
		Market:            request.Market,
		LastFundingRound:  request.LastFundingRound,
		LastFundingAmount: request.LastFundingAmount,
		Note:              request.Note,
	}
	command := cmd.NewUpsertOrganizationCommand(organizationId, request.Tenant, request.Source, request.SourceOfTruth, request.AppSource, request.UserId, coreFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.organizationCommands.UpsertOrganization.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertSyncOrganization.Handle) tenant:%s, organizationID: %s , err: {%v}", request.Tenant, organizationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Upserted organization %s", organizationId)

	return &pb.OrganizationIdGrpcResponse{Id: organizationId}, nil
}

func (s *organizationService) LinkPhoneNumberToOrganization(ctx context.Context, request *pb.LinkPhoneNumberToOrganizationGrpcRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkPhoneNumberToOrganization")
	defer span.Finish()

	command := cmd.NewLinkPhoneNumberCommand(request.OrganizationId, request.Tenant, request.PhoneNumberId, request.Label, request.Primary)
	if err := s.organizationCommands.LinkPhoneNumberCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkPhoneNumberToOrganization.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked phone number {%s} to organization {%s}", request.PhoneNumberId, request.OrganizationId)

	return &pb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) LinkEmailToOrganization(ctx context.Context, request *pb.LinkEmailToOrganizationGrpcRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkEmailToOrganization")
	defer span.Finish()

	command := cmd.NewLinkEmailCommand(request.OrganizationId, request.Tenant, request.EmailId, request.Label, request.Primary)
	if err := s.organizationCommands.LinkEmailCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked email {%s} to organization {%s}", request.EmailId, request.OrganizationId)

	return &pb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) LinkDomainToOrganization(ctx context.Context, request *pb.LinkDomainToOrganizationGrpcRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkDomainToOrganization")
	defer span.Finish()

	command := cmd.NewLinkDomainCommand(request.OrganizationId, request.Tenant, request.Domain, request.UserId)
	if err := s.organizationCommands.LinkDomainCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked domain {%s} to organization {%s}", request.Domain, request.OrganizationId)

	return &pb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) UpdateOrganizationRenewalLikelihood(ctx context.Context, req *pb.OrganizationRenewalLikelihoodRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateOrganizationRenewalLikelihood")
	defer span.Finish()
	span.LogFields(log.Object("request", req))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	fields := models.RenewalLikelihoodFields{
		RenewalLikelihood: mapper.MapRenewalLikelihoodToModels(req.Likelihood),
		Comment:           req.Comment,
		UpdatedBy:         req.UserId,
	}
	command := cmd.NewUpdateRenewalLikelihoodCommand(req.Tenant, req.OrganizationId, fields)
	if err := s.organizationCommands.UpdateRenewalLikelihoodCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed update renewal likelihood for tenant: %s organizationID: %s, err: %s", req.Tenant, req.OrganizationId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Updated renewal likelihood for tenant:%s organizationID: %s", req.Tenant, req.OrganizationId)

	return &pb.OrganizationIdGrpcResponse{Id: req.OrganizationId}, nil
}

func (s *organizationService) UpdateOrganizationRenewalForecast(ctx context.Context, req *pb.OrganizationRenewalForecastRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateOrganizationRenewalForecast")
	defer span.Finish()
	span.LogFields(log.Object("request", req))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	fields := models.RenewalForecastFields{
		Amount:    req.Amount,
		Comment:   req.Comment,
		UpdatedBy: req.UserId,
	}
	command := cmd.NewUpdateRenewalForecastCommand(req.Tenant, req.OrganizationId, fields, "")
	if err := s.organizationCommands.UpdateRenewalForecastCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed update renewal forecast for tenant: %s organizationID: %s, err: %s", req.Tenant, req.OrganizationId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Updated renewal forecast for tenant:%s organizationID: %s", req.Tenant, req.OrganizationId)

	return &pb.OrganizationIdGrpcResponse{Id: req.OrganizationId}, nil
}

func (s *organizationService) UpdateOrganizationBillingDetails(ctx context.Context, req *pb.OrganizationBillingDetailsRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateOrganizationBillingDetails")
	defer span.Finish()
	span.LogFields(log.Object("request", req))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	fields := models.BillingDetailsFields{
		Amount:            req.Amount,
		UpdatedBy:         req.UserId,
		Frequency:         mapper.MapFrequencyToString(req.Frequency),
		RenewalCycle:      mapper.MapFrequencyToString(req.RenewalCycle),
		RenewalCycleStart: utils.TimestampProtoToTime(req.CycleStart),
	}
	command := cmd.NewUpdateBillingDetailsCommand(req.Tenant, req.OrganizationId, fields)
	if err := s.organizationCommands.UpdateBillingDetailsCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed update billing details for tenant: %s organizationID: %s, err: %s", req.Tenant, req.OrganizationId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Updated billing details for tenant:%s organizationID: %s", req.Tenant, req.OrganizationId)

	return &pb.OrganizationIdGrpcResponse{Id: req.OrganizationId}, nil
}

func (s *organizationService) RequestRenewNextCycleDate(ctx context.Context, req *pb.RequestRenewNextCycleDateRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RequestRenewNextCycleDate")
	defer span.Finish()
	span.LogFields(log.Object("request", req))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := cmd.NewRequestNextCycleDateCommand(req.Tenant, req.OrganizationId)
	if err := s.organizationCommands.RequestNextCycleDateCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed request next cycle date for tenant: %s organizationID: %s, err: %s", req.Tenant, req.OrganizationId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Requested next cycle date renewal for tenant:%s organizationID: %s", req.Tenant, req.OrganizationId)

	return &pb.OrganizationIdGrpcResponse{Id: req.OrganizationId}, nil
}

func (s *organizationService) HideOrganization(ctx context.Context, req *pb.OrganizationIdGrpcRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.HideOrganization")
	defer span.Finish()
	span.LogFields(log.Object("request", req))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := cmd.NewHideOrganizationCommand(req.Tenant, req.OrganizationId, req.UserId)
	if err := s.organizationCommands.HideOrganizationCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed hide organization with id  %s for tenant %s, err: %s", req.OrganizationId, req.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Hidden organization with id %s for tenant %s", req.OrganizationId, req.Tenant)

	return &pb.OrganizationIdGrpcResponse{Id: req.OrganizationId}, nil
}

func (s *organizationService) ShowOrganization(ctx context.Context, req *pb.OrganizationIdGrpcRequest) (*pb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.ShowOrganization")
	defer span.Finish()
	span.LogFields(log.Object("request", req))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := cmd.NewShowOrganizationCommand(req.Tenant, req.OrganizationId, req.UserId)
	if err := s.organizationCommands.ShowOrganizationCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed show organization with id  %s for tenant %s, err: %s", req.OrganizationId, req.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Show organization with id %s for tenant %s", req.OrganizationId, req.Tenant)

	return &pb.OrganizationIdGrpcResponse{Id: req.OrganizationId}, nil
}

func (s *organizationService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
