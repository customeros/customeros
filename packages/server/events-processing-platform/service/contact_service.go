package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"strings"
	"time"
)

type contactService struct {
	contactpb.UnimplementedContactGrpcServiceServer
	log      logger.Logger
	services *Services
}

func NewContactService(log logger.Logger, services *Services) *contactService {
	return &contactService{
		log:      log,
		services: services,
	}
}

func (s *contactService) LinkPhoneNumberToContact(ctx context.Context, request *contactpb.LinkPhoneNumberToContactGrpcRequest) (*contactpb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.LinkPhoneNumberToContact")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	agg, err := contact.LoadContactAggregate(ctx, s.services.es, request.Tenant, request.ContactId, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		agg = contact.NewContactAggregateWithTenantAndID(request.Tenant, request.ContactId)
	}

	if eventstore.AllowCheckForNoChanges(request.AppSource, request.LoggedInUserId) {
		if agg.Contact.HasPhoneNumber(request.PhoneNumberId, request.Label, request.Primary) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return &contactpb.ContactIdGrpcResponse{Id: request.ContactId}, nil
		}
	}

	evt, err := event.NewContactLinkPhoneNumberEvent(agg, request.PhoneNumberId, request.Label, request.Primary, time.Now())

	eventstore.EnrichEventWithMetadataExtended(&evt, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	err = agg.Apply(evt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	err = s.services.es.Save(ctx, agg)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	return &contactpb.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) LinkLocationToContact(ctx context.Context, request *contactpb.LinkLocationToContactGrpcRequest) (*contactpb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.LinkLocationToContact")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	agg, err := contact.LoadContactAggregate(ctx, s.services.es, request.Tenant, request.ContactId, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		agg = contact.NewContactAggregateWithTenantAndID(request.Tenant, request.ContactId)
	}

	if eventstore.AllowCheckForNoChanges(request.AppSource, request.LoggedInUserId) {
		if agg.Contact.HasLocation(request.LocationId) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return &contactpb.ContactIdGrpcResponse{Id: request.ContactId}, nil
		}
	}

	evt, err := event.NewContactLinkLocationEvent(agg, request.LocationId, time.Now())

	eventstore.EnrichEventWithMetadataExtended(&evt, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	err = agg.Apply(evt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	err = s.services.es.Save(ctx, agg)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	return &contactpb.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) LinkWithOrganization(ctx context.Context, request *contactpb.LinkWithOrganizationGrpcRequest) (*contactpb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.LinkWithOrganization")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	jobRoleFields := contact.JobRole{
		JobTitle:    request.JobTitle,
		Description: request.Description,
		Primary:     request.Primary,
		StartedAt:   utils.TimestampProtoToTimePtr(request.StartedAt),
		EndedAt:     utils.TimestampProtoToTimePtr(request.EndedAt),
	}

	createdAtNotNil := utils.IfNotNilTimeWithDefault(request.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(request.UpdatedAt, utils.Now())

	agg, err := contact.LoadContactAggregate(ctx, s.services.es, request.Tenant, request.ContactId, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		agg = contact.NewContactAggregateWithTenantAndID(request.Tenant, request.ContactId)
	}

	if eventstore.AllowCheckForNoChanges(request.AppSource, request.LoggedInUserId) {
		if agg.Contact.HasJobRoleInOrganization(request.OrganizationId, jobRoleFields, sourceFields) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return &contactpb.ContactIdGrpcResponse{Id: request.ContactId}, nil
		}
	}

	evt, err := event.NewContactLinkWithOrganizationEvent(agg, request.OrganizationId, request.JobTitle, request.Description,
		request.Primary, sourceFields, createdAtNotNil, updatedAtNotNil, utils.TimestampProtoToTimePtr(request.StartedAt), utils.TimestampProtoToTimePtr(request.EndedAt))

	eventstore.EnrichEventWithMetadataExtended(&evt, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	err = agg.Apply(evt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	err = s.services.es.Save(ctx, agg)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	return &contactpb.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) RemoveSocial(ctx context.Context, request *contactpb.ContactRemoveSocialGrpcRequest) (*contactpb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.RemoveSocial")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.ContactId)

	initAggregateFunc := func() eventstore.Aggregate {
		return contact.NewContactAggregateWithTenantAndID(request.Tenant, request.ContactId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemoveSocial.HandleGRPCRequest) tenant:{%s}, contact ID: {%s}, err: %s", request.Tenant, request.ContactId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &contactpb.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}

func normalizeTimezone(timezone string) string {
	if timezone == "" {
		return ""
	}
	output := strings.Replace(timezone, "_slash_", "/", -1)
	output = utils.CapitalizeAllParts(output, []string{"/", "_"})
	return output
}

func (s *contactService) EnrichContact(ctx context.Context, request *contactpb.EnrichContactGrpcRequest) (*contactpb.ContactIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.EnrichContact")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	initAggregateFunc := func() eventstore.Aggregate {
		return contact.NewContactTempAggregateWithTenantAndID(request.Tenant, request.ContactId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{SkipLoadEvents: true}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(EnrichContact.HandleGRPCRequest) tenant:{%s}, contact ID: {%s}, err: %s", request.Tenant, request.ContactId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &contactpb.ContactIdGrpcResponse{Id: request.ContactId}, nil
}

func (s *contactService) AddLocation(ctx context.Context, request *contactpb.ContactAddLocationGrpcRequest) (*locationpb.LocationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContactService.AddLocation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.ContactId)

	initAggregateFunc := func() eventstore.Aggregate {
		return contact.NewContactAggregateWithTenantAndID(request.Tenant, request.ContactId)
	}
	locationId, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(AddLocation.HandleGRPCRequest) tenant:{%s}, contact ID: {%s}, err: %s", request.Tenant, request.ContactId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &locationpb.LocationIdGrpcResponse{Id: locationId.(string)}, nil
}
