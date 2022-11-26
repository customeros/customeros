package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Visitor holds the schema definition for the Visitor entity.
type Visitor struct {
	ent.Schema
}

// Annotations of the Visitor.
func (Visitor) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "visitors"},
	}
}

// Fields of the Visitor.
func (Visitor) Fields() []ent.Field {
	return []ent.Field{
		field.String(appId),
		field.String(trackerName),
		field.String(tenant),
		field.String(visitorId).Optional(),
		field.String(customerOsContactId).Optional(),
		field.String(domainUserId),
		field.String(networkUserId),
		field.Int(pageViews),
		field.Int(sessions),
		field.Int(engagedTimeInSec),
		field.Bool(syncedToCustomerOs),
	}
}

func (Visitor) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(domainUserId, appId, trackerName),
	}
}
