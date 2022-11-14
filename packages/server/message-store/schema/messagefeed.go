package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MessageFeed holds the schema definition for the MessageFeed entity.
type MessageFeed struct {
	ent.Schema
}

// Fields of the MessageFeed.
func (MessageFeed) Fields() []ent.Field {
	return []ent.Field{
		field.String("contactId").
			Unique(),
		field.String("firstName"),
		field.String("lastName"),
	}
}

// Edges of the MessageFeed.
func (MessageFeed) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("message_item", MessageItem.Type),
	}
}

func (MessageFeed) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("contactId").Unique(),
	}
}
