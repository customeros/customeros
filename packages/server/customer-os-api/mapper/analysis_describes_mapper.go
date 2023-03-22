package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"reflect"
)

func MapEntityToActionDescribes(analysisDescribe *entity.AnalysisDescribe) any {
	switch (*analysisDescribe).AnalysisDescribeLabel() {
	case entity.NodeLabel_InteractionSession:
		sessionEntity := (*analysisDescribe).(*entity.InteractionSessionEntity)
		return MapEntityToInteractionSession(sessionEntity)
	case entity.NodeLabel_InteractionEvent:
		eventEntity := (*analysisDescribe).(*entity.InteractionEventEntity)
		return MapEntityToInteractionEvent(eventEntity)
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
