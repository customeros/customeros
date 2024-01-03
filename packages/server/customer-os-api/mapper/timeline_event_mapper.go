package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToTimelineEvent(timelineEventEntity *entity.TimelineEvent) model.TimelineEvent {
	if timelineEventEntity == nil || *timelineEventEntity == nil {
		return nil
	}
	switch (*timelineEventEntity).TimelineEventLabel() {
	case neo4jentity.NodeLabel_PageView:
		pageViewEntity := (*timelineEventEntity).(*entity.PageViewEntity)
		return MapEntityToPageView(pageViewEntity)
	case neo4jentity.NodeLabel_InteractionSession:
		interactionSessionEntity := (*timelineEventEntity).(*entity.InteractionSessionEntity)
		return MapEntityToInteractionSession(interactionSessionEntity)
	case neo4jentity.NodeLabel_Issue:
		issueEntity := (*timelineEventEntity).(*entity.IssueEntity)
		return MapEntityToIssue(issueEntity)
	case neo4jentity.NodeLabel_Note:
		noteEntity := (*timelineEventEntity).(*entity.NoteEntity)
		return MapEntityToNote(noteEntity)
	case neo4jentity.NodeLabel_InteractionEvent:
		interactionEventEntity := (*timelineEventEntity).(*entity.InteractionEventEntity)
		return MapEntityToInteractionEvent(interactionEventEntity)
	case neo4jentity.NodeLabel_Analysis:
		analysisEntity := (*timelineEventEntity).(*entity.AnalysisEntity)
		return MapEntityToAnalysis(analysisEntity)
	case neo4jentity.NodeLabel_Meeting:
		meetingEntity := (*timelineEventEntity).(*entity.MeetingEntity)
		return MapEntityToMeeting(meetingEntity)
	case neo4jentity.NodeLabel_Action:
		actionEntity := (*timelineEventEntity).(*entity.ActionEntity)
		return MapEntityToAction(actionEntity)
	case neo4jentity.NodeLabel_LogEntry:
		logEntryEntity := (*timelineEventEntity).(*entity.LogEntryEntity)
		return MapEntityToLogEntry(logEntryEntity)
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
