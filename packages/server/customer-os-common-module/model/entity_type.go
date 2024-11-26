package model

type EntityType string

const (
	TENANT                EntityType = "TENANT"
	TENANT_SETTINGS       EntityType = "TENANT_SETTINGS"
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
	MARKDOWN_EVENT        EntityType = "MARKDOWN_EVENT"
	OPPORTUNITY           EntityType = "OPPORTUNITY"
	SERVICE_LINE_ITEM     EntityType = "SERVICE_LINE_ITEM"
	REMINDER              EntityType = "REMINDER"
	ATTACHMENT            EntityType = "ATTACHMENT"
	NOTE                  EntityType = "NOTE"
	FLOW                  EntityType = "FLOW"
	FLOW_ACTION           EntityType = "FLOW_ACTION"
	FLOW_PARTICIPANT      EntityType = "FLOW_PARTICIPANT"
	FLOW_SENDER           EntityType = "FLOW_SENDER"
	CUSTOM_FIELD          EntityType = "CUSTOM_FIELD"
	CUSTOM_FIELD_TEMPLATE EntityType = "CUSTOM_FIELD_TEMPLATE"
	SOCIAL                EntityType = "SOCIAL"
	MAILSTACK_BUY_REQUEST EntityType = "MAILSTACK_BUY_REQUEST"
)

func (entityType EntityType) String() string {
	return string(entityType)
}

func (entityType EntityType) Neo4jLabel() string {
	switch entityType {
	case TENANT:
		return NodeLabelTenant
	case TENANT_SETTINGS:
		return NodeLabelTenantSettings
	case CONTACT:
		return NodeLabelContact
	case USER:
		return NodeLabelUser
	case ORGANIZATION:
		return NodeLabelOrganization
	case OPPORTUNITY:
		return NodeLabelOpportunity
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
	case FLOW_ACTION:
		return NodeLabelFlowAction
	case FLOW_PARTICIPANT:
		return NodeLabelFlowParticipant
	case FLOW_SENDER:
		return NodeLabelFlowSender
	case CUSTOM_FIELD:
		return NodeLabelCustomField
	case CUSTOM_FIELD_TEMPLATE:
		return NodeLabelCustomFieldTemplate
	case SOCIAL:
		return NodeLabelSocial
	case MARKDOWN_EVENT:
		return NodeLabelMarkdownEvent
	}
	return "Unknown"
}

func DecodeEntityType(s string) EntityType {
	return EntityType(s)
}
