package entity

type EntityType string

const (
	CONTACT      EntityType = "CONTACT"
	USER         EntityType = "USER"
	ORGANIZATION EntityType = "ORGANIZATION"
	MEETING      EntityType = "MEETING"
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
	}
	return "Unknown"
}
