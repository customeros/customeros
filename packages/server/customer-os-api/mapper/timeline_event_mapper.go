package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
)

func MapEntityToTimelineEvent(timelineEventEntity *entity.TimelineEvent) model.TimelineEvent {
	if timelineEventEntity == nil || *timelineEventEntity == nil {
		return nil
	}
	switch (*timelineEventEntity).TimelineEventLabel() {
	case neo4jutil.NodeLabelPageView:
		pageViewEntity := (*timelineEventEntity).(*entity.PageViewEntity)
		return MapEntityToPageView(pageViewEntity)
	case neo4jutil.NodeLabelInteractionSession:
		interactionSessionEntity := (*timelineEventEntity).(*entity.InteractionSessionEntity)
		return MapEntityToInteractionSession(interactionSessionEntity)
	case neo4jutil.NodeLabelIssue:
		issueEntity := (*timelineEventEntity).(*entity.IssueEntity)
		return MapEntityToIssue(issueEntity)
	case neo4jutil.NodeLabelNote:
		noteEntity := (*timelineEventEntity).(*entity.NoteEntity)
		return MapEntityToNote(noteEntity)
	case neo4jutil.NodeLabelInteractionEvent:
		interactionEventEntity := (*timelineEventEntity).(*entity.InteractionEventEntity)
		return MapEntityToInteractionEvent(interactionEventEntity)
	case neo4jutil.NodeLabelAnalysis:
		analysisEntity := (*timelineEventEntity).(*entity.AnalysisEntity)
		return MapEntityToAnalysis(analysisEntity)
	case neo4jutil.NodeLabelMeeting:
		meetingEntity := (*timelineEventEntity).(*entity.MeetingEntity)
		return MapEntityToMeeting(meetingEntity)
	case neo4jutil.NodeLabelAction:
		actionEntity := (*timelineEventEntity).(*entity.ActionEntity)
		return MapEntityToAction(actionEntity)
	case neo4jutil.NodeLabelLogEntry:
		logEntryEntity := (*timelineEventEntity).(*neo4jentity.LogEntryEntity)
		return MapEntityToLogEntry(logEntryEntity)
	case neo4jutil.NodeLabelOrder:
		orderEntity := (*timelineEventEntity).(*entity.OrderEntity)
		return MapEntityToOrder(orderEntity)
	}
	return nil
}

func MapEntitiesToTimelineEvents(entities *entity.TimelineEventEntities) []model.TimelineEvent {
	var timelineEvents []model.TimelineEvent
	if entities == nil {
		return timelineEvents
	}
	for _, timelineEventEntity := range *entities {
		timelineEvent := MapEntityToTimelineEvent(&timelineEventEntity)
		if timelineEvent != nil {
			timelineEvents = append(timelineEvents, timelineEvent.(model.TimelineEvent))
		}
	}
	return timelineEvents
}
