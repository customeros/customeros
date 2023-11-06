package entity

type EntityType string

const (
	CONTACT           EntityType = "CONTACT"
	USER              EntityType = "USER"
	ORGANIZATION      EntityType = "ORGANIZATION"
	MEETING           EntityType = "MEETING"
	COMMENT           EntityType = "COMMENT"
	ISSUE             EntityType = "ISSUE"
	LOG_ENTRY         EntityType = "LOG_ENTRY"
	INTERACTION_EVENT EntityType = "INTERACTION_EVENT"
)

func (entityType EntityType) String() string {
	return string(entityType)
}

func (entityType EntityType) Neo4jLabel() string {
	switch entityType {
	case CONTACT:
		return NodeLabel_Contact
	case USER:
		return NodeLabel_User
	case ORGANIZATION:
		return NodeLabel_Organization
	case MEETING:
		return NodeLabel_Meeting
	case COMMENT:
		return NodeLabel_Comment
	case ISSUE:
		return NodeLabel_Issue
	case LOG_ENTRY:
		return NodeLabel_LogEntry
	case INTERACTION_EVENT:
		return NodeLabel_InteractionEvent
	}
	return "Unknown"
}
