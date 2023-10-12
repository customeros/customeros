package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

var NodeLabelsByTimelineEventType = map[string]string{
	//model.TimelineEventTypePageView.String():           NodeLabel_PageView,
	//model.TimelineEventTypeInteractionSession.String(): NodeLabel_InteractionSession,
	model.TimelineEventTypeIssue.String(): NodeLabel_Issue,
	//model.TimelineEventTypeNote.String():               NodeLabel_Note,
	model.TimelineEventTypeInteractionEvent.String(): NodeLabel_InteractionEvent,
	model.TimelineEventTypeMeeting.String():          NodeLabel_Meeting,
	model.TimelineEventTypeAction.String():           NodeLabel_Action,
	model.TimelineEventTypeLogEntry.String():         NodeLabel_LogEntry,
}

type TimelineEvent interface {
	IsTimelineEvent()
	TimelineEventLabel() string
	SetDataloaderKey(key string)
	GetDataloaderKey() string
}

type TimelineEventEntities []TimelineEvent
