package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

const (
	NodeLabel_PageView           = "PageView"
	NodeLabel_InteractionSession = "InteractionSession"
	NodeLabel_Ticket             = "Ticket"
	NodeLabel_Conversation       = "Conversation"
	NodeLabel_Note               = "Note"
	NodeLabel_InteractionEvent   = "InteractionEvent"
)

var NodeLabelsByTimelineEventType = map[string]string{
	model.TimelineEventTypePageView.String():           NodeLabel_PageView,
	model.TimelineEventTypeInteractionSession.String(): NodeLabel_InteractionSession,
	model.TimelineEventTypeTicket.String():             NodeLabel_Ticket,
	model.TimelineEventTypeConversation.String():       NodeLabel_Conversation,
	model.TimelineEventTypeNote.String():               NodeLabel_Note,
	model.TimelineEventTypeInteractionEvent.String():   NodeLabel_InteractionEvent,
}

type TimelineEvent interface {
	TimelineEvent()
	TimelineEventName() string
}

type TimelineEventEntities []TimelineEvent
