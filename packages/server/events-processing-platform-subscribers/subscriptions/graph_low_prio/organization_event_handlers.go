package graph_low_prio

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OrganizationEventHandler struct {
	services    *service.Services
	log         logger.Logger
	grpcClients *grpc_client.Clients
}

func NewOrganizationEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *OrganizationEventHandler {
	return &OrganizationEventHandler{
		services:    services,
		log:         log,
		grpcClients: grpcClients,
	}
}

func (h *OrganizationEventHandler) OnRefreshLastTouchPointV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshLastTouchPointV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshLastTouchpointEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    eventData.Tenant,
		AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
	})

	err := h.services.CommonServices.OrganizationService.RefreshLastTouchpoint(ctx, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to refresh last touchpoint for organization %s: %v", organizationId, err.Error())
	}

	return nil
}
