package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
)

type LastTouchpointType string

var NodeLabelsByTimelineEventType = map[string]string{
	//model.TimelineEventTypePageView.String():           neo4jutil.NodeLabelPageView,
	//model.TimelineEventTypeInteractionSession.String(): neo4jutil.NodeLabelInteractionSession,
	//model.TimelineEventTypeNote.String():               neo4jutil.NodeLabelNote,
	model.TimelineEventTypeIssue.String():            neo4jutil.NodeLabelIssue,
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
