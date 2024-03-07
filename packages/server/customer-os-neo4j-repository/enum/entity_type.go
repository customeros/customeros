package enum

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
)

type EntityType string

const (
	CONTACT      EntityType = "CONTACT"
	USER         EntityType = "USER"
	ORGANIZATION EntityType = "ORGANIZATION"
	MEETING      EntityType = "MEETING"
	CONTRACT     EntityType = "CONTRACT"
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
	case CONTRACT:
		return neo4jutil.NodeLabelContract
	}
	return "Unknown"
}
