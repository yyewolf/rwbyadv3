package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Trade holds the schema definition for the Trade entity.
type Trade struct {
	ent.Schema
}

// Fields of the Trade.
func (Trade) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("initiator_id"),
		field.String("receiver_id"),
		field.Int("liens").Default(0),
		field.Strings("offer_cards").Default([]string{}),
		field.Strings("receive_cards").Default([]string{}),
	}
}

func (Trade) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Trade.
func (Trade) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("initiator", Player.Type).
			Unique().
			Required().
			Field("initiator_id"),
		edge.To("receiver", Player.Type).
			Unique().
			Required().
			Field("receiver_id"),
	}
}
