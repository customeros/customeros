package enum

import "fmt"

type FlowAction string

const (
	ActionCreateTimelineEvent FlowAction = "create.timeline_event"
	ActionCreateContacts      FlowAction = "create.contacts"
	ActionCreateOrganizations FlowAction = "create_organizations"
)

func (t FlowAction) String() string {
	return string(t)
}

func GetFlowAction(s string) (FlowAction, error) {
	switch FlowAction(s) {
	case ActionCreateTimelineEvent,
		ActionCreateContacts,
		ActionCreateOrganizations:
		return FlowAction(s), nil
	default:
		return "", fmt.Errorf("invalid FlowAction: %s", s)
	}
}
