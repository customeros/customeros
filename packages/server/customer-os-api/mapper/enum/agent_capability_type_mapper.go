package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var agentCapabilityTypeByModel = map[model.CapabilityType]enum.AgentCapability{
	model.CapabilityTypeAnalyzeWebSessionIntent:         enum.CapabilityAnalyzeWebSessionIntent,
	model.CapabilityTypeApplyTagToCompany:               enum.CapabilityApplyTagToCompany,
	model.CapabilityTypeCreateContacts:                  enum.CapabilityCreateAndEnrichContact,
	model.CapabilityTypeCreateMarkdownTimelineEvent:     enum.CapabilityCreateMarkdownTimelineEvent,
	model.CapabilityTypeCreateOrganization:              enum.CapabilityCreateAndEnrichCompany,
	model.CapabilityTypeExtractMeetingHighlights:        enum.CapabilityExtractMeetingHighlights,
	model.CapabilityTypeGatherCompanyIntelligence:       enum.CapabilityGatherCompanyIntelligence,
	model.CapabilityTypeIcpQualify:                      enum.CapabilityEvaluateCompanyICPFit,
	model.CapabilityTypeIdentifyMeetingParticipants:     enum.CapabilityIdentifyMeetingParticipants,
	model.CapabilityTypeIdentifyWebVisitor:              enum.CapabilityIdentifyWebVisitor,
	model.CapabilityTypeSendSLACkNotification:           enum.CapabilitySendSlackNotification,
	model.CapabilityTypeUpdateCompanyStatus:             enum.CapabilityUpdateCompanyStatus,
	model.CapabilityTypeWebVisitorSendSLACkNotification: enum.CapabilitySendWebVisitorSlackNotification,
}

var agentCapabilityTypeByValue = utils.ReverseMap(agentCapabilityTypeByModel)

func MapAgentCapabilityTypeFromModel(input model.CapabilityType) enum.AgentCapability {
	return agentCapabilityTypeByModel[input]
}

func MapAgentCapabilityTypeToModel(input enum.AgentCapability) model.CapabilityType {
	return agentCapabilityTypeByValue[input]
}
