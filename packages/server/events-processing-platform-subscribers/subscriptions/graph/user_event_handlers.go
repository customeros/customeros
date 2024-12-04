package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type UserEventHandler struct {
	log      logger.Logger
	services *service.Services
}

func NewUserEventHandler(log logger.Logger, services *service.Services) *UserEventHandler {
	return &UserEventHandler{
		log:      log,
		services: services,
	}
}

func (h *UserEventHandler) OnJobRoleLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnJobRoleLinkedToUser")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserLinkJobRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.JobRoleWriteRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.JobRoleId)

	return err
}

func (h *UserEventHandler) OnPhoneNumberLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnPhoneNumberLinkedToUser")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.PhoneNumberId, eventData.Label, eventData.Primary)

	return err
}

func (h *UserEventHandler) OnAddRole(c context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(c, "UserEventHandler.OnAddRole")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserAddRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: eventData.Tenant})

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.UserWriteRepository.AddRole(ctx, userId, eventData.Role)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while adding role %s to user %s: %s", eventData.Role, userId, err.Error())
	}

	return err
}

func (h *UserEventHandler) OnRemoveRole(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnRemoveRole")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserRemoveRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.UserWriteRepository.RemoveRole(ctx, eventData.Tenant, userId, eventData.Role)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while removing role %s from user %s: %s", eventData.Role, userId, err.Error())
	}

	return err
}
