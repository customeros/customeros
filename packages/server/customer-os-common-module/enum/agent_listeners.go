package enum

import (
	"fmt"
	"strings"
)

type AgentListenerEvent string

const (
	EventNewLead                 AgentListenerEvent = "new_lead"
	EventICPFit                  AgentListenerEvent = "icp_fit"
	EventICPNotAFit              AgentListenerEvent = "icp_not_a_fit"
	EventRunICPQualifierAgent    AgentListenerEvent = "run_icp_qualifier_agent"
	EventNewWebSession           AgentListenerEvent = "new_web_session"
	EventWebVisitorIdentified    AgentListenerEvent = "web_visitor_identified"
	EventWebVisitorNotIdentified AgentListenerEvent = "web_visitor_not_identified"

	EventFathomMeetingSummaryCreated AgentListenerEvent = "fathom.meeting_summary.created"
	EventFlowContactAdded            AgentListenerEvent = "flow.contact.added"
	EventGrainMeetingSummaryCreated  AgentListenerEvent = "grain.meeting_summary.created"
	EventIntentSignal                AgentListenerEvent = "generic.intent_signal.detected"
	NotSet                           AgentListenerEvent = ""
)

func (e AgentListenerEvent) String() string {
	return string(e)
}

func (e AgentListenerEvent) Parse() (system, resource, action string, err error) {
	parts := strings.Split(e.String(), ".")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid event name format: %s", e.String())
	}

	return parts[0], parts[1], parts[2], nil
}

func (e AgentListenerEvent) ExternalSystem() (system Source, err error) {
	systemId, _, _, err := e.Parse()
	if err != nil {
		return SourceUnknown, fmt.Errorf("invalid event name")
	}

	return DecodeSource(systemId), nil
}

func GetAgentListener(s string) (AgentListenerEvent, error) {
	switch AgentListenerEvent(s) {
	case
		EventNewLead,
		EventNewWebSession,
		EventWebVisitorIdentified,
		EventWebVisitorNotIdentified,
		EventICPFit,
		EventICPNotAFit,
		EventRunICPQualifierAgent,
		EventFathomMeetingSummaryCreated,
		EventFlowContactAdded,
		EventGrainMeetingSummaryCreated,
		EventIntentSignal:

		return AgentListenerEvent(s), nil
	default:
		return "", fmt.Errorf("invalid AgentListener: %s", s)
	}
}
