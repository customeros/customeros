package graph

import (
	"context"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type MasterPlanEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewMasterPlanEventHandler(log logger.Logger, repositories *repository.Repositories) *MasterPlanEventHandler {
	return &MasterPlanEventHandler{
		log:          log,
		repositories: repositories,
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
	err := h.repositories.Neo4jRepositories.MasterPlanWriteRepository.Create(ctx, eventData.Tenant, masterPlanId, eventData.Name, source, appSource, eventData.CreatedAt)
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
		UpdatedAt:     eventData.UpdatedAt,
		UpdateName:    eventData.UpdateName(),
		UpdateRetired: eventData.UpdateRetired(),
	}
	err := h.repositories.Neo4jRepositories.MasterPlanWriteRepository.Update(ctx, eventData.Tenant, masterPlanId, data)
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
	err := h.repositories.Neo4jRepositories.MasterPlanWriteRepository.CreateMilestone(ctx, eventData.Tenant, masterPlanId, eventData.MilestoneId,
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
		UpdatedAt:           eventData.UpdatedAt,
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
	err := h.repositories.Neo4jRepositories.MasterPlanWriteRepository.UpdateMilestone(ctx, eventData.Tenant, masterPlanId, eventData.MilestoneId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating master plan milestone %s: %s", eventData.MilestoneId, err.Error())
		return err
	}
	return err
}
