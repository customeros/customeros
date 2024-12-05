package enum

import "fmt"

type FlowAction string

const (
	ActionCreateTimelineEvent FlowAction = "create.timeline_event"
	ActionCreateContact       FlowAction = "create.contact"
	ActionCreateOrganization  FlowAction = "create_organization"
)

func (t FlowAction) String() string {
	return string(t)
}

func GetFlowAction(s string) (FlowAction, error) {
	switch FlowAction(s) {
	case ActionCreateTimelineEvent,
		ActionCreateContact,
		ActionCreateOrganization:
		return FlowAction(s), nil
	default:
		return "", fmt.Errorf("invalid FlowAction: %s", s)
	}
}
