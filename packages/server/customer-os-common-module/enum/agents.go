package enum

import (
	"fmt"
)

type AgentType string

const (
	AgentICPQualification AgentType = "icp_qualification_agent"
	AgentSupport          AgentType = "support_agent"
	AgentVisitorID        AgentType = "visitor_identity_agent"
)

func (t AgentType) String() string {
	return string(t)
}

func GetAgentType(s string) (AgentType, error) {
	switch AgentType(s) {
	case
		AgentICPQualification,
		AgentSupport,
		AgentVisitorID:
		return AgentType(s), nil

	default:
		return "", fmt.Errorf("invalid Agent: %s", s)
	}
}

type AgentGoal string

const (
	AgentGoalIdentifyVisitors AgentGoal = "identify_web_visitor"
	AgentGoalTagSupport       AgentGoal = "tag_support"
)

func (t AgentGoal) String() string {
	return string(t)
}
