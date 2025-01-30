package enum

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"strings"
)

type AgentCapabilityType string

const (
	CapabilityAnalyzeWebSessionIntent         AgentCapabilityType = "analyze_web_session_for_intent"
	CapabilityCreateOrganization              AgentCapabilityType = "create_organization"
	CapabilityIdentifyWebVisitor              AgentCapabilityType = "identify_web_visitor"
	CapabilitySendSlackNotification           AgentCapabilityType = "send_slack_notification"
	CapabilitySendWebVisitorSlackNotification AgentCapabilityType = "send_web_visitor_slack_notification"
	CapabilityApplyTag                        AgentCapabilityType = "apply_tag"
)

func (t AgentCapabilityType) String() string {
	return string(t)
}

func GetAgentCapability(s string) (AgentCapabilityType, error) {
	switch AgentCapabilityType(s) {
	case
		CapabilityAnalyzeWebSessionIntent,
		CapabilityCreateOrganization,
		CapabilityIdentifyWebVisitor,
		CapabilitySendSlackNotification,
		CapabilitySendWebVisitorSlackNotification,
		CapabilityApplyTag:
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
