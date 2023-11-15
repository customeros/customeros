package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

type LastTouchpointType string

const (
	LastTouchpointTypePageView                  LastTouchpointType = "PAGE_VIEW"
	LastTouchpointTypeInteractionSession        LastTouchpointType = "INTERACTION_SESSION"
	LastTouchpointTypeNote                      LastTouchpointType = "NOTE"
	LastTouchpointTypeInteractionEventEmailSent LastTouchpointType = "INTERACTION_EVENT_EMAIL_SENT"
	LastTouchpointTypeInteractionEventPhoneCall LastTouchpointType = "INTERACTION_EVENT_PHONE_CALL"
	LastTouchpointTypeInteractionEventChat      LastTouchpointType = "INTERACTION_EVENT_CHAT"
	LastTouchpointTypeMeeting                   LastTouchpointType = "MEETING"
	LastTouchpointTypeAnalysis                  LastTouchpointType = "ANALYSIS"
	LastTouchpointTypeActionCreated             LastTouchpointType = "ACTION_CREATED"
	LastTouchpointTypeAction                    LastTouchpointType = "ACTION"
	LastTouchpointTypeLogEntry                  LastTouchpointType = "LOG_ENTRY"
	LastTouchpointTypeIssueCreated              LastTouchpointType = "ISSUE_CREATED"
	LastTouchpointTypeIssueUpdated              LastTouchpointType = "ISSUE_UPDATED"
)

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
