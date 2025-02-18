package enum

import (
	"fmt"
)

type CapabilityExecutionStatus string

const (
	CapabilityExecutionRetry     CapabilityExecutionStatus = "RETRY"
	CapabilityExecutionError     CapabilityExecutionStatus = "ERROR"
	CapabilityExecutionPending   CapabilityExecutionStatus = "PENDING"
	CapabilityExecutionCompleted CapabilityExecutionStatus = "COMPLETED"
	CapabilityExecutionStop      CapabilityExecutionStatus = "STOP"
)

func (t CapabilityExecutionStatus) String() string {
	return string(t)
}

func GetCapabilityExecutionStatus(s string) (CapabilityExecutionStatus, error) {
	switch CapabilityExecutionStatus(s) {
	case
		CapabilityExecutionRetry,
		CapabilityExecutionError,
		CapabilityExecutionPending,
		CapabilityExecutionStop,
		CapabilityExecutionCompleted:
		return CapabilityExecutionStatus(s), nil

	default:
		return "", fmt.Errorf("invalid CapabilityExecutionStatus: %s", s)
	}
}

type AgentCapability string

const (
	CapabilityTemplate                           AgentCapability = "template"
	CapabilityAddMeetingNotesToCompany           AgentCapability = "add_meeting_notes_to_company"
	CapabilityAnalyzeWebSessionIntent            AgentCapability = "analyze_web_session_for_intent"
	CapabilityApplyTagToCompany                  AgentCapability = "apply_tag_to_company"
	CapabilityBuildCampaignList                  AgentCapability = "build_campaign_list"
	CapabilityCollectPayments                    AgentCapability = "collect_payment"
	CapabilityCreateAndEnrichContact             AgentCapability = "create_and_enrich_contacts"
	CapabilityCreateAndEnrichCompany             AgentCapability = "create_and_enrich_company"
	CapabilityCreateMarkdownTimelineEvent        AgentCapability = "create_markdown_timeline_event"
	CapabilityDetectNewLead                      AgentCapability = "detect_new_lead"
	CapabilityDetectNewMeetingRecording          AgentCapability = "detect_new_meeting_recording"
	CapabilityDetectSupportWebvisit              AgentCapability = "detect_support_webvisit"
	CapabilityEnrichEmailAddress                 AgentCapability = "enrich_email_address"
	CapabilityEvaluateCompanyICPFit              AgentCapability = "evaluate_company_icp_fit"
	CapabilityExtractMeetingHighlights           AgentCapability = "extract_meeting_highlights"
	CapabilityExtractSupportSignalsFromMeeting   AgentCapability = "extract_support_signals_from_meeting"
	CapabilityForwardEmailReply                  AgentCapability = "forward_email_reply"
	CapabilityGatherCompanyIntelligence          AgentCapability = "gather_company_intelligence"
	CapabilityGenerateInvoice                    AgentCapability = "generate_invoice"
	CapabilityHandleAutoresponder                AgentCapability = "handle_autoresponder"
	CapabilityIdentifyMeetingParticipants        AgentCapability = "identify_meeting_participants"
	CapabilityIdentifyWebVisitor                 AgentCapability = "identify_web_visitor"
	CapabilityLogRequestsForHelp                 AgentCapability = "log_requests_for_help"
	CapabilityManageBouncedEmail                 AgentCapability = "manage_bounced_email"
	CapabilityManageCampaignExecution            AgentCapability = "manage_campaign_execution"
	CapabilityManageEmailDeliveryFailure         AgentCapability = "manage_email_delivery_failure"
	CapabilityMonitorAccountsReceivable          AgentCapability = "monitor_accounts_receivable"
	CapabilityMonitorSupportVisits               AgentCapability = "monitor_support_visits"
	CapabilityProcessRefund                      AgentCapability = "process_refund"
	CapabilityScheduleEmailDelivery              AgentCapability = "schedule_email_delivery"
	CapabilitySelectOptimalSendingMailbox        AgentCapability = "select_optimal_sending_mailbox"
	CapabilitySendLinkedinConnectionRequest      AgentCapability = "send_linkedin_connection_request"
	CapabilitySendLinkedinMessage                AgentCapability = "send_linkedin_message"
	CapabilitySyncLinkedinConnections            AgentCapability = "sync_linkedin_connections"
	CapabilitySendPaymentReminder                AgentCapability = "send_payment_reminder"
	CapabilitySendSlackNotification              AgentCapability = "send_slack_notification"
	CapabilitySendWebVisitorSlackNotification    AgentCapability = "send_web_visitor_slack_notification"
	CapabilitySyncWithAccountingSystem           AgentCapability = "sync_with_accounting_system"
	CapabilityTrackCampaignEngagement            AgentCapability = "track_campaign_engagement"
	CapabilityUpdateCompanyStatus                AgentCapability = "update_company_status"
	CapabilityValidateEmailAddressDeliverability AgentCapability = "validate_email_deliverability"
	CapabilitySendInvoiceViaEmail                AgentCapability = "send_invoice_via_email"
	CapabilityProcessAutopayment                 AgentCapability = "process_autopayment"
	CapabilityClassifyEmail                      AgentCapability = "classify_email"
	CapabilityIdentifyEmailParticipants          AgentCapability = "identify_email_participants"
	CapabilitySummarizeMessage                   AgentCapability = "summarize_message"
	CapabilitySummarizeThread                    AgentCapability = "summarize_thread"
	CapabilityIngestEmail                        AgentCapability = "ingest_email"
	CapabilitySendPaidNotification               AgentCapability = "send_paid_notification"
	CapabilitySendInvoiceVoidedNotification      AgentCapability = "send_invoice_voided_notification"
	CapabilitySendPastDueNotification            AgentCapability = "send_past_due_notification"
)

func (t AgentCapability) String() string {
	return string(t)
}

func GetAgentCapability(s string) (AgentCapability, error) {
	switch AgentCapability(s) {
	case
		CapabilityTemplate,
		CapabilityAddMeetingNotesToCompany,
		CapabilityAnalyzeWebSessionIntent,
		CapabilityApplyTagToCompany,
		CapabilityBuildCampaignList,
		CapabilityCollectPayments,
		CapabilityCreateAndEnrichContact,
		CapabilityCreateAndEnrichCompany,
		CapabilityCreateMarkdownTimelineEvent,
		CapabilityDetectNewLead,
		CapabilityDetectNewMeetingRecording,
		CapabilityDetectSupportWebvisit,
		CapabilityEnrichEmailAddress,
		CapabilityEvaluateCompanyICPFit,
		CapabilityExtractMeetingHighlights,
		CapabilityExtractSupportSignalsFromMeeting,
		CapabilityForwardEmailReply,
		CapabilityGatherCompanyIntelligence,
		CapabilityGenerateInvoice,
		CapabilityHandleAutoresponder,
		CapabilityIdentifyMeetingParticipants,
		CapabilityIdentifyWebVisitor,
		CapabilityLogRequestsForHelp,
		CapabilityManageBouncedEmail,
		CapabilityManageCampaignExecution,
		CapabilityManageEmailDeliveryFailure,
		CapabilityMonitorAccountsReceivable,
		CapabilityMonitorSupportVisits,
		CapabilityProcessRefund,
		CapabilityScheduleEmailDelivery,
		CapabilitySelectOptimalSendingMailbox,
		CapabilitySendLinkedinConnectionRequest,
		CapabilitySendLinkedinMessage,
		CapabilitySyncLinkedinConnections,
		CapabilitySendPaymentReminder,
		CapabilitySendSlackNotification,
		CapabilitySendWebVisitorSlackNotification,
		CapabilitySyncWithAccountingSystem,
		CapabilityTrackCampaignEngagement,
		CapabilityUpdateCompanyStatus,
		CapabilitySendInvoiceViaEmail,
		CapabilityValidateEmailAddressDeliverability,
		CapabilityProcessAutopayment,
		CapabilityClassifyEmail,
		CapabilityIdentifyEmailParticipants,
		CapabilitySummarizeMessage,
		CapabilitySummarizeThread,
		CapabilityIngestEmail,
		CapabilitySendPaidNotification,
		CapabilitySendInvoiceVoidedNotification,
		CapabilitySendPastDueNotification:
		return AgentCapability(s), nil

	default:
		return "", fmt.Errorf("invalid Agent Capability: %s", s)
	}
}
