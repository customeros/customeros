package graph

import (
	"context"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type MasterPlanEventHandler struct {
	log      logger.Logger
	services *service.Services
}

func NewMasterPlanEventHandler(log logger.Logger, services *service.Services) *MasterPlanEventHandler {
	return &MasterPlanEventHandler{
		log:      log,
		services: services,
	}
}

func (h *MasterPlanEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.MasterPlanCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	masterPlanId := aggregate.GetMasterPlanObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)

	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)
	err := h.services.CommonServices.Neo4jRepositories.MasterPlanWriteRepository.Create(ctx, eventData.Tenant, masterPlanId, eventData.Name, source, appSource, eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving master plan %s: %s", masterPlanId, err.Error())
		return err
	}
	return err
}

func (h *MasterPlanEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.MasterPlanUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	masterPlanId := aggregate.GetMasterPlanObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)

	data := neo4jrepository.MasterPlanUpdateFields{
		Name:          eventData.Name,
		Retired:       eventData.Retired,
		UpdateName:    eventData.UpdateName(),
		UpdateRetired: eventData.UpdateRetired(),
	}
	err := h.services.CommonServices.Neo4jRepositories.MasterPlanWriteRepository.Update(ctx, eventData.Tenant, masterPlanId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating master plan %s: %s", masterPlanId, err.Error())
		return err
	}
	return err
}

func (h *MasterPlanEventHandler) OnCreateMilestone(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanEventHandler.OnCreateMilestone")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.MasterPlanMilestoneCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	masterPlanId := aggregate.GetMasterPlanObjectID(evt.GetAggregateID(), eventData.Tenant)
	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)
	err := h.services.CommonServices.Neo4jRepositories.MasterPlanWriteRepository.CreateMilestone(ctx, eventData.Tenant, masterPlanId, eventData.MilestoneId,
		eventData.Name, source, appSource, eventData.Order, eventData.DurationHours, eventData.Items, eventData.Optional, eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving master plan milestone %s: %s", masterPlanId, err.Error())
		return err
	}
	return err
}

func (h *MasterPlanEventHandler) OnUpdateMilestone(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanEventHandler.OnUpdateMilestone")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.MasterPlanMilestoneUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	masterPlanId := aggregate.GetMasterPlanObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, eventData.MilestoneId)

	data := neo4jrepository.MasterPlanMilestoneUpdateFields{
		Name:                eventData.Name,
		Order:               eventData.Order,
		DurationHours:       eventData.DurationHours,
		Items:               eventData.Items,
		Optional:            eventData.Optional,
		Retired:             eventData.Retired,
		UpdateName:          eventData.UpdateName(),
		UpdateOrder:         eventData.UpdateOrder(),
		UpdateDurationHours: eventData.UpdateDurationHours(),
		UpdateItems:         eventData.UpdateItems(),
		UpdateOptional:      eventData.UpdateOptional(),
		UpdateRetired:       eventData.UpdateRetired(),
	}
	err := h.services.CommonServices.Neo4jRepositories.MasterPlanWriteRepository.UpdateMilestone(ctx, eventData.Tenant, masterPlanId, eventData.MilestoneId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating master plan milestone %s: %s", eventData.MilestoneId, err.Error())
		return err
	}
	return err
}

func (h *MasterPlanEventHandler) OnReorderMilestones(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanEventHandler.OnReorderMilestones")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.MasterPlanMilestoneReorderEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	masterPlanId := aggregate.GetMasterPlanObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)

	for i, milestoneId := range eventData.MilestoneIds {
		data := neo4jrepository.MasterPlanMilestoneUpdateFields{
			Order:       int64(i),
			UpdateOrder: true,
		}
		err := h.services.CommonServices.Neo4jRepositories.MasterPlanWriteRepository.UpdateMilestone(ctx, eventData.Tenant, masterPlanId, milestoneId, data)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while updating master plan milestone order %s: %s", milestoneId, err.Error())
			return err
		}
	}
	return nil
}
