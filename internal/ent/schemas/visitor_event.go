package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
		field.Int64("page_id"),
		field.String("event_type"),
		field.String("event_category"),
		field.String("event_name"),
		field.Time("occurred_at").Optional(),
		field.Time("server_received_at").Optional(),
		field.Int("sequence"),
		field.String("schema_version").Optional(),
		field.String("sdk_name").Optional(),
		field.String("sdk_version").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the VisitorEvent.
func (VisitorEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("custom_props", EventCustomProp.Type),
		edge.To("platform_clicks", VisitorPlatformClick.Type).Unique(),
		edge.To("target_props", EventTargetProp.Type).Unique(),
		edge.To("commerce_items", EventCommerceItem.Type),
		edge.To("search_filters", EventTargetSearchFilter.Type),
	}
}
