package entity

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
		return NodeLabelContact
	case USER:
		return NodeLabelUser
	case ORGANIZATION:
		return NodeLabelOrganization
	case MEETING:
		return NodeLabelMeeting
	case CONTRACT:
		return NodeLabelContract
	}
	return "Unknown"
}
