package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MessageFeed holds the schema definition for the MessageFeed entity.
type Conversation struct {
	ent.Schema
}

// Fields of the MessageFeed.
func (Conversation) Fields() []ent.Field {
	return []ent.Field{
		field.String("contactId").Unique(),
		field.Time("createdOn").Annotations(&entsql.Annotation{Default: "CURRENT_TIMESTAMP"}),
		field.Time("updatedOn").Annotations(&entsql.Annotation{Default: "CURRENT_TIMESTAMP"}),
		field.Enum("state").Values("NEW", "IN_PROGRESS", "CLOSED"),

		field.String("lastMessage").SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.String("lastSenderId"),
		field.Enum("lastSenderType").Values("CONTACT", "USER"),
	}
}

// Edges of the MessageFeed.
func (Conversation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("conversation_item", ConversationItem.Type),
	}
}

func (Conversation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("contactId").Unique(),
	}
}
