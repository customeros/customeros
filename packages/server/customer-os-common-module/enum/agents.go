package enum

import (
	"fmt"
)

type AgentType string

const (
	AgentICPQualifier          AgentType = "icp_qualifier"
	AgentSupportSignalDetector AgentType = "support_signal_detector"
	AgentWebVisitorIdentifier  AgentType = "web_visitor_identifier"
)

func (t AgentType) String() string {
	return string(t)
}

func GetAgentType(s string) (AgentType, error) {
	switch AgentType(s) {
	case
		AgentICPQualifier,
		AgentSupportSignalDetector,
		AgentWebVisitorIdentifier:
		return AgentType(s), nil

	default:
		return "", fmt.Errorf("invalid Agent: %s", s)
	}
}

type AgentGoal string

const (
	AgentGoalIdentifyWebVisitor  AgentGoal = "identify_web_visitor"
	AgentGoalDetectSupportSignal AgentGoal = "detect_support_signal"
	AgentGoalEvaluateICPFit      AgentGoal = "evaluate_icp_fit"
)

func (t AgentGoal) String() string {
	return string(t)
}
