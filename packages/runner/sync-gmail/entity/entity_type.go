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
