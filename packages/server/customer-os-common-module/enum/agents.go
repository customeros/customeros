package enum

import (
	"fmt"
)

type AgentType string

const (
	AgentCampaignManager      AgentType = "campaign_manager"
	AgentICPQualifier         AgentType = "icp_qualifier"
	AgentMeetingKeeper        AgentType = "meeting_keeper"
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
		AgentCampaignManager,
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
	AgentGoalReceiveReply           AgentGoal = "receive_reply"
	AgentGoalSpotHelpNeeded         AgentGoal = "spot_help_needed"
	AgentGoalGetPaid                AgentGoal = "get_paid"
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
		AgentGoalReceiveReply,
		AgentGoalSpotHelpNeeded,
		AgentGoalGetPaid:
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
