package enum

import (
	"fmt"
)

type AgentExecutionStatus string

const (
	AgentExecutionFail      AgentExecutionStatus = "ERROR"
	AgentExecutionPending   AgentExecutionStatus = "PENDING"
	AgentExecutionRunning   AgentExecutionStatus = "RUNNING"
	AgentExecutionCompleted AgentExecutionStatus = "COMPLETED"
)

func (t AgentExecutionStatus) String() string {
	return string(t)
}

func GetAgentExecutionStatus(s string) (AgentExecutionStatus, error) {
	switch AgentExecutionStatus(s) {
	case
		AgentExecutionFail,
		AgentExecutionPending,
		AgentExecutionRunning,
		AgentExecutionCompleted:
		return AgentExecutionStatus(s), nil

	default:
		return "", fmt.Errorf("invalid Agent: %s", s)
	}
}
