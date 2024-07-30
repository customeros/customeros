package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/generic"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type GenericEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewGenericEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *GenericEventHandler {
	return &GenericEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
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

	err := h.repositories.Neo4jRepositories.CommonWriteRepository.LinkEntityWithEntity(ctx, eventData.Tenant, eventData.EntityId, eventData.EntityType, eventData.RelationshipName, eventData.WithEntityId, eventData.WithEntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
