package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type organizationService struct {
	organizationpb.UnimplementedOrganizationGrpcServiceServer
	services                   *Services
	log                        logger.Logger
	organizationCommands       *command_handler.CommandHandlers
	organizationRequestHandler organization.OrganizationRequestHandler
}

func NewOrganizationService(log logger.Logger, organizationCommands *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore, cfg *config.Config, services *Services) *organizationService {
	return &organizationService{
		log:                        log,
		services:                   services,
		organizationCommands:       organizationCommands,
		organizationRequestHandler: organization.NewOrganizationRequestHandler(log, aggregateStore, cfg.Utils),
	}
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

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) LinkLocationToOrganization(ctx context.Context, request *organizationpb.LinkLocationToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkLocationToOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewLinkLocationCommand(request.OrganizationId, request.Tenant, request.LoggedInUserId, request.LocationId)
	if err := s.organizationCommands.LinkLocationCommand.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkLocationCommand.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) UnlinkDomainFromOrganization(ctx context.Context, request *organizationpb.UnLinkDomainFromOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UnlinkDomainFromOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, utils.StringFirstNonEmpty(request.LoggedInUserId, request.LoggedInUserId))
	tracing.LogObjectAsJson(span, "request", request)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	// Check if the contract ID is valid
	if request.OrganizationId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOrganizationAggregateWithTenantAndID(request.Tenant, request.OrganizationId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{SkipLoadEvents: true}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UnlinkDomainFromOrganization.Handle) tenant:{%s}, err: %s", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) RefreshRenewalSummary(ctx context.Context, request *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RefreshRenewalSummary")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	_, err := s.organizationRequestHandler.HandleTempWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to refresh renewal summary for organization with id {%s} for tenant {%s}, err: %s", request.OrganizationId, request.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) RefreshDerivedData(ctx context.Context, request *organizationpb.RefreshDerivedDataGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RefreshDerivedData")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOrganizationTempAggregateWithTenantAndID(request.Tenant, request.OrganizationId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{SkipLoadEvents: true}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RefreshDerivedData.Handle) tenant:{%s}, err: %s", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) RefreshArr(ctx context.Context, request *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RefreshArr")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	cmd := command.NewRefreshArrCommand(request.Tenant, request.OrganizationId, request.LoggedInUserId, request.AppSource)
	if err := s.organizationCommands.RefreshArr.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to refresh ARR for organization with id  %s in tenant %s, err: %s", request.OrganizationId, request.Tenant, err.Error())
		return nil, s.errResponse(err)
	}

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
		CustomFieldValue: neo4jmodel.CustomFieldValue{
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

func (s *organizationService) EnrichOrganization(ctx context.Context, request *organizationpb.EnrichOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.EnrichOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOrganizationTempAggregateWithTenantAndID(request.Tenant, request.OrganizationId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{SkipLoadEvents: true}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(EnrichOrganization.HandleGRPCRequest) tenant:{%s}, organization ID: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) UpdateOnboardingStatus(ctx context.Context, request *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateOnboardingStatus")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	cmd := command.NewUpdateOnboardingStatusCommand(request.Tenant, request.OrganizationId, request.LoggedInUserId, request.AppSource,
		model.OnboardingStatus(request.OnboardingStatus).String(), request.Comments, request.CausedByContractId, utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err := s.organizationCommands.UpdateOnboardingStatus.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOnboardingStatus.Handle) tenant:{%s}, organization ID: {%s}, err: {%v}", request.Tenant, request.OrganizationId, err)
		return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) RemoveSocial(ctx context.Context, request *organizationpb.RemoveSocialGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.RemoveSocial")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOrganizationAggregateWithTenantAndID(request.Tenant, request.OrganizationId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemoveSocial.HandleGRPCRequest) tenant:{%s}, organization ID: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) AddLocation(ctx context.Context, request *organizationpb.OrganizationAddLocationGrpcRequest) (*locationpb.LocationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.AddLocation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOrganizationAggregateWithTenantAndID(request.Tenant, request.OrganizationId)
	}
	locationId, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddLocation.HandleGRPCRequest) tenant:{%s}, organization ID: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &locationpb.LocationIdGrpcResponse{Id: locationId.(string)}, nil
}

func (s *organizationService) UpdateOrganizationOwner(ctx context.Context, request *organizationpb.UpdateOrganizationOwnerGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateOrganizationOwner")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	cmd := command.NewUpdateOrganizationOwnerCommand(request.OrganizationId, request.Tenant, request.LoggedInUserId, request.OwnerUserId, request.AppSource)

	if err := s.organizationCommands.UpdateOrganizationOwner.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOrganizationOwner.Handle) tenant:{%s}, organization ID: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, s.errResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}

func (s *organizationService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}

func (s *organizationService) CreateBillingProfile(ctx context.Context, request *organizationpb.CreateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.CreateBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	result, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	billingProfileId := ""
	if result != nil {
		billingProfileId = result.(string)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: billingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: billingProfileId}, nil
}

func (s *organizationService) UpdateBillingProfile(ctx context.Context, request *organizationpb.UpdateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UpdateBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) LinkEmailToBillingProfile(ctx context.Context, request *organizationpb.LinkEmailToBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkEmailToBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkEmailToBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) UnlinkEmailFromBillingProfile(ctx context.Context, request *organizationpb.UnlinkEmailFromBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UnlinkEmailFromBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UnlinkEmailFromBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) LinkLocationToBillingProfile(ctx context.Context, request *organizationpb.LinkLocationToBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.LinkLocationToBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(LinkLocationToBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) UnlinkLocationFromBillingProfile(ctx context.Context, request *organizationpb.UnlinkLocationFromBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.UnlinkLocationFromBillingProfile")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.organizationRequestHandler.HandleWithRetry(ctx, request.Tenant, request.OrganizationId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UnlinkLocationFromBillingProfile) tenant:{%s}, organization id: {%s}, err: %s", request.Tenant, request.OrganizationId, err.Error())
		return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, s.errResponse(err)
	}

	return &organizationpb.BillingProfileIdGrpcResponse{Id: request.BillingProfileId}, nil
}

func (s *organizationService) AdjustIndustry(ctx context.Context, request *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrganizationService.AdjustIndustry")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OrganizationId)

	// handle deadlines
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "Context canceled")
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOrganizationTempAggregateWithTenantAndID(request.Tenant, request.OrganizationId)
	}
	params := map[string]interface{}{"rpc": "AdjustIndustry"}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{SkipLoadEvents: true}, request, params); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AdjustIndustry.Handle) tenant:{%s}, err: %s", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &organizationpb.OrganizationIdGrpcResponse{Id: request.OrganizationId}, nil
}
