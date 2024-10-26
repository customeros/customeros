package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ContactEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewContactEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *ContactEventHandler {
	return &ContactEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

func (h *ContactEventHandler) OnPhoneNumberLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnPhoneNumberLinkToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.PhoneNumberId, eventData.Label, eventData.Primary)

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTACT.String(), contactId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return err
}

func (h *ContactEventHandler) OnLocationLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnLocationLinkToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.LocationId)

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTACT.String(), contactId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return err
}

func (h *ContactEventHandler) OnContactLinkToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactLinkToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactLinkWithOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.JobRoleFields{
		Description: eventData.Description,
		JobTitle:    eventData.JobTitle,
		Primary:     eventData.Primary,
		CreatedAt:   eventData.CreatedAt,
		StartedAt:   eventData.StartedAt,
		EndedAt:     eventData.EndedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:        helper.GetSource(eventData.SourceFields.Source),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceFields.SourceOfTruth),
			AppSource:     helper.GetAppSource(eventData.SourceFields.AppSource),
		},
	}
	err := h.services.CommonServices.Neo4jRepositories.JobRoleWriteRepository.LinkContactWithOrganization(ctx, eventData.Tenant, contactId, eventData.OrganizationId, data)

	// Request last touch point update
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: eventData.OrganizationId,
			AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while refreshing last touchpoint for organization %s: %s", eventData.OrganizationId, err.Error())
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTACT.String(), contactId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
	utils.EventCompleted(ctx, eventData.Tenant, model.ORGANIZATION.String(), eventData.OrganizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return err
}

func (h *ContactEventHandler) OnSocialAddedToContactV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnSocialAddedToContactV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTACT.String(), contactId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *ContactEventHandler) OnSocialRemovedFromContactV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnSocialRemovedFromContactV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactRemoveSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	if eventData.SocialId != "" {
		err := h.services.CommonServices.Neo4jRepositories.SocialWriteRepository.RemoveSocialForEntityById(ctx, eventData.Tenant, contactId, model.NodeLabelContact, eventData.SocialId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil
		}
	} else {
		err := h.services.CommonServices.Neo4jRepositories.SocialWriteRepository.RemoveSocialForEntityByUrl(ctx, eventData.Tenant, contactId, model.NodeLabelContact, eventData.Url)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil
		}
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTACT.String(), contactId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *ContactEventHandler) OnLocationAddedToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnLocationAddedToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactAddLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	data := neo4jrepository.LocationCreateFields{
		RawAddress: eventData.RawAddress,
		Name:       eventData.Name,
		CreatedAt:  eventData.CreatedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:        helper.GetSource(eventData.Source),
			SourceOfTruth: helper.GetSource(eventData.SourceOfTruth),
			AppSource:     helper.GetSource(eventData.AppSource),
		},
		AddressDetails: neo4jrepository.AddressDetails{
			Latitude:      eventData.Latitude,
			Longitude:     eventData.Longitude,
			Country:       eventData.Country,
			CountryCodeA2: eventData.CountryCodeA2,
			CountryCodeA3: eventData.CountryCodeA3,
			Region:        eventData.Region,
			District:      eventData.District,
			Locality:      eventData.Locality,
			Street:        eventData.Street,
			Address:       eventData.AddressLine1,
			Address2:      eventData.AddressLine2,
			Zip:           eventData.ZipCode,
			AddressType:   eventData.AddressType,
			HouseNumber:   eventData.HouseNumber,
			PostalCode:    eventData.PostalCode,
			PlusFour:      eventData.PlusFour,
			Commercial:    eventData.Commercial,
			Predirection:  eventData.Predirection,
			TimeZone:      eventData.TimeZone,
			UtcOffset:     eventData.UtcOffset,
		},
	}

	err := h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.CreateLocation(ctx, eventData.Tenant, eventData.LocationId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while creating location %s: %s", eventData.LocationId, err.Error())
		return err
	}
	err = h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.LocationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while linking location %s to contact %s: %s", eventData.LocationId, contactId, err.Error())
		return err
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTACT.String(), contactId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}
