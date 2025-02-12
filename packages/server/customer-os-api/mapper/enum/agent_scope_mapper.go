package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var agentScopeByModel = map[model.AgentScope]enum.AgentScope{
	model.AgentScopePersonal:  enum.AgentScopePersonal,
	model.AgentScopeWorkspace: enum.AgentScopeWorkspace,
}

var agentScopeByValue = utils.ReverseMap(agentScopeByModel)

func MapAgentScopeFromModel(input model.AgentScope) enum.AgentScope {
	return agentScopeByModel[input]
}

func MapAgentScopeToModel(input enum.AgentScope) model.AgentScope {
	return agentScopeByValue[input]
}
