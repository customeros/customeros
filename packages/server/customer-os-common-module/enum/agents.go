package enum

import (
	"fmt"
)

type AgentType string

const (
	AgentICPQualifier         AgentType = "icp_qualifier"
	AgentSupportSpotter       AgentType = "support_spotter"
	AgentWebVisitorIdentifier AgentType = "web_visitor_identifier"
	AgentCashflowGuardian     AgentType = "cashflow_guardian"
)

func (t AgentType) String() string {
	return string(t)
}

func GetAgentType(s string) (AgentType, error) {
	switch AgentType(s) {
	case
		AgentICPQualifier,
		AgentSupportSpotter,
		AgentCashflowGuardian,
		AgentWebVisitorIdentifier:
		return AgentType(s), nil

	default:
		return "", fmt.Errorf("invalid Agent: %s", s)
	}
}

type AgentGoal string

const (
	AgentGoalEvaluateICPFit     AgentGoal = "evaluate_icp_fit"
	AgentGoalIdentifyWebVisitor AgentGoal = "identify_web_visitor"
	AgentGoalSpotHelpNeeded     AgentGoal = "spot_help_needed"
)

func (t AgentGoal) String() string {
	return string(t)
}

func GetAgentGoal(s string) (AgentGoal, error) {
	switch AgentGoal(s) {
	case
		AgentGoalEvaluateICPFit,
		AgentGoalIdentifyWebVisitor,
		AgentGoalSpotHelpNeeded:
		return AgentGoal(s), nil

	default:
		return "", fmt.Errorf("invalid Agent Goal: %s", s)
	}
}
