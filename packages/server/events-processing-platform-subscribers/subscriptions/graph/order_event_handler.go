package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/order"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OrderEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewOrderEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *OrderEventHandler {
	return &OrderEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

func (h *OrderEventHandler) OnUpsertOrderV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderEventHandler.OnUpsertOrderV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData order.OrderUpsertEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	orderId := order.GetOrderObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, orderId)

	err := h.services.CommonServices.Neo4jRepositories.OrderWriteRepository.UpsertOrder(ctx, eventData.Tenant, eventData.OrganizationId, orderId, eventData.CreatedAt, eventData.ConfirmedAt, eventData.PaidAt, eventData.FulfilledAt, eventData.CanceledAt, neo4jmodel.Source{
		Source:    eventData.SourceFields.Source,
		AppSource: eventData.SourceFields.AppSource,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving order %s: %s", orderId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		externalSystemData := neo4jmodel.ExternalSystem{
			ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
			ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
			ExternalId:       eventData.ExternalSystem.ExternalId,
			ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
			ExternalSource:   eventData.ExternalSystem.ExternalSource,
			SyncDate:         eventData.ExternalSystem.SyncDate,
		}
		err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, orderId, model.NodeLabelOrder, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link order %s with external system %s: %s", orderId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return nil
}
