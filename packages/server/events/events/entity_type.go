package events

type EntityType string

const (
	CONTACT EntityType = "contact"
)

func (entityType EntityType) String() string {
	return string(entityType)
}
