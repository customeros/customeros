package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var agentTypeByModel = map[model.AgentType]enum.AgentType{
	model.AgentTypeCampaignManager:    enum.AgentCampaignManager,
	model.AgentTypeIcpQualifier:       enum.AgentICPQualifier,
	model.AgentTypeMeetingKeeper:      enum.AgentMeetingKeeper,
	model.AgentTypeSupportSpotter:     enum.AgentSupportSpotter,
	model.AgentTypeWebVisitIdentifier: enum.AgentWebVisitorIdentifier,
	model.AgentTypeCashflowGuardian:   enum.AgentCashflowGuardian,
}

var agentTypeByValue = utils.ReverseMap(agentTypeByModel)

func MapAgentTypeFromModel(input model.AgentType) enum.AgentType {
	return agentTypeByModel[input]
}

func MapAgentTypeToModel(input enum.AgentType) model.AgentType {
	return agentTypeByValue[input]
}
