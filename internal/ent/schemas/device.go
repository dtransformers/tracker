package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Device struct {
	ent.Schema
}

func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("uuid").
			Unique().
			NotEmpty(),

		field.Int64("visitor_id").
			Positive(),

		field.String("fingerprint_hash").
			NotEmpty(),

		field.String("device_type").
			Optional(),

		field.String("os").
			Optional(),

		field.String("address_ip").
			Optional(),

		field.String("country").
			Optional(),

		field.String("region").
			Optional(),

		field.String("city").
			Optional(),

		field.String("browser").
			Optional(),

		field.String("language").
			Optional(),

		field.String("timezone").
			Optional(),

		field.Time("first_seen_at").
			Default(time.Now),

		field.Time("last_seen_at").
			Default(time.Now),
	}
}
