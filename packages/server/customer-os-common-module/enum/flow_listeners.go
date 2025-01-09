package enum

import (
	"fmt"
	"strings"
)

type FlowListenerEvent string

const (
	EventFathomMeetingSummaryCreated FlowListenerEvent = "fathom.meeting_summary.created"
	EventFlowContactAdded            FlowListenerEvent = "flow.contact.added"
	EventGrainMeetingSummaryCreated  FlowListenerEvent = "grain.meeting_summary.created"
	EventRevealWebsiteVisit          FlowListenerEvent = "reveal.website_visit"
	NotSet                           FlowListenerEvent = ""
)

func (e FlowListenerEvent) String() string {
	return string(e)
}

func (e FlowListenerEvent) Parse() (system, resource, action string, err error) {
	parts := strings.Split(e.String(), ".")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid event name format: %s", e.String())
	}

	return parts[0], parts[1], parts[2], nil
}

func (e FlowListenerEvent) ExternalSystem() (system Source, err error) {
	systemId, _, _, err := e.Parse()
	if err != nil {
		return SourceUnknown, fmt.Errorf("invalid event name")
	}

	return DecodeSource(systemId), nil
}

func GetFlowListenerEvent(s string) (FlowListenerEvent, error) {
	switch FlowListenerEvent(s) {
	case
		EventFathomMeetingSummaryCreated,
		EventFlowContactAdded,
		EventGrainMeetingSummaryCreated,
		EventRevealWebsiteVisitNew,
		EventRevealWebsiteVisitRepeat:

		return FlowListenerEvent(s), nil
	default:
		return "", fmt.Errorf("invalid FlowListenerEvent: %s", s)
	}
}
