package events

type EntityType string

const (
	CONTACT           EntityType = "CONTACT"
	USER              EntityType = "USER"
	ORGANIZATION      EntityType = "ORGANIZATION"
	MEETING           EntityType = "MEETING"
	CONTRACT          EntityType = "CONTRACT"
	INVOICE           EntityType = "INVOICE"
	INTERACTION_EVENT EntityType = "INTERACTION_EVENT"
)

func (entityType EntityType) String() string {
	return string(entityType)
}
