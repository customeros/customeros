package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type organizationService struct {
	organizationpb.UnimplementedOrganizationGrpcServiceServer
	log                  logger.Logger
	repositories         *repository.Repositories
	organizationCommands *command_handler.CommandHandlers
}

func NewOrganizationService(log logger.Logger, repositories *repository.Repositories, organizationCommands *command_handler.CommandHandlers) *organizationService {
	return &organizationService{
		log:                  log,
		repositories:         repositories,
		organizationCommands: organizationCommands,
	}
}

func (s *organizationService) UpsertOrganization(ctx context.Context, request *organizationpb.UpsertOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpsertOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.LoggedInUserId))
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	organizationId := utils.NewUUIDIfEmpty(request.Id)

	dataFields := models.OrganizationDataFields{
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
		IsCustomer:        request.IsCustomer,
		Employees:         request.Employees,
		Market:            request.Market,
		LastFundingRound:  request.LastFundingRound,
		LastFundingAmount: request.LastFundingAmount,
		ReferenceId:       request.ReferenceId,
		Note:              request.Note,
	}
	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.Source = utils.StringFirstNonEmpty(sourceFields.Source, request.Source)
	sourceFields.SourceOfTruth = utils.StringFirstNonEmpty(sourceFields.SourceOfTruth, request.SourceOfTruth)
	sourceFields.AppSource = utils.StringFirstNonEmpty(sourceFields.AppSource, request.AppSource)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	command := command.NewUpsertOrganizationCommand(organizationId, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.LoggedInUserId),
		sourceFields, externalSystem, dataFields, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt), request.IgnoreEmptyFields)
	if err := s.organizationCommands.UpsertOrganization.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertSyncOrganization.Handle) tenant:%s, organizationID: %s , err: {%v}", request.Tenant, organizationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Upserted organization %s", organizationId)

	return &organizationpb.OrganizationIdGrpcResponse{Id: organizationId}, nil
}

func (s *organizationService) LinkPhoneNumberToOrganization(ctx context.Context, request *organizationpb.LinkPhoneNumberToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkPhoneNumberToOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewLinkPhoneNumberCommand(request.OrganizationId, request.Tenant, request.LoggedInUserId, request.PhoneNumberId, request.Label, request.Primary)
	if err := s.organizationCommands.LinkPhoneNumberCommand.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkPhoneNumberToOrganization.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked phone number {%s} to organization {%s}", request.PhoneNumberId, request.OrganizationId)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) LinkEmailToOrganization(ctx context.Context, request *organizationpb.LinkEmailToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkEmailToOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewLinkEmailCommand(request.OrganizationId, request.Tenant, request.LoggedInUserId, request.EmailId, request.Label, request.Primary)
	if err := s.organizationCommands.LinkEmailCommand.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkEmailCommand.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked email {%s} to organization {%s}", request.EmailId, request.OrganizationId)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) LinkLocationToOrganization(ctx context.Context, request *organizationpb.LinkLocationToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkLocationToOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	command := command.NewLinkLocationCommand(request.OrganizationId, request.Tenant, request.LoggedInUserId, request.LocationId)
	if err := s.organizationCommands.LinkLocationCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkLocationCommand.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("Linked location {%s} to organization {%s}", request.LocationId, request.OrganizationId)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) LinkDomainToOrganization(ctx context.Context, request *organizationpb.LinkDomainToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkDomainToOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.LoggedInUserId))
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewLinkDomainCommand(request.OrganizationId, request.Tenant, request.Domain, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId), request.AppSource)
	if err := s.organizationCommands.LinkDomainCommand.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) UpdateOrganizationRenewalLikelihood(ctx context.Context, request *organizationpb.OrganizationRenewalLikelihoodRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateOrganizationRenewalLikelihood")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.LoggedInUserId))
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	fields := models.RenewalLikelihoodFields{
		RenewalLikelihood: mapper.MapRenewalLikelihoodToModels(request.Likelihood),
		Comment:           request.Comment,
		UpdatedBy:         utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId),
	}
	command := command.NewUpdateRenewalLikelihoodCommand(request.Tenant, request.OrganizationId, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId), fields)
	if err := s.organizationCommands.UpdateRenewalLikelihoodCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed update renewal likelihood for tenant: %s organizationID: %s, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Updated renewal likelihood for tenant:%s organizationID: %s", request.Tenant, request.OrganizationId)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) UpdateOrganizationRenewalForecast(ctx context.Context, request *organizationpb.OrganizationRenewalForecastRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateOrganizationRenewalForecast")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId))
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	fields := models.RenewalForecastFields{
		Amount:    request.Amount,
		Comment:   request.Comment,
		UpdatedBy: utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId),
	}
	command := command.NewUpdateRenewalForecastCommand(request.Tenant, request.OrganizationId, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId), fields, "")
	if err := s.organizationCommands.UpdateRenewalForecastCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed update renewal forecast for tenant: %s organizationID: %s, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Updated renewal forecast for tenant:%s organizationID: %s", request.Tenant, request.OrganizationId)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) UpdateOrganizationBillingDetails(ctx context.Context, request *organizationpb.OrganizationBillingDetailsRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateOrganizationBillingDetails")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId))
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	fields := models.BillingDetailsFields{
		Amount:            request.Amount,
		UpdatedBy:         utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId),
		Frequency:         mapper.MapFrequencyToString(request.Frequency),
		RenewalCycle:      mapper.MapFrequencyToString(request.RenewalCycle),
		RenewalCycleStart: utils.TimestampProtoToTime(request.CycleStart),
	}
	command := command.NewUpdateBillingDetailsCommand(request.Tenant, request.OrganizationId, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId), fields)
	if err := s.organizationCommands.UpdateBillingDetailsCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed update billing details for tenant: %s organizationID: %s, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Updated billing details for tenant:%s organizationID: %s", request.Tenant, request.OrganizationId)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) RequestRenewNextCycleDate(ctx context.Context, request *organizationpb.RequestRenewNextCycleDateRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RequestRenewNextCycleDate")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := command.NewRequestNextCycleDateCommand(request.Tenant, request.OrganizationId, request.LoggedInUserId)
	if err := s.organizationCommands.RequestNextCycleDateCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed request next cycle date for tenant: %s organizationID: %s, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Requested next cycle date renewal for tenant:%s organizationID: %s", request.Tenant, request.OrganizationId)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) HideOrganization(ctx context.Context, request *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.HideOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId))
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := command.NewHideOrganizationCommand(request.Tenant, request.OrganizationId, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId))
	if err := s.organizationCommands.HideOrganizationCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed hide organization with id  %s for tenant %s, err: %s", request.OrganizationId, request.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Hidden organization with id %s for tenant %s", request.OrganizationId, request.Tenant)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) ShowOrganization(ctx context.Context, request *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.ShowOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId))
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := command.NewShowOrganizationCommand(request.Tenant, request.OrganizationId, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId))
	if err := s.organizationCommands.ShowOrganizationCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed show organization with id  %s for tenant %s, err: %s", request.OrganizationId, request.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Show organization with id %s for tenant %s", request.OrganizationId, request.Tenant)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) UpsertCustomFieldToOrganization(ctx context.Context, request *organizationpb.CustomFieldForOrganizationGrpcRequest) (*organizationpb.CustomFieldIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpsertCustomFieldToOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.LoggedInUserId))
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	customFieldId := request.CustomFieldId
	if customFieldId == "" {
		customFieldId = uuid.New().String()
	}
	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	customField := models.CustomField{
		Id:         customFieldId,
		Name:       request.CustomFieldName,
		TemplateId: request.CustomFieldTemplateId,
		CustomFieldValue: models.CustomFieldValue{
			Str:     request.CustomFieldValue.StringValue,
			Bool:    request.CustomFieldValue.BoolValue,
			Time:    utils.TimestampProtoToTime(request.CustomFieldValue.DatetimeValue),
			Int:     request.CustomFieldValue.IntegerValue,
			Decimal: request.CustomFieldValue.DecimalValue,
		},
		CustomFieldDataType: mapper.MapCustomFieldDataType(request.CustomFieldDataType),
	}

	command := command.NewUpsertCustomFieldCommand(request.OrganizationId, request.Tenant,
		sourceFields.Source, sourceFields.SourceOfTruth, sourceFields.AppSource, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId),
		utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt), customField)
	if err := s.organizationCommands.UpsertCustomFieldCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	return &organizationpb.CustomFieldIdGrpcResponse{Id: customFieldId}, nil
}

func (s *organizationService) AddParentOrganization(ctx context.Context, request *organizationpb.AddParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.AddParentOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewAddParentCommand(request.OrganizationId, request.Tenant, request.LoggedInUserId, request.ParentOrganizationId, request.Type, request.AppSource)
	if err := s.organizationCommands.AddParentCommand.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddParentCommand.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) RemoveParentOrganization(ctx context.Context, request *organizationpb.RemoveParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RemoveParentOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.String("request", fmt.Sprintf("%+v", request)))

	cmd := command.NewRemoveParentCommand(request.OrganizationId, request.Tenant, request.LoggedInUserId, request.ParentOrganizationId, request.AppSource)
	if err := s.organizationCommands.RemoveParentCommand.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemoveParentCommand.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
