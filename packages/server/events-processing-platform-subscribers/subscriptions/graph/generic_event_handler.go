package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/generic"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type GenericEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewGenericEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *GenericEventHandler {
	return &GenericEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

func (h *GenericEventHandler) OnLinkEntityWithEntityV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GenericEventHandler.OnLinkEntityWithEntityV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData generic.LinkEntityWithEntity
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, eventData.EntityId)

	err := h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, eventData.Tenant, repository.LinkDetails{
		FromEntityId:   eventData.EntityId,
		FromEntityType: eventData.EntityType,
		Relationship:   eventData.Relationship,
		ToEntityId:     eventData.WithEntityId,
		ToEntityType:   eventData.WithEntityType,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
