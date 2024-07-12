package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
)

type LastTouchpointType string

var NodeLabelsByTimelineEventType = map[string]string{
	//model.TimelineEventTypePageView.String():           neo4jutil.NodeLabelPageView,
	//model.TimelineEventTypeInteractionSession.String(): neo4jutil.NodeLabelInteractionSession,
	//model.TimelineEventTypeNote.String():               neo4jutil.NodeLabelNote,
	model.TimelineEventTypeIssue.String():            model2.NodeLabelIssue,
	model.TimelineEventTypeInteractionEvent.String(): model2.NodeLabelInteractionEvent,
	model.TimelineEventTypeMeeting.String():          model2.NodeLabelMeeting,
	model.TimelineEventTypeAction.String():           model2.NodeLabelAction,
	model.TimelineEventTypeLogEntry.String():         model2.NodeLabelLogEntry,
}

type TimelineEvent interface {
	IsTimelineEvent()
	TimelineEventLabel() string
	SetDataloaderKey(key string)
	GetDataloaderKey() string
}

type TimelineEventEntities []TimelineEvent
