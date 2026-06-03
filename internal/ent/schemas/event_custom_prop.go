package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EventCustomProp holds the schema definition for the EventCustomProp entity.
type EventCustomProp struct {
	ent.Schema
}

// Fields of the EventCustomProp.
func (EventCustomProp) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.Int64("event_id"),
		field.String("key").NotEmpty(),
		field.String("value").Optional(),
	}
}

// Edges of the EventCustomProp.
func (EventCustomProp) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", VisitorEvent.Type).
			Ref("custom_props").
			Field("event_id").
			Unique().
			Required(),
	}
}
