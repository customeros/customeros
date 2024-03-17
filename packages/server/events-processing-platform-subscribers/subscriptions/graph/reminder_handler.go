package graph

import (
	"context"

	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ReminderEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewReminderEventHandler(log logger.Logger, repositories *repository.Repositories) *ReminderEventHandler {
	return &ReminderEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *ReminderEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.ReminderCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	reminderId := aggregate.GetReminderObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, reminderId)

	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)

	err := h.repositories.Neo4jRepositories.ReminderWriteRepository.CreateReminder(ctx, eventData.Tenant, reminderId, eventData.UserId, eventData.OrganizationId, eventData.Content, source, appSource, eventData.CreatedAt, eventData.DueDate)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving reminder %s: %s", reminderId, err.Error())
		return err
	}
	err = h.repositories.Neo4jRepositories.ReminderWriteRepository.LinkReminderToUser(ctx, eventData.Tenant, reminderId, eventData.UserId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while linking reminder %s to user %s: %s", reminderId, eventData.UserId, err.Error())
		return err
	}

	err = h.repositories.Neo4jRepositories.ReminderWriteRepository.LinkReminderToOrganization(ctx, eventData.Tenant, reminderId, eventData.OrganizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while linking reminder %s to organization %s: %s", reminderId, eventData.OrganizationId, err.Error())
		return err
	}

	return err
}

func (h *ReminderEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.ReminderUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	reminderId := aggregate.GetReminderObjectID(evt.AggregateID, eventData.Tenant)

	due := eventData.DueDate
	span.SetTag(tracing.SpanTagEntityId, reminderId)

	updateData := neo4jrepo.ReminderUpdateFields{
		Content:         &eventData.Content,
		DueDate:         &due,
		Dismissed:       &eventData.Dismissed,
		UpdateContent:   eventData.UpdateContent(),
		UpdateDueDate:   eventData.UpdateDueDate(),
		UpdateDismissed: eventData.UpdateDismissed(),
	}

	err := h.repositories.Neo4jRepositories.ReminderWriteRepository.UpdateReminder(ctx, eventData.Tenant, reminderId, updateData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating reminder plan %s: %s", reminderId, err.Error())
		return err
	}

	return err
}
