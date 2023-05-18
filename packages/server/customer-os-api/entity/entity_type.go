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
		return "Contact"
	case USER:
		return "User"
	case ORGANIZATION:
		return "Organization"
	case MEETING:
		return "Meeting"
	}
	return "Unknown"
}
