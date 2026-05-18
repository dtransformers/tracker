package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// VisitorConsent holds the schema definition for the VisitorConsent entity.
type VisitorConsent struct {
	ent.Schema
}

// Fields of the VisitorConsent.
func (VisitorConsent) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.String("uuid").Unique().NotEmpty(),
		field.Int64("visitor_id"),
		field.Int64("site_id"),
		field.String("consent_type"),
		field.String("status"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the VisitorConsent.
func (VisitorConsent) Edges() []ent.Edge {
	return nil
}
