package enum

import (
	"fmt"
	"strings"
)

type AgentListenerEvent string

const (
	EventDoesNotNeedHelp         AgentListenerEvent = "does_not_need_help"
	EventICPFit                  AgentListenerEvent = "icp_fit"
	EventICPNotAFit              AgentListenerEvent = "icp_not_a_fit"
	EventNeedsHelp               AgentListenerEvent = "needs_help"
	EventNewLead                 AgentListenerEvent = "new_lead"
	EventNewMeetingRecording     AgentListenerEvent = "new_meeting_recording"
	EventNewSupportVisit         AgentListenerEvent = "new_support_visit"
	EventNewWebSession           AgentListenerEvent = "new_web_session"
	EventRunICPQualifierAgent    AgentListenerEvent = "run_icp_qualifier_agent"
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
		EventDoesNotNeedHelp,
		EventICPFit,
		EventICPNotAFit,
		EventNeedsHelp,
		EventNewLead,
		EventNewMeetingRecording,
		EventNewWebSession,
		EventNewSupportVisit,
		EventWebVisitorIdentified,
		EventWebVisitorNotIdentified,

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
