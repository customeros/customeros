package graph

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
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

	reminderId := eventData.Id
	span.SetTag(tracing.SpanTagEntityId, reminderId)

	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)

	err := h.repositories.Neo4jRepositories.ReminderWriteRepository.CreateReminder(ctx, eventData.Tenant, reminderId, eventData.UserId, eventData.OrganizationId, eventData.Content, source, appSource, eventData.CreatedAt, eventData.DueDate)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving reminder plan %s: %s", eventData.Id, err.Error())
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

	reminderId := eventData.ReminderId
	due := eventData.DueDate
	span.SetTag(tracing.SpanTagEntityId, reminderId)

	err := h.repositories.Neo4jRepositories.ReminderWriteRepository.UpdateReminder(ctx, eventData.Tenant, reminderId, &eventData.Content, &due, &eventData.Dismissed)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating reminder plan %s: %s", eventData.ReminderId, err.Error())
		return err
	}

	return err
}
