package enum

import (
	"fmt"
)

type AgentExecutionStatus string

const (
	AgentExecutionError     AgentExecutionStatus = "ERROR"
	AgentExecutionPending   AgentExecutionStatus = "PENDING"
	AgentExecutionRetrying  AgentExecutionStatus = "RETRYING"
	AgentExecutionRunning   AgentExecutionStatus = "RUNNING"
	AgentExecutionCompleted AgentExecutionStatus = "COMPLETED"
)

func (t AgentExecutionStatus) String() string {
	return string(t)
}

func GetAgentExecutionStatus(s string) (AgentExecutionStatus, error) {
	switch AgentExecutionStatus(s) {
	case
		AgentExecutionError,
		AgentExecutionPending,
		AgentExecutionRetrying,
		AgentExecutionRunning,
		AgentExecutionCompleted:
		return AgentExecutionStatus(s), nil

	default:
		return "", fmt.Errorf("invalid Agent: %s", s)
	}
}
