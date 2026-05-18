package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type VisitorSession struct {
	ent.Schema
}

func (VisitorSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("uuid").
			Unique().
			NotEmpty(),

		field.Int64("visitor_id").
			Positive(),

		field.Int64("device_id").
			Positive(),

		field.String("session_token").
			Unique().
			NotEmpty(),

		field.Time("started_at").
			Default(time.Now),

		field.Time("ended_at").
			Optional().
			Nillable(),

		field.Time("last_activity_at").
			Default(time.Now),

		field.String("entry_page").
			Optional(),

		field.String("exit_page").
			Optional(),

		field.String("acquisition_source").
			Optional(),

		field.String("acquisition_medium").
			Optional(),

		field.String("acquisition_campaign").
			Optional(),

		field.String("referrer_host").
			Optional(),

		field.Time("created_at").
			Default(time.Now),
	}
}
