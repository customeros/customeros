package enum

import (
	"fmt"
)

type FlowAgentExecutionStatus string

const (
	FlowAgentExecutionFail    FlowAgentExecutionStatus = "FAIL"
	FlowAgentExecutionPending FlowAgentExecutionStatus = "PENDING"
	FlowAgentExecutionSuccess FlowAgentExecutionStatus = "SUCCESS"
)

func (t FlowAgentExecutionStatus) String() string {
	return string(t)
}

func GetFlowAgentExecutionStatus(s string) (FlowAgentExecutionStatus, error) {
	switch FlowAgentExecutionStatus(s) {
	case
		FlowAgentExecutionFail,
		FlowAgentExecutionPending,
		FlowAgentExecutionSuccess:
		return FlowAgentExecutionStatus(s), nil

	default:
		return "", fmt.Errorf("invalid FlowAgent: %s", s)
	}
}
