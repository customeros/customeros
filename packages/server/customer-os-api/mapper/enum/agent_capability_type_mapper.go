package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

var agentCapabilityTypeByModel = map[model.CapabilityType]enum.AgentCapabilityType{
	model.CapabilityTypeSendSLACkNotification:   enum.CapabilitySendSlackNotification,
	model.CapabilityTypeAnalyzeWebSessionIntent: enum.CapabilityAnalyzeWebSessionIntent,
	model.CapabilityTypeCreateOrganization:      enum.CapabilityCreateOrganization,
	model.CapabilityTypeIdentifyWebVisitor:      enum.CapabilityIdentifyWebVisitor,
}

var agentCapabilityTypeByValue = utils.ReverseMap(agentCapabilityTypeByModel)

func MapAgentCapabilityTypeFromModel(input model.CapabilityType) enum.AgentCapabilityType {
	return agentCapabilityTypeByModel[input]
}

func MapAgentCapabilityTypeToModel(input enum.AgentCapabilityType) model.CapabilityType {
	return agentCapabilityTypeByValue[input]
}
