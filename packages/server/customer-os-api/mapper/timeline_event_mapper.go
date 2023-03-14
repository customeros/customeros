package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"reflect"
)

func MapEntityToTimelineEvent(timelineEventEntity *entity.TimelineEvent) any {
	switch (*timelineEventEntity).TimelineEventName() {
	case entity.NodeLabel_PageView:
		pageViewEntity := (*timelineEventEntity).(*entity.PageViewEntity)
		return MapEntityToPageView(pageViewEntity)
	case entity.NodeLabel_InteractionSession:
		interactionSessionEntity := (*timelineEventEntity).(*entity.InteractionSessionEntity)
		return MapEntityToInteractionSession(interactionSessionEntity)
	case entity.NodeLabel_Ticket:
		ticketEntity := (*timelineEventEntity).(*entity.TicketEntity)
		return MapEntityToTicket(ticketEntity)
	case entity.NodeLabel_Conversation:
		conversationEntity := (*timelineEventEntity).(*entity.ConversationEntity)
		return MapEntityToConversation(conversationEntity)
	case entity.NodeLabel_Note:
		noteEntity := (*timelineEventEntity).(*entity.NoteEntity)
		return MapEntityToNote(noteEntity)
	}
	fmt.Errorf("timeline event of type %s not identified", reflect.TypeOf(timelineEventEntity))
	return nil
}

func MapEntitiesToTimelineEvents(entities *entity.TimelineEventEntities) []model.TimelineEvent {
	var timelineEvents []model.TimelineEvent
	for _, timelineEventEntity := range *entities {
		timelineEvent := MapEntityToTimelineEvent(&timelineEventEntity)
		if timelineEvent != nil {
			timelineEvents = append(timelineEvents, timelineEvent.(model.TimelineEvent))
		}
	}
	return timelineEvents
}
