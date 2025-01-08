package enum

import (
	"fmt"
)

type Agent string

const (
	AgentICPQualification Agent = "ICP Qualification Agent"
	AgentVisitorID        Agent = "Visitor Identity Agent"
)

func (t Agent) String() string {
	return string(t)
}

func GetAgent(s string) (Agent, error) {
	switch Agent(s) {
	case
		AgentICPQualification,
		AgentVisitorID:
		return Agent(s), nil

	default:
		return "", fmt.Errorf("invalid Agent: %s", s)
	}
}
