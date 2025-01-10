package model

type EntityType string

const (
	ATTACHMENT               EntityType = "ATTACHMENT"
	COMMENT                  EntityType = "COMMENT"
	CONTACT                  EntityType = "CONTACT"
	CONTRACT                 EntityType = "CONTRACT"
	CUSTOM_FIELD             EntityType = "CUSTOM_FIELD"
	CUSTOM_FIELD_TEMPLATE    EntityType = "CUSTOM_FIELD_TEMPLATE"
	DOMAIN                   EntityType = "DOMAIN"
	EMAIL                    EntityType = "EMAIL"
	FLOW                     EntityType = "FLOW"
	FLOW_ACTION              EntityType = "FLOW_ACTION"
	FLOW_ACTION_EVENT        EntityType = "FLOW_ACTION_EVENT"
	FLOW_ACTION_RESULT_EVENT EntityType = "FLOW_ACTION_RESULT_EVENT"
	FLOW_PARTICIPANT         EntityType = "FLOW_PARTICIPANT"
	FLOW_SENDER              EntityType = "FLOW_SENDER"
	INTERACTION_EVENT        EntityType = "INTERACTION_EVENT"
	INTERACTION_SESSION      EntityType = "INTERACTION_SESSION"
	INVOICE                  EntityType = "INVOICE"
	ISSUE                    EntityType = "ISSUE"
	LOG_ENTRY                EntityType = "LOG_ENTRY"
	MAILBOX                  EntityType = "MAILBOX"
	MAILSTACK_BUY_REQUEST    EntityType = "MAILSTACK_BUY_REQUEST"
	MARKDOWN_EVENT           EntityType = "MARKDOWN_EVENT"
	MEETING                  EntityType = "MEETING"
	NOTE                     EntityType = "NOTE"
	OPPORTUNITY              EntityType = "OPPORTUNITY"
	ORGANIZATION             EntityType = "ORGANIZATION"
	PHONE_NUMBER             EntityType = "PHONE_NUMBER"
	REMINDER                 EntityType = "REMINDER"
	SERVICE_LINE_ITEM        EntityType = "SERVICE_LINE_ITEM"
	SOCIAL                   EntityType = "SOCIAL"
	TAG                      EntityType = "TAG"
	TENANT                   EntityType = "TENANT"
	TENANT_SETTINGS          EntityType = "TENANT_SETTINGS"
	USER                     EntityType = "USER"
	WEBHOOK_EVENT            EntityType = "WEBHOOK"
	LOCATION                 EntityType = "LOCATION"
	JOB_ROLE                 EntityType = "JOB_ROLE"
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
	case DOMAIN:
		return NodeLabelDomain
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
	case LOCATION:
		return NodeLabelLocation
	case TAG:
		return NodeLabelTag
	case JOB_ROLE:
		return NodeLabelJobRole
	}
	return "Unknown"
}

func DecodeEntityType(s string) EntityType {
	return EntityType(s)
}
