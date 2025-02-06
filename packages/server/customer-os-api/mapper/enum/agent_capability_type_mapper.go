package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var agentCapabilityTypeByModel = map[model.CapabilityType]enum.AgentCapability{
	model.CapabilityTypeSendSLACkNotification:           enum.CapabilitySendSlackNotification,
	model.CapabilityTypeAnalyzeWebSessionIntent:         enum.CapabilityAnalyzeWebSessionIntent,
	model.CapabilityTypeCreateOrganization:              enum.CapabilityCreateAndEnrichCompany,
	model.CapabilityTypeIdentifyWebVisitor:              enum.CapabilityIdentifyWebVisitor,
	model.CapabilityTypeWebVisitorSendSLACkNotification: enum.CapabilitySendWebVisitorSlackNotification,
	model.CapabilityTypeApplyTag:                        enum.CapabilityApplyTag,
	model.CapabilityTypeCreateMarkdownTimelineEvent:     enum.CapabilityCreateMarkdownTimelineEvent,
	model.CapabilityTypeIcpQualify:                      enum.CapabilityEvaluateCompanyICPFit,
	model.CapabilityTypeGatherCompanyIntelligence:       enum.CapabilityGatherCompanyIntelligence,
	model.CapabilityTypeUpdateCompanyStatus:             enum.CapabilityUpdateCompanyStatus,
}

var agentCapabilityTypeByValue = utils.ReverseMap(agentCapabilityTypeByModel)

func MapAgentCapabilityTypeFromModel(input model.CapabilityType) enum.AgentCapability {
	return agentCapabilityTypeByModel[input]
}

func MapAgentCapabilityTypeToModel(input enum.AgentCapability) model.CapabilityType {
	return agentCapabilityTypeByValue[input]
}
