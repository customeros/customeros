package reminder

import (
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type ReminderEventHandler struct {
	log                     logger.Logger
	repositories            *repository.Repositories
	cfg                     config.Config
	eventBufferStoreService *eventbuffer.EventBufferStoreService
}

func NewReminderEventHandler(log logger.Logger, repositories *repository.Repositories, cfg config.Config, eventBufferStoreService *eventbuffer.EventBufferStoreService) *ReminderEventHandler {
	return &ReminderEventHandler{
		log:                     log,
		repositories:            repositories,
		cfg:                     cfg,
		eventBufferStoreService: eventBufferStoreService,
	}
}

func (h *ReminderEventHandler) onReminderCreateV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderEventHandler.onReminderCreateV1")
	defer span.Finish()

	var eventData event.ReminderCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	if eventData.DueDate.Before(utils.Now()) {
		return nil
	}

	notificationEvent := createReminderParkedEvent(eventData.Tenant, eventData.EntityId, eventData.Content, eventData.UserId, eventData.OrganizationId)

	err := h.eventBufferStoreService.ParkBaseEventWithId(ctx, &notificationEvent, eventData.Tenant, eventData.DueDate, "reminder-"+eventData.Tenant+"-"+eventData.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (h *ReminderEventHandler) onReminderUpdateV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderEventHandler.onReminderUpdateV1")
	defer span.Finish()

	var eventData event.ReminderUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	uuid := "reminder-" + eventData.Tenant + "-" + eventData.EntityId
	parkedReminder, err := h.eventBufferStoreService.GetById(uuid)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if eventData.DueDate.Before(utils.Now()) {
		//todo fix FE to stop sending due date in past
		return nil
	}

	if eventData.Dismissed {
		err := h.eventBufferStoreService.Delete(parkedReminder)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil
	}

	//if reminder is missing, we recreate it
	if parkedReminder == nil {
		return errors.New("reminder not found")
	} else {
		//overrite it
		var parkedReminderEventData event.ReminderNotificationEvent
		if err := json.Unmarshal(parkedReminder.EventData, &parkedReminderEventData); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		parkedReminderEventData.Content = eventData.Content

		parkedReminder.EventData, err = json.Marshal(parkedReminderEventData)
		parkedReminder.ExpiryTimestamp = eventData.DueDate.UTC()

		err := h.eventBufferStoreService.Update(parkedReminder)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func createReminderParkedEvent(tenant, reminderId, content, userId, organizationId string) event.ReminderNotificationEvent {
	return event.ReminderNotificationEvent{
		BaseEvent: event.BaseEvent{
			Tenant:     tenant,
			EventName:  event.ReminderNotificationV1,
			CreatedAt:  time.Now().UTC(),
			AppSource:  constants.AppSourceEventProcessingPlatformSubscribers,
			Source:     neo4jentity.DataSourceOpenline.String(),
			EntityId:   reminderId,
			EntityType: model.REMINDER,
		},
		Content:        content,
		UserId:         userId,
		OrganizationId: organizationId,
	}
}
