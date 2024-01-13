package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
)

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
		return neo4jutil.NodeLabelContact
	case USER:
		return neo4jutil.NodeLabelUser
	case ORGANIZATION:
		return neo4jutil.NodeLabelOrganization
	case MEETING:
		return neo4jutil.NodeLabelMeeting
	case COMMENT:
		return neo4jutil.NodeLabelComment
	case ISSUE:
		return neo4jutil.NodeLabelIssue
	case LOG_ENTRY:
		return neo4jutil.NodeLabelLogEntry
	case INTERACTION_EVENT:
		return neo4jutil.NodeLabelInteractionEvent
	}
	return "Unknown"
}
