package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"reflect"
)

func MapEntityToActionDescribes(analysisDescribe *entity.AnalysisDescribe) any {
	switch (*analysisDescribe).AnalysisDescribeLabel() {
	case model2.NodeLabelInteractionSession:
		sessionEntity := (*analysisDescribe).(*entity.InteractionSessionEntity)
		return MapEntityToInteractionSession(sessionEntity)
	case model2.NodeLabelInteractionEvent:
		eventEntity := (*analysisDescribe).(*entity.InteractionEventEntity)
		return MapEntityToInteractionEvent(eventEntity)
	case model2.NodeLabelMeeting:
		meetingEntity := (*analysisDescribe).(*entity.MeetingEntity)
		return MapEntityToMeeting(meetingEntity)
	}
	fmt.Errorf("Describes of type %s not identified", reflect.TypeOf(analysisDescribe))
	return nil
}

func MapEntitiesToDescriptionNodes(entities *entity.AnalysisDescribes) []model.DescriptionNode {
	var interactionEventParticipants []model.DescriptionNode
	for _, interactionEventParticipantEntity := range *entities {
		interactionEventParticipant := MapEntityToActionDescribes(&interactionEventParticipantEntity)
		if interactionEventParticipant != nil {
			interactionEventParticipants = append(interactionEventParticipants, interactionEventParticipant.(model.DescriptionNode))
		}
	}
	return interactionEventParticipants
}
