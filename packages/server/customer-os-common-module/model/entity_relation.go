package model

type EntityRelation string

const (
	BELONGS_TO_TENANT EntityRelation = "BELONGS_TO_TENANT"
	HAS               EntityRelation = "HAS"
	HAS_OPPORTUNITY   EntityRelation = "HAS_OPPORTUNITY"
	INCLUDES          EntityRelation = "INCLUDES"
	RECORDING         EntityRelation = "RECORDING"
	PART_OF           EntityRelation = "PART_OF"
	REPLIES_TO        EntityRelation = "REPLIES_TO"
	SENT_BY           EntityRelation = "SENT_BY"
	SENT_TO           EntityRelation = "SENT_TO"
	NEXT              EntityRelation = "NEXT"
	CONNECTED_WITH    EntityRelation = "CONNECTED_WITH"
)

func (entityRelation EntityRelation) String() string {
	return string(entityRelation)
}
