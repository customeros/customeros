package events

type EntityType string

const (
	CONTACT EntityType = "contact"
	USER    EntityType = "user"
)

func (entityType EntityType) String() string {
	return string(entityType)
}
