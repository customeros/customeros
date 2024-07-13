package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToTimelineEvent(timelineEventEntity *entity.TimelineEvent) model.TimelineEvent {
	if timelineEventEntity == nil || *timelineEventEntity == nil {
		return nil
	}
	switch (*timelineEventEntity).TimelineEventLabel() {
	case model2.NodeLabelPageView:
		pageViewEntity := (*timelineEventEntity).(*entity.PageViewEntity)
		return MapEntityToPageView(pageViewEntity)
	case model2.NodeLabelInteractionSession:
		interactionSessionEntity := (*timelineEventEntity).(*entity.InteractionSessionEntity)
		return MapEntityToInteractionSession(interactionSessionEntity)
	case model2.NodeLabelIssue:
		issueEntity := (*timelineEventEntity).(*entity.IssueEntity)
		return MapEntityToIssue(issueEntity)
	case model2.NodeLabelNote:
		noteEntity := (*timelineEventEntity).(*entity.NoteEntity)
		return MapEntityToNote(noteEntity)
	case model2.NodeLabelInteractionEvent:
		interactionEventEntity := (*timelineEventEntity).(*entity.InteractionEventEntity)
		return MapEntityToInteractionEvent(interactionEventEntity)
	case model2.NodeLabelAnalysis:
		analysisEntity := (*timelineEventEntity).(*entity.AnalysisEntity)
		return MapEntityToAnalysis(analysisEntity)
	case model2.NodeLabelMeeting:
		meetingEntity := (*timelineEventEntity).(*entity.MeetingEntity)
		return MapEntityToMeeting(meetingEntity)
	case model2.NodeLabelAction:
		actionEntity := (*timelineEventEntity).(*neo4jentity.ActionEntity)
		return MapEntityToAction(actionEntity)
	case model2.NodeLabelLogEntry:
		logEntryEntity := (*timelineEventEntity).(*neo4jentity.LogEntryEntity)
		return MapEntityToLogEntry(logEntryEntity)
	case model2.NodeLabelOrder:
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
