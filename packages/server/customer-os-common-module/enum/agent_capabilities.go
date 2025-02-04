package enum

import (
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AgentCapability string

const (
	CapabilityAnalyzeWebSessionIntent            AgentCapability = "analyze_web_session_for_intent"
	CapabilityApplyTag                           AgentCapability = "apply_tag"
	CapabilityBuildCampaignList                  AgentCapability = "build_campaign_list"
	CapabilityCollectPayments                    AgentCapability = "collect_payment"
	CapabilityCreateAndEnrichCompany             AgentCapability = "create_and_enrich_company"
	CapabilityCreateMarkdownTimelineEvent        AgentCapability = "create_markdown_timeline_event"
	CapabilityDetectNewLead                      AgentCapability = "detect_new_lead"
	CapabilityDetectNewMeetingRecording          AgentCapability = "detect_new_meeting_recording"
	CapabilityEnrichEmailAddress                 AgentCapability = "enrich_email_address"
	CapabilityEvaluateCompanyICPFit              AgentCapability = "evaluate_company_icp_fit"
	CapabilityExtractMeetingHighlights           AgentCapability = "extract_meeting_highlights"
	CapabilityForwardEmailReply                  AgentCapability = "forward_email_reply"
	CapabilityGatherCompanyIntelligence          AgentCapability = "gather_company_intelligence"
	CapabilityGenerateInvoice                    AgentCapability = "generate_invoice"
	CapabilityHandleAutoresponder                AgentCapability = "handle_autoresponder"
	CapabilityIdentifyMeetingParticipants        AgentCapability = "identify_meeting_participants"
	CapabilityIdentifyWebVisitor                 AgentCapability = "identify_web_visitor"
	CapabilityManageBouncedEmail                 AgentCapability = "manage_bounced_email"
	CapabilityMonitorAccountsReceivable          AgentCapability = "monitor_accounts_receivable"
	CapabilityMonitorSupportPages                AgentCapability = "monitor_support_pages"
	CapabilityProcessRefund                      AgentCapability = "process_refund"
	CapabilityScheduleEmailDelivery              AgentCapability = "schedule_email_delivery"
	CapabilitySelectOptimalSendingMailbox        AgentCapability = "select_optimal_sending_mailbox"
	CapabilitySendInvoice                        AgentCapability = "send_invoice"
	CapabilitySendLinkedinConnectionRequest      AgentCapability = "send_linkedin_connection_request"
	CapabililtySendLinkedinMessage               AgentCapability = "send_linkedin_message"
	CapabilitySyncLinkedinConnections            AgentCapability = "sync_linkedin_connections"
	CapabilitySendPaymentReminder                AgentCapability = "send_payment_reminder"
	CapabilitySendSlackNotification              AgentCapability = "send_slack_notification"
	CapabilitySendWebVisitorSlackNotification    AgentCapability = "send_web_visitor_slack_notification"
	CapabilitySyncWithAccountingSystem           AgentCapability = "sync_with_accounting_system"
	CapabilityTrackCampaignEngagement            AgentCapability = "track_campaign_engagement"
	CapabilityUpdateCompanyStatus                AgentCapability = "update_company_status"
	CapabilityValidateEmailAddressDeliverability AgentCapability = "validate_email_address_deliverability"
)

func (t AgentCapability) String() string {
	return string(t)
}

func GetAgentCapability(s string) (AgentCapability, error) {
	switch AgentCapability(s) {
	case
		CapabilityAnalyzeWebSessionIntent,
		CapabilityApplyTag,
		CapabilityBuildCampaignList,
		CapabilityCollectPayments,
		CapabilityCreateAndEnrichCompany,
		CapabilityCreateMarkdownTimelineEvent,
		CapabilityDetectNewLead,
		CapabilityDetectNewMeetingRecording,
		CapabilityEnrichEmailAddress,
		CapabilityEvaluateCompanyICPFit,
		CapabilityExtractMeetingHighlights,
		CapabilityForwardEmailReply,
		CapabilityGatherCompanyIntelligence,
		CapabilityGenerateInvoice,
		CapabilityHandleAutoresponder,
		CapabilityIdentifyMeetingParticipants,
		CapabilityIdentifyWebVisitor,
		CapabilityManageBouncedEmail,
		CapabilityMonitorAccountsReceivable,
		CapabilityMonitorSupportPages,
		CapabilityProcessRefund,
		CapabilityScheduleEmailDelivery,
		CapabilitySelectOptimalSendingMailbox,
		CapabilitySendInvoice,
		CapabilitySendLinkedinConnectionRequest,
		CapabililtySendLinkedinMessage,
		CapabilitySyncLinkedinConnections,
		CapabilitySendPaymentReminder,
		CapabilitySendSlackNotification,
		CapabilitySendWebVisitorSlackNotification,
		CapabilitySyncWithAccountingSystem,
		CapabilityTrackCampaignEngagement,
		CapabilityUpdateCompanyStatus,
		CapabilityValidateEmailAddressDeliverability:
		return AgentCapability(s), nil

	default:
		return "", fmt.Errorf("invalid Agent Capability: %s", s)
	}
}

// Get name returns a friendly name for the capability, by removing _ and capitalizing the first letter of each word
func (t AgentCapability) GetName() string {
	name := strings.ReplaceAll(string(t), "_", " ")
	return utils.CapitalizeAllParts(name, []string{" "})
}
