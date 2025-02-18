package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var agentCapabilityTypeByModel = map[model.CapabilityType]enum.AgentCapability{
	model.CapabilityTypeAddMeetingNotesToCompany:         enum.CapabilityAddMeetingNotesToCompany,
	model.CapabilityTypeAnalyzeWebSessionIntent:          enum.CapabilityAnalyzeWebSessionIntent,
	model.CapabilityTypeApplyTagToCompany:                enum.CapabilityApplyTagToCompany,
	model.CapabilityTypeCreateContacts:                   enum.CapabilityCreateAndEnrichContact,
	model.CapabilityTypeCreateMarkdownTimelineEvent:      enum.CapabilityCreateMarkdownTimelineEvent,
	model.CapabilityTypeCreateOrganization:               enum.CapabilityCreateAndEnrichCompany,
	model.CapabilityTypeDetectSupportWebvisit:            enum.CapabilityDetectSupportWebvisit,
	model.CapabilityTypeEnrichEmailAddress:               enum.CapabilityEnrichEmailAddress,
	model.CapabilityTypeExtractMeetingHighlights:         enum.CapabilityExtractMeetingHighlights,
	model.CapabilityTypeExtractSupportSignalsFromMeeting: enum.CapabilityExtractSupportSignalsFromMeeting,
	model.CapabilityTypeForwardEmailReply:                enum.CapabilityForwardEmailReply,
	model.CapabilityTypeGatherCompanyIntelligence:        enum.CapabilityGatherCompanyIntelligence,
	model.CapabilityTypeIcpQualify:                       enum.CapabilityEvaluateCompanyICPFit,
	model.CapabilityTypeIdentifyMeetingParticipants:      enum.CapabilityIdentifyMeetingParticipants,
	model.CapabilityTypeIdentifyWebVisitor:               enum.CapabilityIdentifyWebVisitor,
	model.CapabilityTypeLogRequestsForHelp:               enum.CapabilityLogRequestsForHelp,
	model.CapabilityTypeManageCampaignExecution:          enum.CapabilityManageCampaignExecution,
	model.CapabilityTypeManageEmailDeliveryFailure:       enum.CapabilityManageEmailDeliveryFailure,
	model.CapabilityTypeSelectOptimalSendingMailbox:      enum.CapabilitySelectOptimalSendingMailbox,
	model.CapabilityTypeSendSLACkNotification:            enum.CapabilitySendSlackNotification,
	model.CapabilityTypeUpdateCompanyStatus:              enum.CapabilityUpdateCompanyStatus,
	model.CapabilityTypeValidateEmailDeliverability:      enum.CapabilityValidateEmailAddressDeliverability,
	model.CapabilityTypeWebVisitorSendSLACkNotification:  enum.CapabilitySendWebVisitorSlackNotification,
	model.CapabilityTypeSendInvoiceViaEmail:              enum.CapabilitySendInvoiceViaEmail,
	model.CapabilityTypeGenerateInvoice:                  enum.CapabilityGenerateInvoice,
	model.CapabilityTypeProcessAutopayment:               enum.CapabilityProcessAutopayment,
	model.CapabilityTypeCreayePaymentLink:                enum.CapabilityCreatePaymentLink,
	model.CapabilityTypeSendPaidNotification:             enum.CapabilitySendPaidNotification,
	model.CapabilityTypeSendInvoiceVoidedNotification:    enum.CapabilitySendInvoiceVoidedNotification,
	model.CapabilityTypeClassifyEmail:                    enum.CapabilityClassifyEmail,
	model.CapabilityTypeIdentifyParticipants:             enum.CapabilityIdentifyEmailParticipants,
	model.CapabilityTypeSummarizeMessage:                 enum.CapabilitySummarizeMessage,
	model.CapabilityTypeSummarizeThread:                  enum.CapabilitySummarizeThread,
	model.CapabilityTypeIngestEmail:                      enum.CapabilityIngestEmail,
	model.CapabilityTypeSendPastDueNotification:          enum.CapabilitySendPastDueNotification,
}

var agentCapabilityTypeByValue = utils.ReverseMap(agentCapabilityTypeByModel)

func MapAgentCapabilityTypeFromModel(input model.CapabilityType) enum.AgentCapability {
	return agentCapabilityTypeByModel[input]
}

func MapAgentCapabilityTypeToModel(input enum.AgentCapability) model.CapabilityType {
	return agentCapabilityTypeByValue[input]
}
