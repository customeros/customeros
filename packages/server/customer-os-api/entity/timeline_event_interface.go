package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

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
	//model.TimelineEventTypePageView.String():           neo4jentity.NodeLabelPageView,
	//model.TimelineEventTypeInteractionSession.String(): neo4jentity.NodeLabelInteractionSession,
	model.TimelineEventTypeIssue.String(): neo4jentity.NodeLabelIssue,
	//model.TimelineEventTypeNote.String():               neo4jentity.NodeLabelNote,
	model.TimelineEventTypeInteractionEvent.String(): neo4jentity.NodeLabelInteractionEvent,
	model.TimelineEventTypeMeeting.String():          neo4jentity.NodeLabelMeeting,
	model.TimelineEventTypeAction.String():           neo4jentity.NodeLabelAction,
	model.TimelineEventTypeLogEntry.String():         neo4jentity.NodeLabelLogEntry,
}

type TimelineEvent interface {
	IsTimelineEvent()
	TimelineEventLabel() string
	SetDataloaderKey(key string)
	GetDataloaderKey() string
}

type TimelineEventEntities []TimelineEvent
