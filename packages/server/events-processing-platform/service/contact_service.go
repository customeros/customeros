package service

import (
	"context"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
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

func (s *contactService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
