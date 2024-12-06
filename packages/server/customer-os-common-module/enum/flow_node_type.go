package enum

import "fmt"

type FlowNodeType string

const (
	NodeFlowAction        FlowNodeType = "flowAction"
	NodeFlowListenerEvent FlowNodeType = "listenerEvent"
)

func (e FlowNodeType) String() string {
	return string(e)
}

func GetFlowNodeType(s string) (FlowNodeType, error) {
	switch FlowNodeType(s) {
	case
		NodeFlowAction,
		NodeFlowListenerEvent:

		return FlowNodeType(s), nil
	default:
		return "", fmt.Errorf("invalid FlowNodeType: %s", s)
	}
}
