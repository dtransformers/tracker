package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type SiteVisitor struct {
	ent.Schema
}

func (SiteVisitor) Fields() []ent.Field {
	return []ent.Field{
		field.String("uuid").
			Unique().
			NotEmpty(),

		field.Int64("visitor_id").
			Positive(),

		field.Int64("session_id").
			Positive(),

		field.Int64("site_id").
			Positive(),

		field.Time("created_at").
			Default(time.Now),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
