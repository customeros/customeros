package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
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
