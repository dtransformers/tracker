package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// VisitorPlatformClick holds the schema definition for the VisitorPlatformClick entity.
type VisitorPlatformClick struct {
	ent.Schema
}

// Fields of the VisitorPlatformClick.
func (VisitorPlatformClick) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.Int64("event_id").Unique(), // O2O
		field.String("gclid").Optional(),
		field.String("gbraid").Optional(),
		field.String("wbraid").Optional(),
		field.String("dclid").Optional(),
		field.String("fbclid").Optional(),
		field.String("fbp").Optional(),
		field.String("fbc").Optional(),
		field.String("ttclid").Optional(),
		field.String("ttp").Optional(),
		field.String("msclkid").Optional(),
		field.String("twclid").Optional(),
		field.String("li_fat_id").Optional(),
		field.String("sc_cid").Optional(),
		field.String("epik").Optional(),
		field.String("liveramp_ats_id").Optional(),
	}
}

// Edges of the VisitorPlatformClick.
func (VisitorPlatformClick) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", VisitorEvent.Type).
			Ref("platform_clicks").
			Field("event_id").
			Unique().
			Required(),
	}
}
