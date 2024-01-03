package entity

import neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

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
	CONTRACT          EntityType = "CONTRACT"
	OPPORTUNITY       EntityType = "OPPORTUNITY"
	SERVICE_LINE_ITEM EntityType = "SERVICE_LINE_ITEM"
)

func (entityType EntityType) String() string {
	return string(entityType)
}

func (entityType EntityType) Neo4jLabel() string {
	switch entityType {
	case CONTACT:
		return neo4jentity.NodeLabel_Contact
	case USER:
		return neo4jentity.NodeLabel_User
	case ORGANIZATION:
		return neo4jentity.NodeLabel_Organization
	case MEETING:
		return neo4jentity.NodeLabel_Meeting
	case COMMENT:
		return neo4jentity.NodeLabel_Comment
	case ISSUE:
		return neo4jentity.NodeLabel_Issue
	case LOG_ENTRY:
		return neo4jentity.NodeLabel_LogEntry
	case INTERACTION_EVENT:
		return neo4jentity.NodeLabel_InteractionEvent
	}
	return "Unknown"
}
