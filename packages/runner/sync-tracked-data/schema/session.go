package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Annotations of the Session.
func (Session) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sessions"},
	}
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String(appId),
		field.String(trackerName),
		field.String(tenant),
		field.String(domainSessionId),
		field.Int(domainSessionIdx),
		field.Bool(syncedToCustomerOs),
		field.Time(startTimestamp),
		field.Time(endTimestamp),
		field.String(domainUserId),
		field.String(networkUserId),
		field.String(visitorId).Optional(),
		field.String(customerOsContactId).Optional(),
	}
}

func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(domainSessionId, appId, trackerName).Unique(),
	}
}
