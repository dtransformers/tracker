package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// PageView holds the schema definition for the PageView entity.
type PageView struct {
	ent.Schema
}

// Fields of the PageView.
func (PageView) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.String("uuid").Unique().NotEmpty(),
		field.Int64("session_id"),
		field.Int64("visitor_id"),
		field.Int64("device_id"),
		field.String("url").Optional(),
		field.String("title").Optional(),
		field.String("referrer").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the PageView.
func (PageView) Edges() []ent.Edge {
	return nil
}
