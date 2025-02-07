package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var agentListenerEventByModel = map[model.AgentListenerEvent]enum.AgentListenerEvent{
	model.AgentListenerEventIcpFit:                  enum.EventICPFit,
	model.AgentListenerEventIcpNotAFit:              enum.EventICPNotAFit,
	model.AgentListenerEventNewLead:                 enum.EventNewLead,
	model.AgentListenerEventRunIcpQualifierAgent:    enum.EventRunICPQualifierAgent,
	model.AgentListenerEventNewWebSession:           enum.EventNewWebSession,
	model.AgentListenerEventWebVisitorIdentified:    enum.EventWebVisitorIdentified,
	model.AgentListenerEventWebVisitorNotIdentified: enum.EventWebVisitorNotIdentified,
}

var agentListenerEventByValue = utils.ReverseMap(agentListenerEventByModel)

func MapAgentListenerTypeFromModel(input model.AgentListenerEvent) enum.AgentListenerEvent {
	return agentListenerEventByModel[input]
}

func MapAgentListenerTypeToModel(input enum.AgentListenerEvent) model.AgentListenerEvent {
	return agentListenerEventByValue[input]
}
