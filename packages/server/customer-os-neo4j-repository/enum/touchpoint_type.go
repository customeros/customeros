package enum

type TouchpointType string

const (
	TouchpointTypePageView                      TouchpointType = "PAGE_VIEW"
	TouchpointTypeInteractionSession            TouchpointType = "INTERACTION_SESSION"
	TouchpointTypeNote                          TouchpointType = "NOTE"
	TouchpointTypeInteractionEventEmailSent     TouchpointType = "INTERACTION_EVENT_EMAIL_SENT"
	TouchpointTypeInteractionEventEmailReceived TouchpointType = "INTERACTION_EVENT_EMAIL_RECEIVED"
	TouchpointTypeInteractionEventPhoneCall     TouchpointType = "INTERACTION_EVENT_PHONE_CALL"
	TouchpointTypeInteractionEventChat          TouchpointType = "INTERACTION_EVENT_CHAT"
	TouchpointTypeMeeting                       TouchpointType = "MEETING"
	TouchpointTypeActionCreated                 TouchpointType = "ACTION_CREATED"
	TouchpointTypeAction                        TouchpointType = "ACTION"
	TouchpointTypeLogEntry                      TouchpointType = "LOG_ENTRY"
	TouchpointTypeIssueCreated                  TouchpointType = "ISSUE_CREATED"
	TouchpointTypeIssueUpdated                  TouchpointType = "ISSUE_UPDATED"
)

func (e TouchpointType) String() string {
	return string(e)
}

func DecodeTouchpointType(str string) TouchpointType {
	switch str {
	case TouchpointTypePageView.String():
		return TouchpointTypePageView
	case TouchpointTypeInteractionSession.String():
		return TouchpointTypeInteractionSession
	case TouchpointTypeNote.String():
		return TouchpointTypeNote
	case TouchpointTypeInteractionEventEmailSent.String():
		return TouchpointTypeInteractionEventEmailSent
	case TouchpointTypeInteractionEventEmailReceived.String():
		return TouchpointTypeInteractionEventEmailReceived
	case TouchpointTypeInteractionEventPhoneCall.String():
		return TouchpointTypeInteractionEventPhoneCall
	case TouchpointTypeInteractionEventChat.String():
		return TouchpointTypeInteractionEventChat
	case TouchpointTypeMeeting.String():
		return TouchpointTypeMeeting
	case TouchpointTypeActionCreated.String():
		return TouchpointTypeActionCreated
	case TouchpointTypeAction.String():
		return TouchpointTypeAction
	case TouchpointTypeLogEntry.String():
		return TouchpointTypeLogEntry
	case TouchpointTypeIssueCreated.String():
		return TouchpointTypeIssueCreated
	case TouchpointTypeIssueUpdated.String():
		return TouchpointTypeIssueUpdated
	default:
		return ""
	}
}
