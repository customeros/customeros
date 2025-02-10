package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var agentListenerEventByModel = map[model.AgentListenerEvent]enum.AgentListenerEvent{
	model.AgentListenerEventCompanyIdentified:       enum.EventCompanyIdentified,
	model.AgentListenerEventIcpFit:                  enum.EventICPFit,
	model.AgentListenerEventIcpNotAFit:              enum.EventICPNotAFit,
	model.AgentListenerEventNewLead:                 enum.EventNewLead,
	model.AgentListenerEventNewMeetingRecording:     enum.EventNewMeetingRecording,
	model.AgentListenerEventNewWebSession:           enum.EventNewWebSession,
	model.AgentListenerEventRunIcpQualifierAgent:    enum.EventRunICPQualifierAgent,
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
