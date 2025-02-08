package enum

import (
	"fmt"
)

type AgentType string

const (
	AgentCashflowGuardian     AgentType = "cashflow_guardian"
	AgentICPQualifier         AgentType = "icp_qualifier"
	AgentMeetingKeeper        AgentType = "meeting_keeper"
	AgentSupportSpotter       AgentType = "support_spotter"
	AgentWebVisitorIdentifier AgentType = "web_visitor_identifier"
)

func (t AgentType) String() string {
	return string(t)
}

func GetAgentType(s string) (AgentType, error) {
	switch AgentType(s) {
	case
		AgentCashflowGuardian,
		AgentMeetingKeeper,
		AgentICPQualifier,
		AgentSupportSpotter,
		AgentWebVisitorIdentifier:
		return AgentType(s), nil

	default:
		return "", fmt.Errorf("invalid Agent: %s", s)
	}
}

type AgentGoal string

const (
	AgentGoalCaptureExternalMeeting AgentGoal = "capture_external_meeting"
	AgentGoalEvaluateICPFit         AgentGoal = "evaluate_icp_fit"
	AgentGoalIdentifyWebVisitor     AgentGoal = "identify_web_visitor"
	AgentGoalSpotHelpNeeded         AgentGoal = "spot_help_needed"
)

func (t AgentGoal) String() string {
	return string(t)
}

func GetAgentGoal(s string) (AgentGoal, error) {
	switch AgentGoal(s) {
	case
		AgentGoalCaptureExternalMeeting,
		AgentGoalEvaluateICPFit,
		AgentGoalIdentifyWebVisitor,
		AgentGoalSpotHelpNeeded:
		return AgentGoal(s), nil

	default:
		return "", fmt.Errorf("invalid Agent Goal: %s", s)
	}
}

type AgentScope string

const (
	AgentScopePersonal  AgentScope = "personal"
	AgentScopeWorkspace AgentScope = "workspace"
)

func (t AgentScope) String() string {
	return string(t)
}

func GetAgentScope(s string) (AgentScope, error) {
	switch AgentScope(s) {
	case
		AgentScopePersonal,
		AgentScopeWorkspace:
		return AgentScope(s), nil

	default:
		return "", fmt.Errorf("invalid Agent Scope: %s", s)
	}
}
