package enum

import (
	"fmt"
)

type AgentID string

const (
	AgentICPQualification AgentID = "icp-qualification-agent"
	AgentVisitorID        AgentID = "visitor-identity-agent"
)

func (t AgentID) String() string {
	return string(t)
}

func GetAgent(s string) (AgentID, error) {
	switch AgentID(s) {
	case
		AgentICPQualification,
		AgentVisitorID:
		return AgentID(s), nil

	default:
		return "", fmt.Errorf("invalid Agent: %s", s)
	}
}
