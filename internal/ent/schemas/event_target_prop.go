package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EventTargetProp holds the schema definition for the EventTargetProp entity.
type EventTargetProp struct {
	ent.Schema
}

// Fields of the EventTargetProp.
func (EventTargetProp) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.Int64("event_id").Unique(), // O2O

		// DOM Element Identification
		field.String("element_type").Optional(),
		field.String("element_id").Optional(),
		field.String("element_name").Optional(),
		field.String("element_text").Optional(),
		field.String("element_class").Optional(),
		field.String("element_selector").Optional(),
		field.String("component").Optional(),

		// Human-Readable Semantic Labels
		field.String("page_name").Optional(),
		field.String("button_name").Optional(),
		field.String("button_id").Optional(),
		field.String("href").Optional(),

		// Form
		field.String("form_name").Optional(),
		field.String("form_id").Optional(),
		field.String("form_field").Optional(),
		field.String("form_value_hash").Optional(),
		field.Int("form_field_count").Optional(),
		field.Float("form_completion_pct").Optional(),

		// Media / Video
		field.String("media_type").Optional(),
		field.String("media_id").Optional(),
		field.String("media_title").Optional(),
		field.String("media_provider").Optional(),
		field.Int64("media_duration_ms").Optional(),
		field.Int64("media_current_time_ms").Optional(),
		field.Float("media_percent").Optional(),
		field.String("media_quality").Optional(),
		field.Bool("media_sound_on").Optional(),
		field.Bool("media_fullscreen").Optional(),

		// Scroll
		field.Float("scroll_percent").Optional(),
		field.String("scroll_direction").Optional(),
		field.Int("scroll_depth_px").Optional(),

		// Mouse / Pointer
		field.String("mouse_button").Optional(),
		field.Int("click_x").Optional(),
		field.Int("click_y").Optional(),
		field.Int("element_x").Optional(),
		field.Int("element_y").Optional(),

		// Commerce
		field.String("transaction_id").Optional(),
		field.String("currency").Optional(),
		field.Float("value").Optional(),
		field.Float("tax").Optional(),
		field.Float("shipping").Optional(),
		field.String("coupon").Optional(),
		field.String("affiliation").Optional(),
		field.String("payment_type").Optional(),
		field.String("shipping_tier").Optional(),

		// Search
		field.String("search_term").Optional(),
		field.Int("search_results_count").Optional(),

		// Authentication
		field.String("auth_provider").Optional(),
		field.String("auth_method").Optional(),
		field.Bool("is_new_user").Optional(),

		// Error / System
		field.String("error_code").Optional(),
		field.String("error_message").Optional(),
		field.String("error_severity").Optional(),
		field.String("error_stack").Optional(),
		field.Int("http_status").Optional(),
		field.String("component_name").Optional(),

		// Engagement / Sharing
		field.String("share_platform").Optional(),
		field.String("share_method").Optional(),

		// Timing
		field.Int64("duration_ms").Optional(),

		// Feature Flags
		field.String("feature").Optional(),
		field.String("priority").Optional(),
	}
}

// Edges of the EventTargetProp.
func (EventTargetProp) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", VisitorEvent.Type).
			Ref("target_props").
			Field("event_id").
			Unique().
			Required(),
	}
}
