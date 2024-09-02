package model

type EntityRelation string

const (
	HAS        EntityRelation = "HAS"
	INCLUDES   EntityRelation = "INCLUDES"
	RECORDING  EntityRelation = "RECORDING"
	PART_OF    EntityRelation = "PART_OF"
	REPLIES_TO EntityRelation = "REPLIES_TO"
	SENT_BY    EntityRelation = "SENT_BY"
	SENT_TO    EntityRelation = "SENT_TO"
)

func (entityRelation EntityRelation) String() string {
	return string(entityRelation)
}
