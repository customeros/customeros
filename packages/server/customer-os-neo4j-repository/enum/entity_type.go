package enum

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

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
		return entity.NodeLabelContact
	case USER:
		return entity.NodeLabelUser
	case ORGANIZATION:
		return entity.NodeLabelOrganization
	case MEETING:
		return entity.NodeLabelMeeting
	case CONTRACT:
		return entity.NodeLabelContract
	}
	return "Unknown"
}
