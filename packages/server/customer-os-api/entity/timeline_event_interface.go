package entity

import (
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

type LastTouchpointType string

var NodeLabelsByTimelineEventType = map[string]string{
	// model.TimelineEventTypePageView.String():           commonmodel.NodeLabelPageView,
	// model.TimelineEventTypeInteractionSession.String(): commonmodel.NodeLabelInteractionSession,
	// model.TimelineEventTypeNote.String():               commonmodel.NodeLabelNote,
	model.TimelineEventTypeIssue.String():            commonmodel.NodeLabelIssue,
	model.TimelineEventTypeInteractionEvent.String(): commonmodel.NodeLabelInteractionEvent,
	model.TimelineEventTypeMeeting.String():          commonmodel.NodeLabelMeeting,
	model.TimelineEventTypeAction.String():           commonmodel.NodeLabelAction,
	model.TimelineEventTypeLogEntry.String():         commonmodel.NodeLabelLogEntry,
	model.TimelineEventTypeMarkdownEvent.String():    commonmodel.NodeLabelMarkdownEvent,
}

type TimelineEvent interface {
	IsTimelineEvent()
	TimelineEventLabel() string
	SetDataloaderKey(key string)
	GetDataloaderKey() string
}

type TimelineEventEntities []TimelineEvent
