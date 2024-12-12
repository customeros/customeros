package enum

import (
	"fmt"
)

type FlowActionExecutionStatus string

const (
	FlowActionExecutionFail    FlowActionExecutionStatus = "FAIL"
	FlowActionExecutionPending FlowActionExecutionStatus = "PENDING"
	FlowActionExecutionSuccess FlowActionExecutionStatus = "SUCCESS"
)

func (t FlowActionExecutionStatus) String() string {
	return string(t)
}

func GetFlowActionExecutionStatus(s string) (FlowActionExecutionStatus, error) {
	switch FlowActionExecutionStatus(s) {
	case
		FlowActionExecutionFail,
		FlowActionExecutionPending,
		FlowActionExecutionSuccess:
		return FlowActionExecutionStatus(s), nil

	default:
		return "", fmt.Errorf("invalid FlowAction: %s", s)
	}
}
