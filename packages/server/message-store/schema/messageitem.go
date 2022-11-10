package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MessageItem holds the schema definition for the MessageItem entity.
type MessageItem struct {
	ent.Schema
}

// Fields of the MessageItem.
func (MessageItem) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").Values("MESSAGE", "FILE"),
		field.String("username"),
		field.String("message").
			SchemaType(map[string]string{
				dialect.Postgres: "text", // Override Postgres.
			}),
		field.Enum("channel").Values("CHAT", "MAIL", "WHATSAPP", "FACEBOOK", "TWITTER", "VOICE"),
		field.Enum("direction").Values("INBOUND", "OUTBOUND"),
		field.Time("time").Optional().
			Annotations(
				&entsql.Annotation{
					Default: "CURRENT_TIMESTAMP",
				},
			),
	}
}

// Edges of the MessageItem.
func (MessageItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("message_feed", MessageFeed.Type).
			Ref("message_item").
			Unique().
			Required().
			Immutable(),
	}
}

func (MessageItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("time").
			Edges("message_feed"),
	}
}
