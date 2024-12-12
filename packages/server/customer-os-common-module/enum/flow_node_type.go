package enum

import "fmt"

type FlowNodeType string

const (
	NodeFlowAction        FlowNodeType = "action"
	NodeFlowEnd           FlowNodeType = "end"
	NodeFlowListenerEvent FlowNodeType = "listener"
	NodeFlowWait          FlowNodeType = "wait"
)

func (e FlowNodeType) String() string {
	return string(e)
}

func GetFlowNodeType(s string) (FlowNodeType, error) {
	switch FlowNodeType(s) {
	case
		NodeFlowAction,
		NodeFlowEnd,
		NodeFlowListenerEvent,
		NodeFlowWait:

		return FlowNodeType(s), nil
	default:
		return "", fmt.Errorf("invalid FlowNodeType: %s", s)
	}
}
