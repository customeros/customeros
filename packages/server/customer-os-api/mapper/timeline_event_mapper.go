package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToTimelineEvent(timelineEventEntity *entity.TimelineEvent) model.TimelineEvent {
	if timelineEventEntity == nil || *timelineEventEntity == nil {
		return nil
	}
	switch (*timelineEventEntity).TimelineEventLabel() {
	case commonmodel.NodeLabelIssue:
		issueEntity := (*timelineEventEntity).(*entity.IssueEntity)
		return MapEntityToIssue(issueEntity)
	case commonmodel.NodeLabelNote:
		noteEntity := (*timelineEventEntity).(*entity.NoteEntity)
		return MapEntityToNote(noteEntity)
	case commonmodel.NodeLabelInteractionEvent:
		interactionEventEntity := (*timelineEventEntity).(*neo4jentity.InteractionEventEntity)
		return MapEntityToInteractionEvent(interactionEventEntity)
	case commonmodel.NodeLabelMeeting:
		meetingEntity := (*timelineEventEntity).(*entity.MeetingEntity)
		return MapEntityToMeeting(meetingEntity)
	case commonmodel.NodeLabelAction:
		actionEntity := (*timelineEventEntity).(*neo4jentity.ActionEntity)
		return MapEntityToAction(actionEntity)
	case commonmodel.NodeLabelLogEntry:
		logEntryEntity := (*timelineEventEntity).(*neo4jentity.LogEntryEntity)
		return MapEntityToLogEntry(logEntryEntity)
	case commonmodel.NodeLabelMarkdownEvent:
		markdownEventEntity := (*timelineEventEntity).(*neo4jentity.MarkdownEventEntity)
		return MapEntityToMarkdownEvent(markdownEventEntity)
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
