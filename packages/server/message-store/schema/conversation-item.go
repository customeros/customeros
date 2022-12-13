package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ConversationItem holds the schema definition for the ConversationItem entity.
type ConversationItem struct {
	ent.Schema
}

// Fields of the ConversationItem.
func (ConversationItem) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").Values("MESSAGE", "FILE"),
		field.String("senderId"),
		field.Enum("senderType").Values("CONTACT", "USER"),
		field.String("message").SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.Enum("channel").Values("CHAT", "MAIL", "WHATSAPP", "FACEBOOK", "TWITTER", "VOICE"),
		field.Enum("direction").Values("INBOUND", "OUTBOUND"),
		field.Time("time").Annotations(&entsql.Annotation{Default: "CURRENT_TIMESTAMP"}),
	}
}

// Edges of the ConversationItem.
func (ConversationItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("conversation", Conversation.Type).
			Ref("conversation_item").
			Unique().
			Required().
			Immutable(),
	}
}

func (ConversationItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("time").
			Edges("conversation"),
	}
}
