package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
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
	//model.TimelineEventTypePageView.String():           neo4jutil.NodeLabelPageView,
	//model.TimelineEventTypeInteractionSession.String(): neo4jutil.NodeLabelInteractionSession,
	model.TimelineEventTypeIssue.String(): neo4jutil.NodeLabelIssue,
	//model.TimelineEventTypeNote.String():               neo4jutil.NodeLabelNote,
	model.TimelineEventTypeInteractionEvent.String(): neo4jutil.NodeLabelInteractionEvent,
	model.TimelineEventTypeMeeting.String():          neo4jutil.NodeLabelMeeting,
	model.TimelineEventTypeAction.String():           neo4jutil.NodeLabelAction,
	model.TimelineEventTypeLogEntry.String():         neo4jutil.NodeLabelLogEntry,
}

type TimelineEvent interface {
	IsTimelineEvent()
	TimelineEventLabel() string
	SetDataloaderKey(key string)
	GetDataloaderKey() string
}

type TimelineEventEntities []TimelineEvent
