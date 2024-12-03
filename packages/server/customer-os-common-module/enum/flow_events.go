package enum

import "fmt"

type FlowEvent string

const (
	EventFathomMeetingSummaryCreated FlowEvent = "fathom.meeting_summary.created"
	EventGrainMeetingSummaryCreated  FlowEvent = "grain.meeting_summary.created"
)

func (e FlowEvent) String() string {
	return string(e)
}

func GetFlowEvent(s string) (FlowEvent, error) {
	switch FlowEvent(s) {
	case EventFathomMeetingSummaryCreated, EventGrainMeetingSummaryCreated:
		return FlowEvent(s), nil
	default:
		return "", fmt.Errorf("invalid FlowEvent: %s", s)
	}
}
