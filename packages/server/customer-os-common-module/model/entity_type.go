package model

type EntityType string

const (
	CONTACT               EntityType = "CONTACT"
	USER                  EntityType = "USER"
	ORGANIZATION          EntityType = "ORGANIZATION"
	EMAIL                 EntityType = "EMAIL"
	PHONE_NUMBER          EntityType = "PHONE_NUMBER"
	MEETING               EntityType = "MEETING"
	CONTRACT              EntityType = "CONTRACT"
	INVOICE               EntityType = "INVOICE"
	INTERACTION_EVENT     EntityType = "INTERACTION_EVENT"
	INTERACTION_SESSION   EntityType = "INTERACTION_SESSION"
	COMMENT               EntityType = "COMMENT"
	ISSUE                 EntityType = "ISSUE"
	LOG_ENTRY             EntityType = "LOG_ENTRY"
	OPPORTUNITY           EntityType = "OPPORTUNITY"
	SERVICE_LINE_ITEM     EntityType = "SERVICE_LINE_ITEM"
	REMINDER              EntityType = "REMINDER"
	ATTACHMENT            EntityType = "ATTACHMENT"
	NOTE                  EntityType = "NOTE"
	FLOW                  EntityType = "FLOW"
	FLOW_SEQUENCE         EntityType = "FLOW_SEQUENCE"
	FLOW_SEQUENCE_CONTACT EntityType = "FLOW_SEQUENCE_CONTACT"
	FLOW_SEQUENCE_SENDER  EntityType = "FLOW_SEQUENCE_SENDER"
)

func (entityType EntityType) String() string {
	return string(entityType)
}

func (entityType EntityType) Neo4jLabel() string {
	switch entityType {
	case CONTACT:
		return NodeLabelContact
	case USER:
		return NodeLabelUser
	case ORGANIZATION:
		return NodeLabelOrganization
	case EMAIL:
		return NodeLabelEmail
	case PHONE_NUMBER:
		return NodeLabelPhoneNumber
	case MEETING:
		return NodeLabelMeeting
	case CONTRACT:
		return NodeLabelContract
	case INVOICE:
		return NodeLabelInvoice
	case INTERACTION_EVENT:
		return NodeLabelInteractionEvent
	case INTERACTION_SESSION:
		return NodeLabelInteractionSession
	case COMMENT:
		return NodeLabelComment
	case ISSUE:
		return NodeLabelIssue
	case LOG_ENTRY:
		return NodeLabelLogEntry
	case REMINDER:
		return NodeLabelReminder
	case ATTACHMENT:
		return NodeLabelAttachment
	case NOTE:
		return NodeLabelNote
	case FLOW:
		return NodeLabelFlow
	case FLOW_SEQUENCE:
		return NodeLabelFlowSequence
	case FLOW_SEQUENCE_CONTACT:
		return NodeLabelFlowSequenceContact
	case FLOW_SEQUENCE_SENDER:
		return NodeLabelFlowSequenceSender
	}
	return "Unknown"
}
