package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Site struct {
	ent.Schema
}

func (Site) Fields() []ent.Field {
	return []ent.Field{

		field.String("uuid").
			Unique().
			NotEmpty(),

		field.String("name").
			NotEmpty(),

		field.String("domain").
			Unique().
			NotEmpty(),

		field.Int64("tenant_id").
			Positive().
			Default(1),

		field.String("public_key").
			Unique().
			NotEmpty(),

		field.String("secret_key").
			Sensitive().
			NotEmpty(),

		field.Bool("is_active").
			Default(true),

		field.Bool("is_deleted").
			Default(false),

		field.Time("created_at").
			Default(time.Now),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),

		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}
