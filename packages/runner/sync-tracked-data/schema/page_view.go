package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PageView holds the schema definition for the PageView entity.
type PageView struct {
	ent.Schema
}

// Annotations of the PageView.
func (PageView) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "page_views"},
	}
}

// Fields of the PageView.
func (PageView) Fields() []ent.Field {
	return []ent.Field{
		field.String(appId),
		field.String(trackerName),
		field.String(tenant),
		field.String(pageViewId),
		field.String(eventId),
		field.Bool(syncedToCustomerOs),
		field.Time(startTimestamp),
		field.Time(endTimestamp),
		field.String(visitorId).Optional(),
		field.String(customerOsContactId).Optional(),
		field.String(domainUserId),
		field.String(networkUserId),
		field.Int(domainSessionId),
		field.Int(domainSessionIdx),
		field.Int(pageViewsInSession),
		field.Int(pageViewInSessionIndex),
		field.Int(engagedTimeInSec),
		field.String(pageUrl),
		field.String(pageTitle),
	}
}

func (PageView) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(eventId).
			Unique(),
	}
}
