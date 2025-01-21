package enum

import (
	"fmt"
)

type AgentCapabilityType string

const (
	CapabilityAnalyzeWebSessionIntent AgentCapabilityType = "analyze_web_session_for_intent"
	CapabilityCreateOrganization      AgentCapabilityType = "create_organization"
	CapabilityIdentifyWebVisitor      AgentCapabilityType = "identify_web_visitor"
	CapabilitySendSlackNotification   AgentCapabilityType = "send_slack_notification"
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
		CapabilitySendSlackNotification:
		return AgentCapabilityType(s), nil

	default:
		return "", fmt.Errorf("invalid Agent Capability: %s", s)
	}
}
