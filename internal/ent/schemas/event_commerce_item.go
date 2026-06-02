package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EventCommerceItem holds the schema definition for the EventCommerceItem entity.
type EventCommerceItem struct {
	ent.Schema
}

// Fields of the EventCommerceItem.
func (EventCommerceItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.Int64("event_id"),
		field.String("item_id"),
		field.String("item_name").Optional(),
		field.String("item_brand").Optional(),
		field.String("item_category").Optional(),
		field.String("item_category_2").Optional(),
		field.String("item_category_3").Optional(),
		field.String("item_variant").Optional(),
		field.Float("price").Optional(),
		field.String("currency").Optional(),
		field.Int("quantity").Optional(),
		field.String("coupon").Optional(),
		field.Float("discount").Optional(),
		field.Int("index").Optional(),
		field.String("item_list_id").Optional(),
		field.String("item_list_name").Optional(),
		field.String("affiliation").Optional(),
	}
}

// Edges of the EventCommerceItem.
func (EventCommerceItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", VisitorEvent.Type).
			Ref("commerce_items").
			Field("event_id").
			Unique().
			Required(),
	}
}
