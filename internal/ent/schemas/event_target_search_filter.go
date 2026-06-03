package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EventTargetSearchFilter holds the schema definition for the EventTargetSearchFilter entity.
type EventTargetSearchFilter struct {
	ent.Schema
}

// Fields of the EventTargetSearchFilter.
func (EventTargetSearchFilter) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.Int64("event_id"),
		field.String("key").NotEmpty(),
		field.String("value").Optional(),
	}
}

// Edges of the EventTargetSearchFilter.
func (EventTargetSearchFilter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", VisitorEvent.Type).
			Ref("search_filters").
			Field("event_id").
			Unique().
			Required(),
	}
}
