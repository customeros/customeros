package dto

import "github.com/customeros/customeros/packages/server/customer-os-common-module/model"

type FlowParticipantGoalAchieved struct {
	ParticipantId   string           `json:"participantId"`
	ParticipantType model.EntityType `json:"participantType"`
}
