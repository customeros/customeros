package enum

import (
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AgentCapabilityType string

const (
	CapabilityAnalyzeWebSessionIntent         AgentCapabilityType = "analyze_web_session_for_intent"
	CapabilityCreateAndEnrichCompany          AgentCapabilityType = "create_and_enrich_company"
	CapabilityIdentifyWebVisitor              AgentCapabilityType = "identify_web_visitor"
	CapabilitySendSlackNotification           AgentCapabilityType = "send_slack_notification"
	CapabilitySendWebVisitorSlackNotification AgentCapabilityType = "send_web_visitor_slack_notification"
	CapabilityApplyTag                        AgentCapabilityType = "apply_tag"
	CapabilityCreateMarkdownTimelineEvent     AgentCapabilityType = "create_markdown_timeline_event"
	CapabilityEvaluateCompanyICPFit           AgentCapabilityType = "evaluate_company_icp_fit"
)

func (t AgentCapabilityType) String() string {
	return string(t)
}

func GetAgentCapability(s string) (AgentCapabilityType, error) {
	switch AgentCapabilityType(s) {
	case
		CapabilityAnalyzeWebSessionIntent,
		CapabilityCreateAndEnrichCompany,
		CapabilityIdentifyWebVisitor,
		CapabilitySendSlackNotification,
		CapabilitySendWebVisitorSlackNotification,
		CapabilityApplyTag,
		CapabilityCreateMarkdownTimelineEvent,
		CapabilityEvaluateCompanyICPFit:
		return AgentCapabilityType(s), nil

	default:
		return "", fmt.Errorf("invalid Agent Capability: %s", s)
	}
}

// Get name returns a friendly name for the capability, by removing _ and capitalizing the first letter of each word
func (t AgentCapabilityType) GetName() string {
	name := strings.ReplaceAll(string(t), "_", " ")
	return utils.CapitalizeAllParts(name, []string{" "})
}
