package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"reflect"
)

func MapEntityToTimelineEvent(timelineEventEntity *entity.TimelineEvent) model.TimelineEvent {
	if timelineEventEntity == nil || *timelineEventEntity == nil {
		return nil
	}
	switch (*timelineEventEntity).TimelineEventLabel() {
	case entity.NodeLabel_PageView:
		pageViewEntity := (*timelineEventEntity).(*entity.PageViewEntity)
		return MapEntityToPageView(pageViewEntity)
	case entity.NodeLabel_InteractionSession:
		interactionSessionEntity := (*timelineEventEntity).(*entity.InteractionSessionEntity)
		return MapEntityToInteractionSession(interactionSessionEntity)
	case entity.NodeLabel_Issue:
		issueEntity := (*timelineEventEntity).(*entity.IssueEntity)
		return MapEntityToIssue(issueEntity)
	case entity.NodeLabel_Note:
		noteEntity := (*timelineEventEntity).(*entity.NoteEntity)
		return MapEntityToNote(noteEntity)
	case entity.NodeLabel_InteractionEvent:
		interactionEventEntity := (*timelineEventEntity).(*entity.InteractionEventEntity)
		return MapEntityToInteractionEvent(interactionEventEntity)
	case entity.NodeLabel_Analysis:
		analysisEntity := (*timelineEventEntity).(*entity.AnalysisEntity)
		return MapEntityToAnalysis(analysisEntity)
	case entity.NodeLabel_Meeting:
		meetingEntity := (*timelineEventEntity).(*entity.MeetingEntity)
		return MapEntityToMeeting(meetingEntity)
	case entity.NodeLabel_Action:
		actionEntity := (*timelineEventEntity).(*entity.ActionEntity)
		return MapEntityToAction(actionEntity)
	}
	fmt.Errorf("timeline event of type %s not identified", reflect.TypeOf(timelineEventEntity))
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
