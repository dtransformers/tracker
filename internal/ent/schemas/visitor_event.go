package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// VisitorEvent holds the schema definition for the VisitorEvent entity.
type VisitorEvent struct {
	ent.Schema
}

// Fields of the VisitorEvent.
func (VisitorEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.String("uuid").Unique().NotEmpty(),
		field.Int64("session_id"),
		field.Int64("visitor_id"),
		field.Int64("device_id"),
		field.String("event_name"),
		field.JSON("event_properties", map[string]interface{}{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the VisitorEvent.
func (VisitorEvent) Edges() []ent.Edge {
	return nil
}
