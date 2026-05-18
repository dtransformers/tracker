package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Visitor struct {
	ent.Schema
}

func (Visitor) Fields() []ent.Field {
	return []ent.Field{
		field.String("uuid").
			Unique().
			NotEmpty(),

		field.String("anonymous_id").
			Unique().
			NotEmpty(),

		field.String("primary_user_id").
			Optional().
			Nillable(),

		field.Time("first_seen_at").
			Default(time.Now),

		field.Time("last_seen_at").
			Default(time.Now),

		field.Bool("is_bot").
			Default(false),

		field.Bool("is_verified_human").
			Default(false),

		field.Time("created_at").
			Default(time.Now),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
