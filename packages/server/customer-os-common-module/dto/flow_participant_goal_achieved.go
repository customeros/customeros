package dto

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"

type FlowParticipantGoalAchieved struct {
	ParticipantId   string           `json:"participantId"`
	ParticipantType model.EntityType `json:"participantType"`
}
