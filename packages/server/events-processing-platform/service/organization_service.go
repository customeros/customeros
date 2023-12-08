package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
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
	tracing.LogObjectAsJson(span, "request", request)

	organizationId := utils.NewUUIDIfEmpty(request.Id)

	dataFields := model.OrganizationDataFields{
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

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	upsertCommand := command.NewUpsertOrganizationCommand(organizationId, request.Tenant, request.LoggedInUserId, sourceFields, externalSystem, dataFields,
		utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt),
		extractOrganizationMaskFields(request.FieldsMask))
	if err := s.organizationCommands.UpsertOrganization.Handle(ctx, upsertCommand); err != nil {
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
	tracing.LogObjectAsJson(span, "request", request)

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
	tracing.LogObjectAsJson(span, "request", request)

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
	tracing.LogObjectAsJson(span, "request", request)

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
	tracing.LogObjectAsJson(span, "request", request)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	cmd := command.NewLinkDomainCommand(request.OrganizationId, request.Tenant, request.Domain, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId), request.AppSource)
	if err := s.organizationCommands.LinkDomainCommand.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) HideOrganization(ctx context.Context, request *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.HideOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := command.NewHideOrganizationCommand(request.Tenant, request.OrganizationId, request.LoggedInUserId)
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
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := command.NewShowOrganizationCommand(request.Tenant, request.OrganizationId, request.LoggedInUserId)
	if err := s.organizationCommands.ShowOrganizationCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed show organization with id  %s for tenant %s, err: %s", request.OrganizationId, request.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Show organization with id %s for tenant %s", request.OrganizationId, request.Tenant)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) RefreshLastTouchpoint(ctx context.Context, request *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RefreshLastTouchpoint")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	command := command.NewRefreshLastTouchpointCommand(request.Tenant, request.OrganizationId, request.LoggedInUserId, request.AppSource)
	if err := s.organizationCommands.RefreshLastTouchpointCommand.Handle(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to refresh the last touchpoint for organization with id  %s in tenant %s, err: %s", request.OrganizationId, request.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

	s.log.Infof("Refresh the last touchpoint for organization with id %s in tenant %s", request.OrganizationId, request.Tenant)

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) UpsertCustomFieldToOrganization(ctx context.Context, request *organizationpb.CustomFieldForOrganizationGrpcRequest) (*organizationpb.CustomFieldIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpsertCustomFieldToOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.LoggedInUserId))
	tracing.LogObjectAsJson(span, "request", request)

	customFieldId := request.CustomFieldId
	if customFieldId == "" {
		customFieldId = uuid.New().String()
	}
	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	customField := model.CustomField{
		Id:         customFieldId,
		Name:       request.CustomFieldName,
		TemplateId: request.CustomFieldTemplateId,
		CustomFieldValue: model.CustomFieldValue{
			Str:     request.CustomFieldValue.StringValue,
			Bool:    request.CustomFieldValue.BoolValue,
			Time:    utils.TimestampProtoToTimePtr(request.CustomFieldValue.DatetimeValue),
			Int:     request.CustomFieldValue.IntegerValue,
			Decimal: request.CustomFieldValue.DecimalValue,
		},
		CustomFieldDataType: mapper.MapCustomFieldDataType(request.CustomFieldDataType),
	}

	command := command.NewUpsertCustomFieldCommand(request.OrganizationId, request.Tenant,
		sourceFields.Source, sourceFields.SourceOfTruth, sourceFields.AppSource, utils.StringFirstNonEmpty(request.LoggedInUserId, request.UserId),
		utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt), customField)
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
	tracing.LogObjectAsJson(span, "request", request)

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
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewRemoveParentCommand(request.OrganizationId, request.Tenant, request.LoggedInUserId, request.ParentOrganizationId, request.AppSource)
	if err := s.organizationCommands.RemoveParentCommand.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemoveParentCommand.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) WebScrapeOrganization(ctx context.Context, request *organizationpb.WebScrapeOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.WebScrapeOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewWebScrapeOrganizationCommand(request.Tenant, request.OrganizationId, request.LoggedInUserId, request.AppSource, request.Url)
	if err := s.organizationCommands.WebScrapeOrganization.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(WebScrapeOrganization.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}

func extractOrganizationMaskFields(requestMaskFields []organizationpb.OrganizationMaskField) []string {
	fieldsMask := make([]string, 0)
	if requestMaskFields == nil || len(requestMaskFields) == 0 {
		return fieldsMask
	}
	if containsOrganizationMaskFieldAll(requestMaskFields) {
		return fieldsMask
	}
	for _, field := range requestMaskFields {
		switch field {
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NAME:
			fieldsMask = append(fieldsMask, model.FieldMaskName)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_TARGET_AUDIENCE:
			fieldsMask = append(fieldsMask, model.FieldMaskTargetAudience)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_VALUE_PROPOSITION:
			fieldsMask = append(fieldsMask, model.FieldMaskValueProposition)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_INDUSTRY:
			fieldsMask = append(fieldsMask, model.FieldMaskIndustry)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_SUB_INDUSTRY:
			fieldsMask = append(fieldsMask, model.FieldMaskSubIndustry)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_INDUSTRY_GROUP:
			fieldsMask = append(fieldsMask, model.FieldMaskIndustryGroup)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_WEBSITE:
			fieldsMask = append(fieldsMask, model.FieldMaskWebsite)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_EMPLOYEES:
			fieldsMask = append(fieldsMask, model.FieldMaskEmployees)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_MARKET:
			fieldsMask = append(fieldsMask, model.FieldMaskMarket)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_LAST_FUNDING_ROUND:
			fieldsMask = append(fieldsMask, model.FieldMaskLastFundingRound)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_LAST_FUNDING_AMOUNT:
			fieldsMask = append(fieldsMask, model.FieldMaskLastFundingAmount)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_REFERENCE_ID:
			fieldsMask = append(fieldsMask, model.FieldMaskReferenceId)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NOTE:
			fieldsMask = append(fieldsMask, model.FieldMaskNote)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_IS_PUBLIC:
			fieldsMask = append(fieldsMask, model.FieldMaskIsPublic)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_IS_CUSTOMER:
			fieldsMask = append(fieldsMask, model.FieldMaskIsCustomer)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_HIDE:
			fieldsMask = append(fieldsMask, model.FieldMaskHide)
		case organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_DESCRIPTION:
			fieldsMask = append(fieldsMask, model.FieldMaskDescription)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}

func containsOrganizationMaskFieldAll(fields []organizationpb.OrganizationMaskField) bool {
	for _, field := range fields {
		if field == organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_ALL {
			return true
		}
	}
	return false
}
