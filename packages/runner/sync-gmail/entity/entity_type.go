package entity

type EntityType string

const (
	ORGANIZATION EntityType = "ORGANIZATION"
)

func (entityType EntityType) String() string {
	return string(entityType)
}

func (entityType EntityType) Neo4jLabel() string {
	switch entityType {
	case ORGANIZATION:
		return NodeLabel_Organization
	}
	return "Unknown"
}
